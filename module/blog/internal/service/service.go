package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"fastduck/treasure-doc/module/blog/data/model"
	"fastduck/treasure-doc/module/blog/data/request"
	"fastduck/treasure-doc/module/blog/data/response"

	"gorm.io/gorm"
)

var (
	ErrPostNotFound      = errors.New("post not found")
	ErrDiaryNotFound     = errors.New("diary not found")
	ErrPortfolioNotFound = errors.New("portfolio item not found")
	ErrToolNotFound      = errors.New("tool not found")
)

type DBProvider func() *gorm.DB

type Service struct {
	db DBProvider
}

func New(db DBProvider) *Service {
	return &Service{db: db}
}

func (s *Service) database(ctx context.Context) (*gorm.DB, error) {
	if s.db == nil || s.db() == nil {
		return nil, errors.New("database is not initialized")
	}
	return s.db().WithContext(ctx), nil
}

func published(db *gorm.DB) *gorm.DB {
	return db.Where("publish_status = ? AND published_at <= ?", model.StatusPublished, time.Now())
}

func orderByDate(sort string) string {
	if sort == "asc" {
		return "published_on ASC, id ASC"
	}
	return "published_on DESC, id ASC"
}

func likePattern(keyword string) string {
	return "%" + request.EscapeLike(strings.ToLower(keyword)) + "%"
}

func (s *Service) Categories(ctx context.Context, scope string) ([]response.Category, error) {
	db, err := s.database(ctx)
	if err != nil {
		return nil, err
	}

	var categories []model.Category
	query := db.Where("scope = ? AND enabled = ?", scope, true)
	switch scope {
	case model.CategoryPost:
		query = query.Where("EXISTS (?)", published(db.Model(&model.Post{})).Select("1").Where("category_id = td_blog_category.slug"))
	case model.CategoryPortfolio:
		query = query.Where("EXISTS (?)", published(db.Model(&model.PortfolioItem{})).Select("1").Where("category_id = td_blog_category.slug"))
	case model.CategoryBookmark:
		query = query.Where("EXISTS (?)", published(db.Model(&model.Bookmark{})).Select("1").Where("category_id = td_blog_category.slug"))
	}
	if err := query.Order("sort_order ASC, slug ASC").Find(&categories).Error; err != nil {
		return nil, err
	}

	result := make([]response.Category, 0, len(categories))
	for _, category := range categories {
		result = append(result, response.Category{ID: category.Slug, Name: category.Name})
	}
	return result, nil
}

func (s *Service) PostTags(ctx context.Context) ([]string, error) {
	return s.publishedTags(ctx, "td_blog_post_tag", "post_id", "td_blog_post")
}

func (s *Service) DiaryTags(ctx context.Context) ([]string, error) {
	return s.publishedTags(ctx, "td_blog_diary_tag", "diary_id", "td_blog_diary")
}

func (s *Service) publishedTags(ctx context.Context, relationTable, ownerColumn, ownerTable string) ([]string, error) {
	db, err := s.database(ctx)
	if err != nil {
		return nil, err
	}

	tags := make([]string, 0)
	err = db.Table("td_blog_tag AS tag").Distinct("tag.name").
		Joins(fmt.Sprintf("JOIN %s AS relation ON relation.tag_id = tag.id", relationTable)).
		Joins(fmt.Sprintf("JOIN %s AS owner ON owner.id = relation.%s", ownerTable, ownerColumn)).
		Where("owner.publish_status = ? AND owner.published_at <= ? AND owner.deleted_at IS NULL AND tag.deleted_at IS NULL", model.StatusPublished, time.Now()).
		Order("tag.name ASC").Pluck("tag.name", &tags).Error
	return tags, err
}

func (s *Service) ListPosts(ctx context.Context, query request.PostQuery) (response.Page, error) {
	db, err := s.database(ctx)
	if err != nil {
		return response.Page{}, err
	}

	q := published(db.Model(&model.Post{}))
	if query.CategoryID != "" {
		q = q.Where("category_id = ?", query.CategoryID)
	}
	if query.Tag != "" {
		q = q.Where("EXISTS (SELECT 1 FROM td_blog_post_tag pt JOIN td_blog_tag t ON t.id = pt.tag_id WHERE pt.post_id = td_blog_post.id AND t.name = ? AND t.deleted_at IS NULL)", query.Tag)
	}
	if query.Keyword != "" {
		pattern := likePattern(query.Keyword)
		q = q.Where("LOWER(title) LIKE ? ESCAPE '\\' OR LOWER(summary) LIKE ? ESCAPE '\\' OR LOWER(content) LIKE ? ESCAPE '\\' OR EXISTS (SELECT 1 FROM td_blog_post_tag pt JOIN td_blog_tag t ON t.id = pt.tag_id WHERE pt.post_id = td_blog_post.id AND LOWER(t.name) LIKE ? ESCAPE '\\' AND t.deleted_at IS NULL)", pattern, pattern, pattern, pattern)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return response.Page{}, err
	}
	var posts []model.Post
	if err := q.Order("pinned DESC").Order(orderByDate(query.Sort)).Offset(query.Offset()).Limit(query.PageSize).Find(&posts).Error; err != nil {
		return response.Page{}, err
	}

	items := make([]response.PostSummary, 0, len(posts))
	for _, post := range posts {
		tags, err := s.tagsFor(ctx, "td_blog_post_tag", "post_id", post.ID)
		if err != nil {
			return response.Page{}, err
		}
		items = append(items, postSummary(post, tags))
	}
	return response.Page{List: items, Pagination: response.Pagination{Page: query.Page, PageSize: query.PageSize, Total: total, OrderBy: "date_" + query.Sort}}, nil
}

