package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	blogmodel "fastduck/treasure-doc/module/blog/data/model"
	"fastduck/treasure-doc/module/blog_mgr/data/request"

	"gorm.io/gorm"
)

// MemoResult 是前台速记本的管理返回：Tags / Images 归一化为字符串数组。
type MemoResult struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Images    []string  `json:"images"`
	Weather   string    `json:"weather"`
	Mood      string    `json:"mood"`
	Pinned    bool      `json:"pinned"`
	Public    bool      `json:"public"`
	Tags      []string  `json:"tags"`
	SortOrder int       `json:"sortOrder"`
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func memoResult(memo blogmodel.Memo) MemoResult {
	tags := make([]string, 0)
	if len(memo.Tags) > 0 {
		_ = json.Unmarshal(memo.Tags, &tags)
		if tags == nil {
			tags = make([]string, 0)
		}
	}
	images := make([]string, 0)
	if len(memo.Images) > 0 {
		_ = json.Unmarshal(memo.Images, &images)
		if images == nil {
			images = make([]string, 0)
		}
	}
	return MemoResult{
		ID: memo.ID, Title: memo.Title, Content: memo.Content, Images: images,
		Weather: memo.Weather, Mood: memo.Mood, Pinned: memo.Pinned,
		Public: memo.Public, Tags: tags, SortOrder: memo.SortOrder,
		Version: memo.Version, CreatedAt: memo.CreatedAt, UpdatedAt: memo.UpdatedAt,
	}
}

// ListMemo 返回博主全部速记（含私有）：支持关键词搜索，排序 pinned DESC + updated_at DESC。
func (s *Service) ListMemo(ctx context.Context, query request.List) (Page, error) {
	db, err := s.database(ctx)
	if err != nil {
		return Page{}, err
	}
	q := db.Model(&blogmodel.Memo{})
	if query.Keyword != "" {
		pattern := "%" + strings.ToLower(query.Keyword) + "%"
		q = q.Where("LOWER(title) LIKE ? OR LOWER(content) LIKE ?", pattern, pattern)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return Page{}, err
	}
	var memos []blogmodel.Memo
	if err := q.Order("pinned DESC").Order("updated_at DESC, id ASC").Offset(query.Offset()).Limit(query.PageSize).Find(&memos).Error; err != nil {
		return Page{}, err
	}
	list := make([]MemoResult, 0, len(memos))
	for _, memo := range memos {
		list = append(list, memoResult(memo))
	}
	return Page{List: list, Pagination: Pagination{Page: query.Page, PageSize: query.PageSize, Total: total, OrderBy: "pinned_updated_at_desc"}}, nil
}

