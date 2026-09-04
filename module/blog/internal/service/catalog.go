package service

import (
	"context"
	"errors"
	"time"

	"fastduck/treasure-doc/module/blog/data/model"
	"fastduck/treasure-doc/module/blog/data/request"
	"fastduck/treasure-doc/module/blog/data/response"

	"gorm.io/gorm"
)

func (s *Service) ListPortfolio(ctx context.Context, query request.PortfolioQuery) ([]response.PortfolioSummary, error) {
	db, err := s.database(ctx)
	if err != nil {
		return nil, err
	}
	q := published(db.Model(&model.PortfolioItem{}))
	if query.CategoryID != "" {
		q = q.Where("category_id = ?", query.CategoryID)
	}
	var records []model.PortfolioItem
	if err := q.Order("published_on DESC, slug ASC").Find(&records).Error; err != nil {
		return nil, err
	}
	items := make([]response.PortfolioSummary, 0, len(records))
	for _, record := range records {
		techStack, err := decodeJSON[string](record.TechStack)
		if err != nil {
			return nil, err
		}
		items = append(items, portfolioSummary(record, techStack))
	}
	return items, nil
}

func (s *Service) GetPortfolio(ctx context.Context, id string) (response.PortfolioItem, error) {
	db, err := s.database(ctx)
	if err != nil {
		return response.PortfolioItem{}, err
	}
	var record model.PortfolioItem
	if err := published(db).Where("slug = ?", id).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.PortfolioItem{}, ErrPortfolioNotFound
		}
		return response.PortfolioItem{}, err
	}
	techStack, err := decodeJSON[string](record.TechStack)
	if err != nil {
		return response.PortfolioItem{}, err
	}
	links, err := decodeJSON[response.PortfolioLink](record.Links)
	if err != nil {
		return response.PortfolioItem{}, err
	}
	return response.PortfolioItem{PortfolioSummary: portfolioSummary(record, techStack), Links: links, Content: record.Content}, nil
}

func portfolioSummary(record model.PortfolioItem, techStack []string) response.PortfolioSummary {
	return response.PortfolioSummary{ID: record.Slug, Title: record.Title, Summary: record.Summary, Category: record.CategoryID, Cover: record.Cover, TechStack: techStack, Date: record.PublishedOn.Format("2006-01-02")}
}

func (s *Service) ListTools(ctx context.Context) ([]response.Tool, error) {
	db, err := s.database(ctx)
	if err != nil {
		return nil, err
	}
	var records []model.Tool
	if err := published(db.Model(&model.Tool{})).Order("sort_order ASC, slug ASC").Find(&records).Error; err != nil {
		return nil, err
	}
	items := make([]response.Tool, 0, len(records))
	for _, record := range records {
		items = append(items, toolResponse(record))
	}
	return items, nil
}

func (s *Service) GetTool(ctx context.Context, id string) (response.Tool, error) {
	db, err := s.database(ctx)
	if err != nil {
		return response.Tool{}, err
	}
	var record model.Tool
	if err := published(db).Where("slug = ? AND kind = ?", id, "own").First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Tool{}, ErrToolNotFound
		}
		return response.Tool{}, err
	}
	return toolResponse(record), nil
}

func toolResponse(record model.Tool) response.Tool {
	item := response.Tool{ID: record.Slug, Type: record.Kind, Name: record.Name, Desc: record.Description}
	if record.Kind == "link" {
		item.URL = record.URL
	} else {
		item.Cover = record.Cover
		item.Status = record.DevelopmentStatus
		item.Content = record.Content
	}
	return item
}