func (s *Service) GetPost(ctx context.Context, id string) (response.Post, error) {
	db, err := s.database(ctx)
	if err != nil {
		return response.Post{}, err
	}
	var post model.Post
	if err := published(db).Where("slug = ?", id).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Post{}, ErrPostNotFound
		}
		return response.Post{}, err
	}
	tags, err := s.tagsFor(ctx, "td_blog_post_tag", "post_id", post.ID)
	if err != nil {
		return response.Post{}, err
	}
	return response.Post{PostSummary: postSummary(post, tags), Content: post.Content}, nil
}

func postSummary(post model.Post, tags []string) response.PostSummary {
	return response.PostSummary{ID: post.Slug, Title: post.Title, Summary: post.Summary, Category: post.CategoryID, Tags: tags, Author: post.Author, Date: post.PublishedOn.Format("2006-01-02"), Pinned: post.Pinned}
}

func (s *Service) ListDiaries(ctx context.Context, query request.DiaryQuery) (response.Page, error) {
	db, err := s.database(ctx)
	if err != nil {
		return response.Page{}, err
	}
	q := published(db.Model(&model.Diary{}))
	if query.Tag != "" {
		q = q.Where("EXISTS (SELECT 1 FROM td_blog_diary_tag dt JOIN td_blog_tag t ON t.id = dt.tag_id WHERE dt.diary_id = td_blog_diary.id AND t.name = ? AND t.deleted_at IS NULL)", query.Tag)
	}
	if query.Keyword != "" {
		pattern := likePattern(query.Keyword)
		q = q.Where("LOWER(title) LIKE ? ESCAPE '\\' OR LOWER(summary) LIKE ? ESCAPE '\\' OR LOWER(content) LIKE ? ESCAPE '\\' OR EXISTS (SELECT 1 FROM td_blog_diary_tag dt JOIN td_blog_tag t ON t.id = dt.tag_id WHERE dt.diary_id = td_blog_diary.id AND LOWER(t.name) LIKE ? ESCAPE '\\' AND t.deleted_at IS NULL)", pattern, pattern, pattern, pattern)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return response.Page{}, err
	}
	var diaries []model.Diary
	if err := q.Order("pinned DESC").Order(orderByDate(query.Sort)).Offset(query.Offset()).Limit(query.PageSize).Find(&diaries).Error; err != nil {
		return response.Page{}, err
	}
	items := make([]response.DiarySummary, 0, len(diaries))
	for _, diary := range diaries {
		tags, err := s.tagsFor(ctx, "td_blog_diary_tag", "diary_id", diary.ID)
		if err != nil {
			return response.Page{}, err
		}
		items = append(items, diarySummary(diary, tags))
	}
	return response.Page{List: items, Pagination: response.Pagination{Page: query.Page, PageSize: query.PageSize, Total: total, OrderBy: "date_" + query.Sort}}, nil
}

func (s *Service) GetDiary(ctx context.Context, id string) (response.Diary, error) {
	db, err := s.database(ctx)
	if err != nil {
		return response.Diary{}, err
	}
	var diary model.Diary
	if err := published(db).Where("public_id = ?", id).First(&diary).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Diary{}, ErrDiaryNotFound
		}
		return response.Diary{}, err
	}
	tags, err := s.tagsFor(ctx, "td_blog_diary_tag", "diary_id", diary.ID)
	if err != nil {
		return response.Diary{}, err
	}
	return response.Diary{DiarySummary: diarySummary(diary, tags), Content: diary.Content}, nil
}

func diarySummary(diary model.Diary, tags []string) response.DiarySummary {
	return response.DiarySummary{ID: diary.PublicID, Title: diary.Title, Summary: diary.Summary, Tags: tags, Date: diary.PublishedOn.Format("2006-01-02"), Mood: diary.Mood, Weather: diary.Weather, Pinned: diary.Pinned}
}

func (s *Service) tagsFor(ctx context.Context, relationTable, ownerColumn, ownerID string) ([]string, error) {
	db, err := s.database(ctx)
	if err != nil {
		return nil, err
	}
	tags := make([]string, 0)
	err = db.Table("td_blog_tag AS tag").Joins(fmt.Sprintf("JOIN %s AS relation ON relation.tag_id = tag.id", relationTable)).
		Where("relation."+ownerColumn+" = ? AND tag.deleted_at IS NULL", ownerID).Order("tag.name ASC").Pluck("tag.name", &tags).Error
	return tags, err
}

func decodeJSON[T any](value model.JSON) ([]T, error) {
	result := make([]T, 0)
	if len(value) == 0 {
		return result, nil
	}
	if err := json.Unmarshal(value, &result); err != nil {
		return nil, err
	}
	return result, nil
}
