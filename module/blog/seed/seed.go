package seed

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"fastduck/treasure-doc/module/blog/data/model"
	blogresponse "fastduck/treasure-doc/module/blog/data/response"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Seed(db *gorm.DB, options Options) (Result, error) {
	if err := Validate(options); err != nil || !options.Enabled {
		return Result{}, err
	}
	if db == nil {
		return Result{}, fmt.Errorf("database is not initialized")
	}
	var result Result
	err := db.Transaction(func(tx *gorm.DB) error {
		categories, err := seedCategories(tx, options, &result)
		if err != nil {
			return err
		}
		tags, err := seedTags(tx, options, &result)
		if err != nil {
			return err
		}
		if err = seedPosts(tx, options, tags, &result); err != nil {
			return err
		}
		if err = seedDiaries(tx, options, tags, &result); err != nil {
			return err
		}
		if err = seedPortfolio(tx, options, &result); err != nil {
			return err
		}
		if err = seedTools(tx, options, &result); err != nil {
			return err
		}
		if err = seedBookmarks(tx, options, tags, &result); err != nil {
			return err
		}
		_ = categories
		return seedSettings(tx, options, &result)
	})
	return result, err
}

func seedCategories(tx *gorm.DB, options Options, result *Result) (map[string]*model.Category, error) {
	definitions := []model.Category{
		{Scope: model.CategoryPost, Slug: "tech", Name: "技术", SortOrder: 10, Enabled: true}, {Scope: model.CategoryPost, Slug: "essay", Name: "随笔", SortOrder: 20, Enabled: true},
		{Scope: model.CategoryPost, Slug: "reading", Name: "阅读", SortOrder: 30, Enabled: true}, {Scope: model.CategoryPost, Slug: "unused", Name: "未使用", SortOrder: 40, Enabled: true},
		{Scope: model.CategoryPost, Slug: "hidden", Name: "隐藏分类", SortOrder: 50, Enabled: false}, {Scope: model.CategoryPortfolio, Slug: "website", Name: "网站", SortOrder: 10, Enabled: true},
		{Scope: model.CategoryPortfolio, Slug: "application", Name: "应用", SortOrder: 20, Enabled: true}, {Scope: model.CategoryPortfolio, Slug: "unused-work", Name: "未使用作品分类", SortOrder: 30, Enabled: true},
		{Scope: model.CategoryBookmark, Slug: "dev", Name: "开发", SortOrder: 10, Enabled: true}, {Scope: model.CategoryBookmark, Slug: "design", Name: "设计", SortOrder: 20, Enabled: true},
		{Scope: model.CategoryBookmark, Slug: "hidden-link", Name: "隐藏书签分类", SortOrder: 30, Enabled: false},
	}
	items := make(map[string]*model.Category, len(definitions))
	for index := range definitions {
		item := definitions[index]
		if err := ensure(tx, "scope = ? AND slug = ?", []interface{}{item.Scope, item.Slug}, &item, options, result); err != nil {
			return nil, err
		}
		items[item.Scope+":"+item.Slug] = &item
	}
	return items, nil
}

func seedTags(tx *gorm.DB, options Options, result *Result) (map[string]*model.Tag, error) {
	names := []string{"Go", "Vue", "复盘", "开源", "MixedCase", "50%_技巧", "DraftOnly", "ArchivedOnly", "ScheduledOnly", "RemovedTag"}
	items := make(map[string]*model.Tag, len(names))
	for _, name := range names {
		item := model.Tag{Name: name, NormalizedName: strings.ToLower(name)}
		if err := ensure(tx, "normalized_name = ?", []interface{}{item.NormalizedName}, &item, options, result); err != nil {
			return nil, err
		}
		items[name] = &item
	}
	return items, nil
}