func (s *Service) ListBookmarks(ctx context.Context, query request.BookmarkQuery) ([]response.Bookmark, error) {
	db, err := s.database(ctx)
	if err != nil {
		return nil, err
	}
	q := published(db.Model(&model.Bookmark{}))
	if query.CategoryID != "" {
		q = q.Where("category_id = ?", query.CategoryID)
	}
	if query.Keyword != "" {
		pattern := likePattern(query.Keyword)
		q = q.Where("LOWER(title) LIKE ? ESCAPE '\\' OR LOWER(description) LIKE ? ESCAPE '\\' OR LOWER(url) LIKE ? ESCAPE '\\' OR EXISTS (SELECT 1 FROM td_blog_bookmark_tag bt JOIN td_blog_tag t ON t.id = bt.tag_id WHERE bt.bookmark_id = td_blog_bookmark.id AND LOWER(t.name) LIKE ? ESCAPE '\\' AND t.deleted_at IS NULL)", pattern, pattern, pattern, pattern)
	}
	var records []model.Bookmark
	if err := q.Order("sort_order ASC, public_id ASC").Find(&records).Error; err != nil {
		return nil, err
	}
	items := make([]response.Bookmark, 0, len(records))
	for _, record := range records {
		tags, err := s.tagsFor(ctx, "td_blog_bookmark_tag", "bookmark_id", record.ID)
		if err != nil {
			return nil, err
		}
		items = append(items, response.Bookmark{ID: record.PublicID, Title: record.Title, URL: record.URL, Desc: record.Description, Category: record.CategoryID, Tags: tags, Icon: record.Icon})
	}
	return items, nil
}

func (s *Service) Profile(ctx context.Context) (response.Profile, error) {
	db, err := s.database(ctx)
	if err != nil {
		return response.Profile{}, err
	}
	var record model.Profile
	if err := db.Where("profile_key = ?", "default").First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Profile{Links: []response.ProfileLink{}, Skills: []response.ProfileSkill{}}, nil
		}
		return response.Profile{}, err
	}
	links, err := decodeJSON[response.ProfileLink](record.Links)
	if err != nil {
		return response.Profile{}, err
	}
	skills, err := decodeJSON[response.ProfileSkill](record.Skills)
	if err != nil {
		return response.Profile{}, err
	}
	return response.Profile{Name: record.Name, Avatar: record.Avatar, Role: record.Role, Location: record.Location, Motto: record.Motto, Bio: record.Bio, Links: links, Skills: skills}, nil
}

func (s *Service) Site(ctx context.Context) (response.Site, error) {
	db, err := s.database(ctx)
	if err != nil {
		return response.Site{}, err
	}
	var record model.Site
	if err := db.Where("site_key = ?", "default").First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Site{TechStack: []string{}, Modules: []response.SiteModule{}, Milestones: []response.SiteMilestone{}}, nil
		}
		return response.Site{}, err
	}
	techStack, err := decodeJSON[string](record.TechStack)
	if err != nil {
		return response.Site{}, err
	}
	modules, err := decodeJSON[response.SiteModule](record.Modules)
	if err != nil {
		return response.Site{}, err
	}
	milestones, err := decodeJSON[response.SiteMilestone](record.Milestones)
	if err != nil {
		return response.Site{}, err
	}
	return response.Site{Name: record.Name, Slogan: record.Slogan, Intro: record.Intro, TechStack: techStack, Modules: modules, Milestones: milestones}, nil
}

func (s *Service) Stats(ctx context.Context) (response.Stats, error) {
	db, err := s.database(ctx)
	if err != nil {
		return response.Stats{}, err
	}
	now := time.Now()
	counts := []struct {
		model interface{}
		value *int64
	}{
		{model: &model.Post{}},
		{model: &model.Diary{}},
		{model: &model.PortfolioItem{}},
	}
	var result response.Stats
	counts[0].value = &result.Posts
	counts[1].value = &result.Diaries
	counts[2].value = &result.Works
	for _, count := range counts {
		if err := db.Model(count.model).Where("publish_status = ? AND published_at <= ?", model.StatusPublished, now).Count(count.value).Error; err != nil {
			return response.Stats{}, err
		}
	}
	return result, nil
}
