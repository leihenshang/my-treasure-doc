package service

import (
	"context"
	"encoding/json"
	"errors"

	"fastduck/treasure-doc/module/blog/data/model"
	"fastduck/treasure-doc/module/blog_mgr/data/request"

	"gorm.io/gorm"
)

// maxEditHistory 每个 ref（文章/日记）保留的最新历史条数，超过删除最旧。
const maxEditHistory = 30

// EditHistoryMeta 历史版本列表项。
type EditHistoryMeta struct {
	Seq       int    `json:"seq"`
	Version   int    `json:"version"`
	CreatedAt string `json:"createdAt"`
	Summary   string `json:"summary"`
}

// EditHistoryDetail 历史版本详情（含完整快照）。
type EditHistoryDetail struct {
	Seq       int         `json:"seq"`
	Version   int         `json:"version"`
	CreatedAt string      `json:"createdAt"`
	Summary   string      `json:"summary"`
	Snapshot  interface{} `json:"snapshot"`
}

// recordEditHistory 为文章/日记的每次保存追加一条历史快照（仅 posts/diaries 生效），并裁剪到最近 N 条。
func recordEditHistory(tx *gorm.DB, resource string, snapshot interface{}) error {
	refID := historyRefID(snapshot)
	if refID == "" {
		return nil
	}
	summary, version := historySummaryAndVersion(snapshot)
	var maxSeq int64
	if err := tx.Model(&model.EditHistory{}).
		Where("resource = ? AND ref_id = ?", resource, refID).
		Select("COALESCE(MAX(seq), 0)").Scan(&maxSeq).Error; err != nil {
		return err
	}
	data, err := json.Marshal(snapshot)
	if err != nil {
		return err
	}
	rec := &model.EditHistory{Resource: resource, RefID: refID, Seq: int(maxSeq) + 1, Version: version, Summary: summary, Snapshot: data}
	if err := tx.Create(rec).Error; err != nil {
		return err
	}
	// 保留最近 maxEditHistory 条：删除排序更旧的
	var stale []string
	if err := tx.Model(&model.EditHistory{}).
		Where("resource = ? AND ref_id = ?", resource, refID).
		Order("seq DESC").Offset(maxEditHistory).Limit(int(maxSeq)+1).Pluck("id", &stale).Error; err == nil && len(stale) > 0 {
		tx.Where("id IN ?", stale).Delete(&model.EditHistory{})
	}
	return nil
}

func historyRefID(snapshot interface{}) string {
	switch v := snapshot.(type) {
	case PostWithTags:
		return v.ID
	case *PostWithTags:
		return v.ID
	case DiaryWithTags:
		return v.ID
	case *DiaryWithTags:
		return v.ID
	}
	return ""
}

func historySummaryAndVersion(snapshot interface{}) (string, int) {
	switch v := snapshot.(type) {
	case PostWithTags:
		return v.Title, v.Version
	case *PostWithTags:
		return v.Title, v.Version
	case DiaryWithTags:
		return v.Title, v.Version
	case *DiaryWithTags:
		return v.Title, v.Version
	}
	return "", 0
}

// ListEditHistory 列出某文章/日记的历史版本（倒序，不含完整快照）。
func (s *Service) ListEditHistory(ctx context.Context, resource, id string) ([]EditHistoryMeta, error) {
	if !isPublishableContent(resource) || !request.ValidID(id) {
		return nil, ErrInvalid
	}
	db, err := s.database(ctx)
	if err != nil {
		return nil, err
	}
	var rows []model.EditHistory
	if err := db.Model(&model.EditHistory{}).
		Where("resource = ? AND ref_id = ?", resource, id).
		Order("seq DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]EditHistoryMeta, 0, len(rows))
	for _, row := range rows {
		out = append(out, EditHistoryMeta{Seq: row.Seq, Version: row.Version, Summary: row.Summary, CreatedAt: row.CreatedAt.Format("2006-01-02 15:04:05")})
	}
	return out, nil
}

// GetEditHistory 读取某历史版本的完整快照。
func (s *Service) GetEditHistory(ctx context.Context, resource, id string, seq int) (EditHistoryDetail, error) {
	if !isPublishableContent(resource) || !request.ValidID(id) || seq < 1 {
		return EditHistoryDetail{}, ErrInvalid
	}
	db, err := s.database(ctx)
	if err != nil {
		return EditHistoryDetail{}, err
	}
	var row model.EditHistory
	if err := db.Model(&model.EditHistory{}).
		Where("resource = ? AND ref_id = ? AND seq = ?", resource, id, seq).
		First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return EditHistoryDetail{}, ErrNotFound
		}
		return EditHistoryDetail{}, err
	}
	var snapshot interface{}
	_ = json.Unmarshal(row.Snapshot, &snapshot)
	return EditHistoryDetail{Seq: row.Seq, Version: row.Version, Summary: row.Summary, CreatedAt: row.CreatedAt.Format("2006-01-02 15:04:05"), Snapshot: snapshot}, nil
}

// RestoreEditHistory 把某历史版本恢复为当前内容：写回字段、替换标签、version+1，并记为一条新历史。
func (s *Service) RestoreEditHistory(ctx context.Context, resource, id string, seq int) (interface{}, error) {
	if !isPublishableContent(resource) || !request.ValidID(id) || seq < 1 {
		return nil, ErrInvalid
	}
	db, err := s.database(ctx)
	if err != nil {
		return nil, err
	}
	var row model.EditHistory
	if err := db.Model(&model.EditHistory{}).
		Where("resource = ? AND ref_id = ? AND seq = ?", resource, id, seq).
		First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	var result interface{}
	err = db.Transaction(func(tx *gorm.DB) error {
		switch resource {
		case "posts":
			var snap PostWithTags
			if err := json.Unmarshal(row.Snapshot, &snap); err != nil {
				return ErrInvalid
			}
			updates := map[string]interface{}{
				"slug": snap.Slug, "title": snap.Title, "summary": snap.Summary, "category_id": snap.CategoryID,
				"author": snap.Author, "content": snap.Content, "publish_status": snap.PublishStatus,
				"published_on": snap.PublishedOn, "published_at": snap.PublishedAt, "pinned": snap.Pinned,
				"version": gorm.Expr("version + 1"),
			}
			r := tx.Model(&model.Post{}).Where("id = ?", id).Updates(updates)
			if r.Error != nil {
				return mapDBError(r.Error)
			}
			if err := replaceTags(tx, "td_blog_post_tag", id, snap.TagIDs); err != nil {
				return err
			}
		case "diaries":
			var snap DiaryWithTags
			if err := json.Unmarshal(row.Snapshot, &snap); err != nil {
				return ErrInvalid
			}
			updates := map[string]interface{}{
				"public_id": snap.PublicID, "title": snap.Title, "summary": snap.Summary, "content": snap.Content,
				"mood": snap.Mood, "weather": snap.Weather, "publish_status": snap.PublishStatus,
				"published_on": snap.PublishedOn, "published_at": snap.PublishedAt, "pinned": snap.Pinned,
				"version": gorm.Expr("version + 1"),
			}
			r := tx.Model(&model.Diary{}).Where("id = ?", id).Updates(updates)
			if r.Error != nil {
				return mapDBError(r.Error)
			}
			if err := replaceTags(tx, "td_blog_diary_tag", id, snap.TagIDs); err != nil {
				return err
			}
		default:
			return ErrInvalid
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
