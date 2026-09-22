package service

import (
	"context"
	"errors"
	"fmt"

	"fastduck/treasure-doc/module/blog/data/model"
	"fastduck/treasure-doc/module/blog_mgr/data/request"

	"gorm.io/gorm"
)

// PublishCategory / PublishTag 供发布端（令牌）接口返回的精简分类/标签结构。
type PublishCategory struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Scope string `json:"scope"`
}

type PublishTag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ListCategoriesForPublish 返回发布端可选的分类（可按 scope 过滤，scope 空=全部）。
func (s *Service) ListCategoriesForPublish(ctx context.Context, scope string) ([]PublishCategory, error) {
	db, err := s.database(ctx)
	if err != nil {
		return nil, err
	}
	q := db.Model(&model.Category{}).Order("scope, sort_order, name")
	if scope != "" {
		q = q.Where("scope = ?", scope)
	}
	var rows []model.Category
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]PublishCategory, 0, len(rows))
	for _, row := range rows {
		out = append(out, PublishCategory{ID: row.ID, Name: row.Name, Slug: row.Slug, Scope: row.Scope})
	}
	return out, nil
}

// ListTagsForPublish 返回发布端可选的标签。
func (s *Service) ListTagsForPublish(ctx context.Context) ([]PublishTag, error) {
	db, err := s.database(ctx)
	if err != nil {
		return nil, err
	}
	var rows []model.Tag
	if err := db.Model(&model.Tag{}).Order("name").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]PublishTag, 0, len(rows))
	for _, row := range rows {
		out = append(out, PublishTag{ID: row.ID, Name: row.Name})
	}
	return out, nil
}

// publishLookup 校验发布支持的资源（文章/日记），这里仅做白名单判断。
func isPublishableContent(resource string) bool {
	return resource == "posts" || resource == "diaries"
}

// PublishStatus 发布状态查询结果（用于插件判断是否已发布、能否覆盖）。
type PublishStatus struct {
	Exists bool   `json:"exists"`
	ID     string `json:"id"`
	Title  string `json:"title"`
	Slug   string `json:"slug"`
}

// PublishLookup 按 slug 查询文章/日记是否已发布（deleted_at 为空）。
func (s *Service) PublishLookup(ctx context.Context, resource, slug string) (PublishStatus, error) {
	status := PublishStatus{}
	if !isPublishableContent(resource) || !request.ValidID(slug) {
		return status, ErrInvalid
	}
	db, err := s.database(ctx)
	if err != nil {
		return status, err
	}
	var table, field string
	switch resource {
	case "posts":
		table, field = "td_blog_post", "slug"
	case "diaries":
		table, field = "td_blog_diary", "public_id"
	}
	type row struct {
		ID    string
		Title string
	}
	var found row
	if err := db.Table(table).Where(field+" = ? AND deleted_at IS NULL", slug).First(&found).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			status.Slug = slug
			return status, nil
		}
		return status, err
	}
	status.Exists = true
	status.ID = found.ID
	status.Title = found.Title
	status.Slug = slug
	return status, nil
}

// PublishCheck 发布前校验：按标题推导 slug（含 -2/-3… 冲突后缀回查）判断该资源是否已发布，
// 供插件在勾选「覆盖」前调用，无需插件本地缓存。已发布返回真实 slug；未发布返回生成的基础 slug。
func (s *Service) PublishCheck(ctx context.Context, resource, title string) (PublishStatus, error) {
	status := PublishStatus{}
	if !isPublishableContent(resource) {
		return status, ErrInvalid
	}
	base := slugify(title)
	if base == "" {
		return status, ErrInvalid
	}
	db, err := s.database(ctx)
	if err != nil {
		return status, err
	}
	var table, field string
	switch resource {
	case "posts":
		table, field = "td_blog_post", "slug"
	case "diaries":
		table, field = "td_blog_diary", "public_id"
	}
	type row struct {
		ID    string
		Title string
	}
	var found row
	for n := 0; n < 9; n++ {
		candidate := base
		if n > 0 {
			candidate = fmt.Sprintf("%s-%d", base, n+1)
		}
		if err := db.Table(table).Where(field+" = ? AND deleted_at IS NULL", candidate).First(&found).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}
			return status, err
		}
		status.Exists = true
		status.ID = found.ID
		status.Title = found.Title
		status.Slug = candidate
		return status, nil
	}
	status.Slug = base
	return status, nil
}

// PublishUpdate 按 slug 强制覆盖文章/日记（发布端覆盖）：不要求乐观锁版本，保留创建时间/浏览量，
// 强制 published，version+1，并记录历史快照。仅 posts/diaries。
func (s *Service) PublishUpdate(ctx context.Context, resource, slug string, payload interface{}) (interface{}, error) {
	if !isPublishableContent(resource) || !request.ValidID(slug) {
		return nil, ErrInvalid
	}
	db, err := s.database(ctx)
	if err != nil {
		return nil, err
	}
	item, tagIDs, relation, err := buildModel(resource, payload)
	if err != nil {
		return nil, err
	}
	// slug 以路径为准，避免覆盖时静默改名
	if v, ok := item.(*model.Post); ok {
		v.Slug = slug
	} else if v, ok := item.(*model.Diary); ok {
		v.PublicID = slug
	}

	var result interface{}
	err = db.Transaction(func(tx *gorm.DB) error {
		// 定位既有记录
		id, prevVersion, err := resolveByPublishKey(tx, resource, slug)
		if err != nil {
			return err
		}
		if err := validateReferences(tx, resource, item, tagIDs); err != nil {
			return err
		}
		if err := setModelOwnership(item, id, prevVersion); err != nil {
			return err
		}
		updates := updateMap(item)
		updates["version"] = gorm.Expr("version + 1")
		r := tx.Model(item).Where("id = ?", id).Updates(updates)
		if r.Error != nil {
			return mapDBError(r.Error)
		}
		if relation != "" {
			if err := replaceTags(tx, relation, id, tagIDs); err != nil {
				return err
			}
		}
		loaded, err := s.getWithDB(tx, resource, id)
		if err != nil {
			return err
		}
		result, err = enrichItemWithTags(tx, resource, loaded)
		if err != nil {
			return err
		}
		return recordEditHistory(tx, resource, result)
	})
	return result, err
}

// resolveByPublishKey 按发布主键（post:slug / diary:public_id）定位记录的 id 与当前 version。
func resolveByPublishKey(tx *gorm.DB, resource, slug string) (string, int, error) {
	var table, field string
	switch resource {
	case "posts":
		table, field = "td_blog_post", "slug"
	case "diaries":
		table, field = "td_blog_diary", "public_id"
	}
	var row struct {
		ID      string
		Version int
	}
	if err := tx.Table(table).Where(field+" = ? AND deleted_at IS NULL", slug).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", 0, ErrNotFound
		}
		return "", 0, err
	}
	return row.ID, row.Version, nil
}

// setModelOwnership 设置模型的自有/乐观锁字段（保留 ID、CreatedAt、ViewCount，version 自增由 SQL 表达式处理）。
func setModelOwnership(item interface{}, id string, prevVersion int) error {
	switch v := item.(type) {
	case *model.Post:
		v.ID = id
		v.Version = prevVersion + 1
	case *model.Diary:
		v.ID = id
		v.Version = prevVersion + 1
	}
	return nil
}
