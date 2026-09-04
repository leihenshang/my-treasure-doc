package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	blogmodel "fastduck/treasure-doc/module/blog/data/model"
	blogresponse "fastduck/treasure-doc/module/blog/data/response"
	"fastduck/treasure-doc/module/blog_mgr/data/request"

	"gorm.io/gorm"
)

var (
	ErrNotFound = errors.New("resource not found")
	ErrConflict = errors.New("resource conflict")
	ErrInvalid  = errors.New("invalid resource")
)

type DBProvider func() *gorm.DB
type Service struct{ db DBProvider }

type Page struct {
	List       interface{} `json:"list"`
	Pagination Pagination  `json:"pagination"`
}
type Pagination struct {
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	Total    int64  `json:"total"`
	OrderBy  string `json:"orderBy"`
}

func New(db DBProvider) *Service { return &Service{db: db} }
func (s *Service) database(ctx context.Context) (*gorm.DB, error) {
	if s.db == nil || s.db() == nil {
		return nil, errors.New("database is not initialized")
	}
	return s.db().WithContext(ctx), nil
}

func modelFor(resource string) (interface{}, interface{}, string, error) {
	switch resource {
	case "categories":
		return &blogmodel.Category{}, &[]blogmodel.Category{}, "name", nil
	case "tags":
		return &blogmodel.Tag{}, &[]blogmodel.Tag{}, "name", nil
	case "posts":
		return &blogmodel.Post{}, &[]blogmodel.Post{}, "title", nil
	case "diaries":
		return &blogmodel.Diary{}, &[]blogmodel.Diary{}, "title", nil
	case "portfolio-items":
		return &blogmodel.PortfolioItem{}, &[]blogmodel.PortfolioItem{}, "title", nil
	case "tools":
		return &blogmodel.Tool{}, &[]blogmodel.Tool{}, "name", nil
	case "bookmarks":
		return &blogmodel.Bookmark{}, &[]blogmodel.Bookmark{}, "title", nil
	default:
		return nil, nil, "", ErrInvalid
	}
}

func applyDeleted(db *gorm.DB, mode string) *gorm.DB {
	if mode == "all" {
		return db.Unscoped()
	}
	if mode == "only" {
		return db.Unscoped().Where("deleted_at IS NOT NULL")
	}
	return db
}

func (s *Service) List(ctx context.Context, resource string, query request.List) (Page, error) {
	db, err := s.database(ctx)
	if err != nil {
		return Page{}, err
	}
	item, list, keywordColumn, err := modelFor(resource)
	if err != nil {
		return Page{}, err
	}
	q := applyDeleted(db.Model(item), query.Deleted)
	if query.Keyword != "" {
		q = q.Where("LOWER("+keywordColumn+") LIKE ?", "%"+strings.ToLower(query.Keyword)+"%")
	}
	if query.Status != "" {
		q = q.Where("publish_status = ?", query.Status)
	}
	if query.Scope != "" {
		q = q.Where("scope = ?", query.Scope)
	}
	if query.CategoryID != "" {
		q = q.Where("category_id = ?", query.CategoryID)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return Page{}, err
	}
	order := "created_at DESC"
	if query.Sort == "asc" {
		order = "created_at ASC"
	}
	if err := q.Order(order).Offset(query.Offset()).Limit(query.PageSize).Find(list).Error; err != nil {
		return Page{}, err
	}
	return Page{List: list, Pagination: Pagination{Page: query.Page, PageSize: query.PageSize, Total: total, OrderBy: "created_at_" + query.Sort}}, nil
}

