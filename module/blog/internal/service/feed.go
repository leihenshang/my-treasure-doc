package service

import (
	"context"

	"fastduck/treasure-doc/module/blog/data/model"
	"fastduck/treasure-doc/module/blog/data/response"
)

// 站点固定页面，始终出现在 sitemap 中。
var staticPages = []string{"/", "/Blog", "/Blog/Diary", "/Blog/Portfolio", "/Blog/Tools", "/Blog/Bookmark", "/Blog/About"}

// Sitemap 汇总所有公开可访问的页面，供 sitemap.xml 使用。
func (s *Service) Sitemap(ctx context.Context) ([]response.SitemapEntry, error) {
	db, err := s.database(ctx)
	if err != nil {
		return nil, err
	}

	entries := make([]response.SitemapEntry, 0, len(staticPages)+32)
	for _, path := range staticPages {
		entries = append(entries, response.SitemapEntry{Loc: path})
	}

	var posts []model.Post
	if err := published(db.Model(&model.Post{})).Select("slug", "updated_at").Find(&posts).Error; err != nil {
		return nil, err
	}
	for _, post := range posts {
		entries = append(entries, response.SitemapEntry{Loc: "/Blog/" + post.Slug, LastMod: post.UpdatedAt})
	}

	var diaries []model.Diary
	if err := published(db.Model(&model.Diary{})).Select("public_id", "updated_at").Find(&diaries).Error; err != nil {
		return nil, err
	}
	for _, diary := range diaries {
		entries = append(entries, response.SitemapEntry{Loc: "/Blog/Diary/" + diary.PublicID, LastMod: diary.UpdatedAt})
	}

	var works []model.PortfolioItem
	if err := published(db.Model(&model.PortfolioItem{})).Select("slug", "updated_at").Find(&works).Error; err != nil {
		return nil, err
	}
	for _, work := range works {
		entries = append(entries, response.SitemapEntry{Loc: "/Blog/Portfolio/" + work.Slug, LastMod: work.UpdatedAt})
	}

	var tools []model.Tool
	if err := published(db.Model(&model.Tool{})).Where("kind = ?", "own").Select("slug", "updated_at").Find(&tools).Error; err != nil {
		return nil, err
	}
	for _, tool := range tools {
		entries = append(entries, response.SitemapEntry{Loc: "/Blog/Tools/" + tool.Slug, LastMod: tool.UpdatedAt})
	}

	return entries, nil
}

// Feed 返回最新发布的文章（含正文），供 RSS 使用。
func (s *Service) Feed(ctx context.Context, limit int) ([]response.Post, error) {
	db, err := s.database(ctx)
	if err != nil {
		return nil, err
	}
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	var posts []model.Post
	if err := published(db.Model(&model.Post{})).Order("published_on DESC, id ASC").Limit(limit).Find(&posts).Error; err != nil {
		return nil, err
	}

	items := make([]response.Post, 0, len(posts))
	for _, post := range posts {
		tags, err := s.tagsFor(ctx, "td_blog_post_tag", "post_id", post.ID)
		if err != nil {
			return nil, err
		}
		items = append(items, response.Post{PostSummary: postSummary(post, tags), Content: post.Content})
	}
	return items, nil
}