func seedPosts(tx *gorm.DB, options Options, tags map[string]*model.Tag, result *Result) error {
	for _, definition := range postSeeds() {
		item := model.Post{Slug: definition.Key, Title: definition.Title, Summary: definition.Summary, CategoryID: definition.Category, Author: "Treasure Doc", Content: definition.Content, PublishStatus: definition.Status, PublishedOn: definition.PublishedOn, PublishedAt: definition.PublishedAt, Pinned: definition.Pinned, Version: 1}
		created, err := ensureCreated(tx, "slug = ?", []interface{}{item.Slug}, &item, options, result)
		if err != nil {
			return err
		}
		if err = postRelations(tx, item.ID, definition.Tags, tags); err != nil {
			return err
		}
		if created && definition.Deleted {
			if err = tx.Delete(&item).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func seedDiaries(tx *gorm.DB, options Options, tags map[string]*model.Tag, result *Result) error {
	for _, definition := range diarySeeds() {
		item := model.Diary{PublicID: definition.Key, Title: definition.Title, Summary: definition.Summary, Content: definition.Content, Mood: "充实", Weather: "晴", PublishStatus: definition.Status, PublishedOn: definition.PublishedOn, PublishedAt: definition.PublishedAt, Pinned: definition.Pinned, Version: 1}
		created, err := ensureCreated(tx, "public_id = ?", []interface{}{item.PublicID}, &item, options, result)
		if err != nil {
			return err
		}
		if err = diaryRelations(tx, item.ID, definition.Tags, tags); err != nil {
			return err
		}
		if created && definition.Deleted {
			if err = tx.Delete(&item).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func seedPortfolio(tx *gorm.DB, options Options, result *Result) error {
	tech, err := toJSON([]string{"Vue 3", "Go", "PostgreSQL"})
	if err != nil {
		return err
	}
	links, err := toJSON([]blogresponse.PortfolioLink{{Label: "GitHub", URL: "https://github.com/"}})
	if err != nil {
		return err
	}
	for index := 1; index <= 8; index++ {
		status, at := otherState(index)
		item := model.PortfolioItem{Slug: fmt.Sprintf("mock-work-%02d", index), Title: fmt.Sprintf("Mock 作品 %02d", index), Summary: "固定作品示例", CategoryID: []string{"website", "application"}[(index-1)%2], Cover: "📚", TechStack: tech, Links: links, Content: "# Mock 作品", PublishStatus: status, PublishedOn: dateFor(index), PublishedAt: at, Version: 1}
		created, err := ensureCreated(tx, "slug = ?", []interface{}{item.Slug}, &item, options, result)
		if err != nil {
			return err
		}
		if created && index == 8 {
			if err = tx.Delete(&item).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func seedTools(tx *gorm.DB, options Options, result *Result) error {
	for index := 1; index <= 8; index++ {
		status, at := otherState(index)
		kind, url, cover, devStatus, content := "own", "", "🧰", "开发中", "# Mock 工具"
		if index%2 == 0 {
			kind, url, cover, devStatus, content = "link", "https://developer.mozilla.org/zh-CN/", "", "", ""
		}
		item := model.Tool{Slug: fmt.Sprintf("mock-tool-%02d", index), Kind: kind, Name: fmt.Sprintf("Mock 工具 %02d", index), Description: "固定工具说明", URL: url, Cover: cover, DevelopmentStatus: devStatus, Content: content, PublishStatus: status, PublishedAt: at, SortOrder: index, Version: 1}
		created, err := ensureCreated(tx, "slug = ?", []interface{}{item.Slug}, &item, options, result)
		if err != nil {
			return err
		}
		if created && index == 8 {
			if err = tx.Delete(&item).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func seedBookmarks(tx *gorm.DB, options Options, tags map[string]*model.Tag, result *Result) error {
	for index := 1; index <= 8; index++ {
		status, at := otherState(index)
		item := model.Bookmark{PublicID: fmt.Sprintf("mock-bookmark-%02d", index), Title: fmt.Sprintf("Mock 书签 %02d", index), URL: fmt.Sprintf("https://example.com/mock/%02d", index), Description: "固定书签描述", CategoryID: []string{"dev", "design"}[(index-1)%2], Icon: "🔖", PublishStatus: status, PublishedAt: at, SortOrder: index, Version: 1}
		created, err := ensureCreated(tx, "public_id = ?", []interface{}{item.PublicID}, &item, options, result)
		if err != nil {
			return err
		}
		if err = bookmarkRelations(tx, item.ID, []string{"开源"}, tags); err != nil {
			return err
		}
		if created && index == 8 {
			if err = tx.Delete(&item).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func seedSettings(tx *gorm.DB, options Options, result *Result) error {
	links, err := toJSON([]blogresponse.ProfileLink{{ID: "github", Label: "GitHub", Value: "github.com", URL: "https://github.com/", Icon: "🔗"}})
	if err != nil {
		return err
	}
	skills, err := toJSON([]blogresponse.ProfileSkill{{Name: "Vue 3", Level: 90, Group: "前端"}, {Name: "Go", Level: 80, Group: "后端"}})
	if err != nil {
		return err
	}
	profile := model.Profile{ProfileKey: "default", Name: "Treasure", Avatar: "T", Role: "全栈开发者", Location: "中国", Motto: "写代码，写生活。", Bio: "Blog Seed 演示资料。", Links: links, Skills: skills}
	if err = ensure(tx, "profile_key = ?", []interface{}{"default"}, &profile, options, result); err != nil {
		return err
	}
	tech, err := toJSON([]string{"Vue 3", "TypeScript", "Go"})
	if err != nil {
		return err
	}
	modules, err := toJSON([]blogresponse.SiteModule{{ID: "blog", Icon: "📝", Name: "文章", Desc: "技术笔记与长文", Path: "/Blog"}})
	if err != nil {
		return err
	}
	milestones, err := toJSON([]blogresponse.SiteMilestone{{Date: "2026-09-04", Title: "Mock 数据上线", Desc: "加入幂等 Seed。"}})
	if err != nil {
		return err
	}
	site := model.Site{SiteKey: "default", Name: "Treasure Blog", Slogan: "records · thinking · life", Intro: "Blog Seed 演示站点。", TechStack: tech, Modules: modules, Milestones: milestones}
	return ensure(tx, "site_key = ?", []interface{}{"default"}, &site, options, result)
}

func ensure(tx *gorm.DB, query string, args []interface{}, value interface{}, options Options, result *Result) error {
	_, err := ensureCreated(tx, query, args, value, options, result)
	return err
}
func ensureCreated(tx *gorm.DB, query string, args []interface{}, value interface{}, options Options, result *Result) (bool, error) {
	err := tx.Unscoped().Where(query, args...).First(value).Error
	if err == nil {
		deleted := deletedAt(value)
		if deleted.Valid {
			if options.RestoreDeleted {
				if err = tx.Unscoped().Model(value).Update("deleted_at", nil).Error; err != nil {
					return false, err
				}
				result.Restored++
			} else {
				result.Skipped++
			}
		} else {
			result.Existing++
		}
		return false, nil
	}
	if err != gorm.ErrRecordNotFound {
		return false, err
	}
	if err = tx.Create(value).Error; err != nil {
		return false, err
	}
	result.Created++
	return true, nil
}
func deletedAt(value interface{}) gorm.DeletedAt {
	switch item := value.(type) {
	case *model.Category:
		return item.DeletedAt
	case *model.Tag:
		return item.DeletedAt
	case *model.Post:
		return item.DeletedAt
	case *model.Diary:
		return item.DeletedAt
	case *model.PortfolioItem:
		return item.DeletedAt
	case *model.Tool:
		return item.DeletedAt
	case *model.Bookmark:
		return item.DeletedAt
	case *model.Profile:
		return item.DeletedAt
	case *model.Site:
		return item.DeletedAt
	}
	return gorm.DeletedAt{}
}
func toJSON(value interface{}) (model.JSON, error) {
	data, err := json.Marshal(value)
	return model.JSON(data), err
}
func tagIDs(names []string, tags map[string]*model.Tag) []string {
	ids := make([]string, 0, len(names))
	for _, name := range names {
		if tag := tags[name]; tag != nil && tag.ID != "" {
			ids = append(ids, tag.ID)
		}
	}
	return ids
}
func postRelations(tx *gorm.DB, id string, names []string, tags map[string]*model.Tag) error {
	for _, tagID := range tagIDs(names, tags) {
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&model.PostTag{PostID: id, TagID: tagID, CreatedAt: time.Now()}).Error; err != nil {
			return err
		}
	}
	return nil
}
func diaryRelations(tx *gorm.DB, id string, names []string, tags map[string]*model.Tag) error {
	for _, tagID := range tagIDs(names, tags) {
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&model.DiaryTag{DiaryID: id, TagID: tagID, CreatedAt: time.Now()}).Error; err != nil {
			return err
		}
	}
	return nil
}
func bookmarkRelations(tx *gorm.DB, id string, names []string, tags map[string]*model.Tag) error {
	for _, tagID := range tagIDs(names, tags) {
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&model.BookmarkTag{BookmarkID: id, TagID: tagID, CreatedAt: time.Now()}).Error; err != nil {
			return err
		}
	}
	return nil
}