func (s *Service) Get(ctx context.Context, resource, id string) (interface{}, error) {
	db, err := s.database(ctx)
	if err != nil {
		return nil, err
	}
	item, _, _, err := modelFor(resource)
	if err != nil {
		return nil, err
	}
	if err := db.Unscoped().Where("id = ?", id).First(item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return item, nil
}

func publishedTimes(status, date string, at *time.Time) (time.Time, time.Time, error) {
	var publishedOn time.Time
	var err error
	if date != "" {
		publishedOn, err = time.Parse("2006-01-02", date)
		if err != nil {
			return time.Time{}, time.Time{}, ErrInvalid
		}
	}
	publishedAt := time.Time{}
	if at != nil {
		publishedAt = *at
	}
	if status == blogmodel.StatusPublished && publishedAt.IsZero() {
		publishedAt = time.Now()
	}
	if status == blogmodel.StatusPublished && publishedOn.IsZero() {
		publishedOn = publishedAt
	}
	return publishedOn, publishedAt, nil
}

func marshal(value interface{}) (blogmodel.JSON, error) {
	data, err := json.Marshal(value)
	return blogmodel.JSON(data), err
}

func (s *Service) Create(ctx context.Context, resource string, payload interface{}) (interface{}, error) {
	db, err := s.database(ctx)
	if err != nil {
		return nil, err
	}
	var result interface{}
	err = db.Transaction(func(tx *gorm.DB) error {
		item, tagIDs, relation, err := buildModel(resource, payload)
		if err != nil {
			return err
		}
		if err := validateReferences(tx, resource, item, tagIDs); err != nil {
			return err
		}
		if err := tx.Create(item).Error; err != nil {
			return mapDBError(err)
		}
		if relation != "" {
			if err := replaceTags(tx, relation, modelID(item), tagIDs); err != nil {
				return err
			}
		}
		result = item
		return nil
	})
	return result, err
}

func (s *Service) Update(ctx context.Context, resource, id string, payload interface{}) (interface{}, error) {
	db, err := s.database(ctx)
	if err != nil {
		return nil, err
	}
	var result interface{}
	err = db.Transaction(func(tx *gorm.DB) error {
		if requiresVersion(resource) && requestedVersion(payload) < 1 {
			return ErrInvalid
		}
		item, tagIDs, relation, err := buildModel(resource, payload)
		if err != nil {
			return err
		}
		if err := validateReferences(tx, resource, item, tagIDs); err != nil {
			return err
		}
		if resource == "categories" {
			if err := updateCategory(tx, id, item.(*blogmodel.Category)); err != nil {
				return err
			}
			loaded, err := s.getWithDB(tx, resource, id)
			result = loaded
			return err
		}
		version := modelVersion(item)
		updates := updateMap(item)
		q := tx.Model(item).Where("id = ?", id)
		if version > 0 {
			q = q.Where("version = ?", version)
			updates["version"] = gorm.Expr("version + 1")
		}
		resultDB := q.Updates(updates)
		if resultDB.Error != nil {
			return mapDBError(resultDB.Error)
		}
		if resultDB.RowsAffected == 0 {
			return ErrConflict
		}
		if relation != "" {
			if err := replaceTags(tx, relation, id, tagIDs); err != nil {
				return err
			}
		}
		loaded, err := s.getWithDB(tx, resource, id)
		result = loaded
		return err
	})
	return result, err
}

func (s *Service) Delete(ctx context.Context, resource, id string) error {
	db, err := s.database(ctx)
	if err != nil {
		return err
	}
	item, _, _, err := modelFor(resource)
	if resource == "categories" {
		var category blogmodel.Category
		if err := db.Where("id = ?", id).First(&category).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}
		table := map[string]string{blogmodel.CategoryPost: "td_blog_post", blogmodel.CategoryPortfolio: "td_blog_portfolio_item", blogmodel.CategoryBookmark: "td_blog_bookmark"}[category.Scope]
		var count int64
		if err := db.Table(table).Where("category_id = ? AND deleted_at IS NULL", category.Slug).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return ErrConflict
		}
	}
	if err != nil {
		return err
	}
	result := db.Where("id = ?", id).Delete(item)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Service) Restore(ctx context.Context, resource, id string) error {
	db, err := s.database(ctx)
	if err != nil {
		return err
	}
	item, _, _, err := modelFor(resource)
	if err != nil {
		return err
	}
	result := db.Unscoped().Model(item).Where("id = ? AND deleted_at IS NOT NULL", id).Update("deleted_at", nil)
	if result.Error != nil {
		return mapDBError(result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Service) getWithDB(db *gorm.DB, resource, id string) (interface{}, error) {
	item, _, _, err := modelFor(resource)
	if err != nil {
		return nil, err
	}
	if err := db.Unscoped().Where("id = ?", id).First(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

func mapDBError(err error) error {
	if strings.Contains(err.Error(), "SQLSTATE 23505") || strings.Contains(strings.ToLower(err.Error()), "duplicate key") {
		return ErrConflict
	}
	return err
}

func modelID(item interface{}) string {
	switch value := item.(type) {
	case *blogmodel.Post:
		return value.ID
	case *blogmodel.Diary:
		return value.ID
	case *blogmodel.Bookmark:
		return value.ID
	}
	return ""
}
func modelVersion(item interface{}) int {
	switch value := item.(type) {
	case *blogmodel.Post:
		return value.Version
	case *blogmodel.Diary:
		return value.Version
	case *blogmodel.PortfolioItem:
		return value.Version
	case *blogmodel.Tool:
		return value.Version
	case *blogmodel.Bookmark:
		return value.Version
	}
	return 0
}

func requiresVersion(resource string) bool {
	return resource == "posts" || resource == "diaries" || resource == "portfolio-items" || resource == "tools" || resource == "bookmarks"
}

func requestedVersion(payload interface{}) int {
	switch value := payload.(type) {
	case request.Post:
		return value.Version
	case request.Diary:
		return value.Version
	case request.Portfolio:
		return value.Version
	case request.Tool:
		return value.Version
	case request.Bookmark:
		return value.Version
	}
	return 0
}

func replaceTags(tx *gorm.DB, relation, ownerID string, tagIDs []string) error {
	ownerColumn := map[string]string{"td_blog_post_tag": "post_id", "td_blog_diary_tag": "diary_id", "td_blog_bookmark_tag": "bookmark_id"}[relation]
	var relationModel interface{}
	switch relation {
	case "td_blog_post_tag":
		relationModel = &blogmodel.PostTag{}
	case "td_blog_diary_tag":
		relationModel = &blogmodel.DiaryTag{}
	case "td_blog_bookmark_tag":
		relationModel = &blogmodel.BookmarkTag{}
	default:
		return ErrInvalid
	}
	if err := tx.Where(ownerColumn+" = ?", ownerID).Delete(relationModel).Error; err != nil {
		return err
	}
	for _, tagID := range tagIDs {
		if err := tx.Table(relation).Create(map[string]interface{}{ownerColumn: ownerID, "tag_id": tagID, "created_at": time.Now()}).Error; err != nil {
			return err
		}
	}
	return nil
}

func validateReferences(tx *gorm.DB, resource string, item interface{}, tagIDs []string) error {
	var scope, category string
	switch value := item.(type) {
	case *blogmodel.Post:
		scope, category = blogmodel.CategoryPost, value.CategoryID
	case *blogmodel.PortfolioItem:
		scope, category = blogmodel.CategoryPortfolio, value.CategoryID
	case *blogmodel.Bookmark:
		scope, category = blogmodel.CategoryBookmark, value.CategoryID
	}
	if category != "" {
		var count int64
		if err := tx.Model(&blogmodel.Category{}).Where("scope = ? AND slug = ?", scope, category).Count(&count).Error; err != nil {
			return err
		}
		if count != 1 {
			return ErrInvalid
		}
	}
	if len(tagIDs) > 0 {
		var count int64
		if err := tx.Model(&blogmodel.Tag{}).Where("id IN ?", tagIDs).Count(&count).Error; err != nil {
			return err
		}
		if count != int64(len(tagIDs)) {
			return ErrInvalid
		}
	}
	_ = resource
	return nil
}

func buildModel(resource string, payload interface{}) (interface{}, []string, string, error) {
	switch value := payload.(type) {
	case request.Category:
		if !request.ValidScope(value.Scope) || !request.ValidID(value.Slug) || strings.TrimSpace(value.Name) == "" {
			return nil, nil, "", ErrInvalid
		}
		enabled := true
		if value.Enabled != nil {
			enabled = *value.Enabled
		}
		return &blogmodel.Category{Scope: value.Scope, Slug: value.Slug, Name: value.Name, SortOrder: value.SortOrder, Enabled: enabled}, nil, "", nil
	case request.Tag:
		name := strings.TrimSpace(value.Name)
		if name == "" {
			return nil, nil, "", ErrInvalid
		}
		return &blogmodel.Tag{Name: name, NormalizedName: strings.ToLower(name)}, nil, "", nil
	case request.Post:
		if !request.ValidID(value.Slug) || value.Title == "" || !request.ValidStatus(value.PublishStatus) {
			return nil, nil, "", ErrInvalid
		}
		on, at, err := publishedTimes(value.PublishStatus, value.PublishedOn, value.PublishedAt)
		ids, err2 := request.NormalizeIDs(value.TagIDs)
		if err != nil || err2 != nil {
			return nil, nil, "", ErrInvalid
		}
		return &blogmodel.Post{Slug: value.Slug, Title: value.Title, Summary: value.Summary, CategoryID: value.CategoryID, Author: value.Author, Content: value.Content, PublishStatus: value.PublishStatus, PublishedOn: on, PublishedAt: at, Pinned: value.Pinned, Version: max(value.Version, 1)}, ids, "td_blog_post_tag", nil
	case request.Diary:
		if !request.ValidID(value.PublicID) || value.Title == "" || !request.ValidStatus(value.PublishStatus) {
			return nil, nil, "", ErrInvalid
		}
		on, at, err := publishedTimes(value.PublishStatus, value.PublishedOn, value.PublishedAt)
		ids, err2 := request.NormalizeIDs(value.TagIDs)
		if err != nil || err2 != nil {
			return nil, nil, "", ErrInvalid
		}
		return &blogmodel.Diary{PublicID: value.PublicID, Title: value.Title, Summary: value.Summary, Content: value.Content, Mood: value.Mood, Weather: value.Weather, PublishStatus: value.PublishStatus, PublishedOn: on, PublishedAt: at, Pinned: value.Pinned, Version: max(value.Version, 1)}, ids, "td_blog_diary_tag", nil
	case request.Portfolio:
		if !request.ValidID(value.Slug) || value.Title == "" || !request.ValidStatus(value.PublishStatus) {
			return nil, nil, "", ErrInvalid
		}
		on, at, err := publishedTimes(value.PublishStatus, value.PublishedOn, value.PublishedAt)
		tech, e1 := marshal(value.TechStack)
		links, e2 := marshal(value.Links)
		if err != nil || e1 != nil || e2 != nil {
			return nil, nil, "", ErrInvalid
		}
		return &blogmodel.PortfolioItem{Slug: value.Slug, Title: value.Title, Summary: value.Summary, CategoryID: value.CategoryID, Cover: value.Cover, TechStack: tech, Links: links, Content: value.Content, PublishStatus: value.PublishStatus, PublishedOn: on, PublishedAt: at, Version: max(value.Version, 1)}, nil, "", nil
	case request.Tool:
		if request.ValidateTool(value) != nil {
			return nil, nil, "", ErrInvalid
		}
		_, at, _ := publishedTimes(value.PublishStatus, "", value.PublishedAt)
		if value.Kind == "link" {
			value.Cover = ""
			value.DevelopmentStatus = ""
			value.Content = ""
		} else {
			value.URL = ""
		}
		return &blogmodel.Tool{Slug: value.Slug, Kind: value.Kind, Name: value.Name, Description: value.Description, URL: value.URL, Cover: value.Cover, DevelopmentStatus: value.DevelopmentStatus, Content: value.Content, PublishStatus: value.PublishStatus, PublishedAt: at, SortOrder: value.SortOrder, Version: max(value.Version, 1)}, nil, "", nil
	case request.Bookmark:
		if !request.ValidID(value.PublicID) || value.Title == "" || !request.ValidURL(value.URL, false) || !request.ValidStatus(value.PublishStatus) {
			return nil, nil, "", ErrInvalid
		}
		_, at, _ := publishedTimes(value.PublishStatus, "", value.PublishedAt)
		ids, err := request.NormalizeIDs(value.TagIDs)
		if err != nil {
			return nil, nil, "", ErrInvalid
		}
		return &blogmodel.Bookmark{PublicID: value.PublicID, Title: value.Title, URL: value.URL, Description: value.Description, CategoryID: value.CategoryID, Icon: value.Icon, PublishStatus: value.PublishStatus, PublishedAt: at, SortOrder: value.SortOrder, Version: max(value.Version, 1)}, ids, "td_blog_bookmark_tag", nil
	default:
		return nil, nil, "", fmt.Errorf("%w: %s", ErrInvalid, resource)
	}
}

func updateMap(item interface{}) map[string]interface{} {
	switch value := item.(type) {
	case *blogmodel.Tag:
		return map[string]interface{}{"name": value.Name, "normalized_name": value.NormalizedName}
	case *blogmodel.Post:
		return map[string]interface{}{"slug": value.Slug, "title": value.Title, "summary": value.Summary, "category_id": value.CategoryID, "author": value.Author, "content": value.Content, "publish_status": value.PublishStatus, "published_on": value.PublishedOn, "published_at": value.PublishedAt, "pinned": value.Pinned}
	case *blogmodel.Diary:
		return map[string]interface{}{"public_id": value.PublicID, "title": value.Title, "summary": value.Summary, "content": value.Content, "mood": value.Mood, "weather": value.Weather, "publish_status": value.PublishStatus, "published_on": value.PublishedOn, "published_at": value.PublishedAt, "pinned": value.Pinned}
	case *blogmodel.PortfolioItem:
		return map[string]interface{}{"slug": value.Slug, "title": value.Title, "summary": value.Summary, "category_id": value.CategoryID, "cover": value.Cover, "tech_stack": value.TechStack, "links": value.Links, "content": value.Content, "publish_status": value.PublishStatus, "published_on": value.PublishedOn, "published_at": value.PublishedAt}
	case *blogmodel.Tool:
		return map[string]interface{}{"slug": value.Slug, "kind": value.Kind, "name": value.Name, "description": value.Description, "url": value.URL, "cover": value.Cover, "development_status": value.DevelopmentStatus, "content": value.Content, "publish_status": value.PublishStatus, "published_at": value.PublishedAt, "sort_order": value.SortOrder}
	case *blogmodel.Bookmark:
		return map[string]interface{}{"public_id": value.PublicID, "title": value.Title, "url": value.URL, "description": value.Description, "category_id": value.CategoryID, "icon": value.Icon, "publish_status": value.PublishStatus, "published_at": value.PublishedAt, "sort_order": value.SortOrder}
	}
	return map[string]interface{}{}
}

func updateCategory(tx *gorm.DB, id string, next *blogmodel.Category) error {
	var current blogmodel.Category
	if err := tx.Where("id = ?", id).First(&current).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	if current.Scope != next.Scope {
		return ErrInvalid
	}
	if current.Slug != next.Slug {
		table := map[string]string{blogmodel.CategoryPost: "td_blog_post", blogmodel.CategoryPortfolio: "td_blog_portfolio_item", blogmodel.CategoryBookmark: "td_blog_bookmark"}[current.Scope]
		if err := tx.Table(table).Where("category_id = ?", current.Slug).Update("category_id", next.Slug).Error; err != nil {
			return err
		}
	}
	return mapDBError(tx.Model(&current).Updates(map[string]interface{}{"slug": next.Slug, "name": next.Name, "sort_order": next.SortOrder, "enabled": next.Enabled}).Error)
}

func (s *Service) GetSetting(ctx context.Context, name string) (interface{}, error) {
	db, err := s.database(ctx)
	if err != nil {
		return nil, err
	}
	if name == "profile" {
		value := &blogmodel.Profile{}
		err = db.Unscoped().Where("profile_key = ?", "default").First(value).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return blogresponse.Profile{Links: []blogresponse.ProfileLink{}, Skills: []blogresponse.ProfileSkill{}}, nil
		}
		if err != nil {
			return nil, err
		}
		var links []blogresponse.ProfileLink
		var skills []blogresponse.ProfileSkill
		if json.Unmarshal(value.Links, &links) != nil || json.Unmarshal(value.Skills, &skills) != nil {
			return nil, ErrInvalid
		}
		return blogresponse.Profile{Name: value.Name, Avatar: value.Avatar, Role: value.Role, Location: value.Location, Motto: value.Motto, Bio: value.Bio, Links: links, Skills: skills}, nil
	}
	value := &blogmodel.Site{}
	err = db.Unscoped().Where("site_key = ?", "default").First(value).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return blogresponse.Site{TechStack: []string{}, Modules: []blogresponse.SiteModule{}, Milestones: []blogresponse.SiteMilestone{}}, nil
	}
	if err != nil {
		return nil, err
	}
	var tech []string
	var modules []blogresponse.SiteModule
	var milestones []blogresponse.SiteMilestone
	if json.Unmarshal(value.TechStack, &tech) != nil || json.Unmarshal(value.Modules, &modules) != nil || json.Unmarshal(value.Milestones, &milestones) != nil {
		return nil, ErrInvalid
	}
	return blogresponse.Site{Name: value.Name, Slogan: value.Slogan, Intro: value.Intro, TechStack: tech, Modules: modules, Milestones: milestones}, nil
}

func (s *Service) PutSetting(ctx context.Context, name string, payload interface{}) (interface{}, error) {
	db, err := s.database(ctx)
	if err != nil {
		return nil, err
	}
	if name == "profile" {
		value := payload.(blogresponse.Profile)
		if err := validateProfile(value); err != nil {
			return nil, err
		}
		links, e1 := marshal(value.Links)
		skills, e2 := marshal(value.Skills)
		if e1 != nil || e2 != nil {
			return nil, ErrInvalid
		}
		item := &blogmodel.Profile{ProfileKey: "default", Name: value.Name, Avatar: value.Avatar, Role: value.Role, Location: value.Location, Motto: value.Motto, Bio: value.Bio, Links: links, Skills: skills}
		err = upsertSetting(db, &blogmodel.Profile{}, "profile_key", "default", item)
		return value, err
	}
	value := payload.(blogresponse.Site)
	if err := validateSite(value); err != nil {
		return nil, err
	}
	tech, e1 := marshal(value.TechStack)
	modules, e2 := marshal(value.Modules)
	milestones, e3 := marshal(value.Milestones)
	if e1 != nil || e2 != nil || e3 != nil {
		return nil, ErrInvalid
	}
	item := &blogmodel.Site{SiteKey: "default", Name: value.Name, Slogan: value.Slogan, Intro: value.Intro, TechStack: tech, Modules: modules, Milestones: milestones}
	err = upsertSetting(db, &blogmodel.Site{}, "site_key", "default", item)
	return value, err
}

func validateProfile(value blogresponse.Profile) error {
	if strings.TrimSpace(value.Name) == "" {
		return ErrInvalid
	}
	seen := map[string]struct{}{}
	for _, link := range value.Links {
		if strings.TrimSpace(link.ID) == "" || strings.TrimSpace(link.Label) == "" {
			return ErrInvalid
		}
		if _, ok := seen[link.ID]; ok {
			return ErrInvalid
		}
		seen[link.ID] = struct{}{}
		if link.URL != "" && !request.ValidURL(link.URL, true) {
			return ErrInvalid
		}
	}
	for _, skill := range value.Skills {
		if strings.TrimSpace(skill.Name) == "" || strings.TrimSpace(skill.Group) == "" || skill.Level < 0 || skill.Level > 100 {
			return ErrInvalid
		}
	}
	return nil
}

func validateSite(value blogresponse.Site) error {
	if strings.TrimSpace(value.Name) == "" {
		return ErrInvalid
	}
	seen := map[string]struct{}{}
	for _, module := range value.Modules {
		if strings.TrimSpace(module.ID) == "" || !strings.HasPrefix(module.Path, "/Blog") {
			return ErrInvalid
		}
		if _, ok := seen[module.ID]; ok {
			return ErrInvalid
		}
		seen[module.ID] = struct{}{}
	}
	for _, milestone := range value.Milestones {
		if !request.ValidDate(milestone.Date) || strings.TrimSpace(milestone.Title) == "" {
			return ErrInvalid
		}
	}
	return nil
}

func upsertSetting(db *gorm.DB, existing interface{}, keyColumn, keyValue string, values interface{}) error {
	return db.Transaction(func(tx *gorm.DB) error {
		result := tx.Unscoped().Where(keyColumn+" = ?", keyValue).First(existing)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return tx.Create(values).Error
		}
		if result.Error != nil {
			return result.Error
		}
		return tx.Unscoped().Model(existing).Updates(values).Update("deleted_at", nil).Error
	})
}