// GetMemo 返回单条速记详情（含私有）。
func (s *Service) GetMemo(ctx context.Context, id string) (interface{}, error) {
	db, err := s.database(ctx)
	if err != nil {
		return nil, err
	}
	var memo blogmodel.Memo
	if err := db.Where("id = ?", id).First(&memo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return memoResult(memo), nil
}

func validateMemo(payload request.Memo) error {
	if strings.TrimSpace(payload.Content) == "" {
		return request.Field("content", "正文不能为空")
	}
	if _, err := request.NormalizeIDs(payload.Tags); err != nil {
		return request.Field("tags", "标签格式不合法")
	}
	if _, err := request.NormalizeIDs(payload.Images); err != nil {
		return request.Field("images", "图片路径不合法")
	}
	if len(payload.Images) > 8 {
		return request.Field("images", "图片最多 8 张")
	}
	return nil
}

func memoFromPayload(payload request.Memo) blogmodel.Memo {
	tagsJSON := blogmodel.NewJSON(payload.Tags)
	imagesJSON := blogmodel.NewJSON(payload.Images)
	return blogmodel.Memo{
		Title: strings.TrimSpace(payload.Title), Content: payload.Content, Images: imagesJSON,
		Weather: strings.TrimSpace(payload.Weather), Mood: strings.TrimSpace(payload.Mood),
		Pinned: payload.Pinned, Public: payload.Public, Tags: tagsJSON, SortOrder: payload.SortOrder,
	}
}

// syncPublicAt 维护 public_at：首次公开时写入当前时间；取消公开则清零。
func syncPublicAt(public bool, publicAt time.Time) time.Time {
	if public {
		if publicAt.IsZero() {
			return time.Now()
		}
		return publicAt
	}
	return time.Time{}
}

// CreateMemo 新建速记，默认可见性由入参 public 决定（前端默认私有）。
func (s *Service) CreateMemo(ctx context.Context, payload request.Memo) (interface{}, error) {
	if err := validateMemo(payload); err != nil {
		return nil, err
	}
	db, err := s.database(ctx)
	if err != nil {
		return nil, err
	}
	memo := memoFromPayload(payload)
	memo.PublicAt = syncPublicAt(payload.Public, time.Time{})
	if err := db.Create(&memo).Error; err != nil {
		return nil, mapDBError(err)
	}
	return memoResult(memo), nil
}

// UpdateMemo 更新速记，要求携带当前 version（乐观锁）。
func (s *Service) UpdateMemo(ctx context.Context, id string, payload request.Memo) (interface{}, error) {
	if payload.Version < 1 {
		return nil, request.Field("version", "缺少 version：请先读取最新数据再提交（乐观锁）")
	}
	if err := validateMemo(payload); err != nil {
		return nil, err
	}
	db, err := s.database(ctx)
	if err != nil {
		return nil, err
	}
	var memo blogmodel.Memo
	if err := db.Where("id = ? AND version = ?", id, payload.Version).First(&memo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrConflict
		}
		return nil, err
	}
	memo.Title = strings.TrimSpace(payload.Title)
	memo.Content = payload.Content
	memo.Images = blogmodel.NewJSON(payload.Images)
	memo.Weather = strings.TrimSpace(payload.Weather)
	memo.Mood = strings.TrimSpace(payload.Mood)
	memo.Pinned = payload.Pinned
	memo.Public = payload.Public
	memo.PublicAt = syncPublicAt(payload.Public, memo.PublicAt)
	memo.Tags = blogmodel.NewJSON(payload.Tags)
	memo.SortOrder = payload.SortOrder
	memo.Version++
	if err := db.Model(&blogmodel.Memo{}).Where("id = ?", id).Updates(map[string]interface{}{
		"title": memo.Title, "content": memo.Content, "images": memo.Images,
		"weather": memo.Weather, "mood": memo.Mood, "pinned": memo.Pinned,
		"public": memo.Public, "public_at": memo.PublicAt, "tags": memo.Tags,
		"sort_order": memo.SortOrder, "version": memo.Version,
	}).Error; err != nil {
		return nil, mapDBError(err)
	}
	return memoResult(memo), nil
}

// UpdateMemoFields 快捷改单条速记字段（白名单：pinned / public），成功 version 自增。
func (s *Service) UpdateMemoFields(ctx context.Context, id string, fields map[string]interface{}) (interface{}, error) {
	updates := map[string]interface{}{}
	if value, ok := fields["pinned"]; ok {
		pinned, ok := value.(bool)
		if !ok {
			return nil, request.Field("pinned", "置顶取值必须是布尔值")
		}
		updates["pinned"] = pinned
	}
	if value, ok := fields["public"]; ok {
		public, ok := value.(bool)
		if !ok {
			return nil, request.Field("public", "可见性取值必须是布尔值")
		}
		updates["public"] = public
	}
	if len(updates) == 0 || len(updates) != len(fields) {
		return nil, request.Field("fields", "包含不支持的字段，快捷修改只支持 pinned 与 public")
	}
	db, err := s.database(ctx)
	if err != nil {
		return nil, err
	}
	var memo blogmodel.Memo
	if err := db.Where("id = ?", id).First(&memo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	memo.Pinned, _ = updates["pinned"].(bool)
	if _, ok := updates["public"]; ok {
		memo.Public, _ = updates["public"].(bool)
		memo.PublicAt = syncPublicAt(memo.Public, memo.PublicAt)
		updates["public_at"] = memo.PublicAt
	}
	updates["version"] = gorm.Expr("version + 1")
	if err := db.Model(&blogmodel.Memo{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return nil, mapDBError(err)
	}
	loaded, err := s.GetMemo(ctx, id)
	if err != nil {
		return nil, err
	}
	return loaded, nil
}

// DeleteMemo 软删除单条速记（不提供恢复入口，保持软删除以便数据保全）。
func (s *Service) DeleteMemo(ctx context.Context, id string) error {
	db, err := s.database(ctx)
	if err != nil {
		return err
	}
	result := db.Where("id = ?", id).Delete(&blogmodel.Memo{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// DeleteMemoBatch 批量软删除速记，返回实际删除条数。
func (s *Service) DeleteMemoBatch(ctx context.Context, ids []string) (int64, error) {
	db, err := s.database(ctx)
	if err != nil {
		return 0, err
	}
	result := db.Where("id IN ?", ids).Delete(&blogmodel.Memo{})
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
