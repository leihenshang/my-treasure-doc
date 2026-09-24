package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	blogmodel "fastduck/treasure-doc/module/blog/data/model"
	blogresponse "fastduck/treasure-doc/module/blog/data/response"
	"fastduck/treasure-doc/module/blog_mgr/data/request"
	"fastduck/treasure-doc/module/user/global"

	"gorm.io/gorm"
)

var (
	ErrNotFound = errors.New("resource not found")
	ErrConflict = errors.New("resource conflict")
	ErrInvalid  = errors.New("invalid resource")
	// ErrReferenceNotFound 表示 payload 引用了不存在的分类或标签，
	// 与「字段格式非法」区分开，便于前端提示用户先创建对应分类/标签。
	ErrReferenceNotFound = errors.New("referenced category or tag not found")
	// 工具字段缺失的细分错误，避免统一落到「请求参数格式错误」这种无法定位的提示。
	ErrToolURLRequired    = errors.New("tool url required")
	ErrToolStatusRequired = errors.New("tool development status required")
)

type Service struct{}

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

type PostWithTags struct {
	blogmodel.Post
	TagIDs []string `json:"tagIds"`
}

type DiaryWithTags struct {
	blogmodel.Diary
	TagIDs []string `json:"tagIds"`
}

type BookmarkWithTags struct {
	blogmodel.Bookmark
	TagIDs []string `json:"tagIds"`
}

func New() *Service { return &Service{} }
func (s *Service) database(ctx context.Context) (*gorm.DB, error) {
	if global.Db == nil {
		return nil, errors.New("database is not initialized")
	}
	return global.Db.WithContext(ctx), nil
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
	enriched, err := enrichListWithTags(db, resource, list)
	if err != nil {
		return Page{}, err
	}
	return Page{List: enriched, Pagination: Pagination{Page: query.Page, PageSize: query.PageSize, Total: total, OrderBy: "created_at_" + query.Sort}}, nil
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
	return enrichItemWithTags(db, resource, item)
}

func publishedTimes(status, date string, at *time.Time) (time.Time, time.Time, error) {
	var publishedOn time.Time
	var err error
	if date != "" {
		publishedOn, err = time.Parse("2006-01-02", date)
		if err != nil {
			return time.Time{}, time.Time{}, request.Field("publishedOn", "创建日期格式必须是 YYYY-MM-DD")
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
		if err := resolveIdentifier(tx, resource, "", item, true); err != nil {
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
		result, err = enrichItemWithTags(tx, resource, item)
		if err != nil {
			return err
		}
		return recordEditHistory(tx, resource, result)
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
			return request.Field("version", "缺少 version：请先读取最新数据再提交（乐观锁）")
		}
		item, tagIDs, relation, err := buildModel(resource, payload)
		if err != nil {
			return err
		}
		if err := resolveIdentifier(tx, resource, id, item, false); err != nil {
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

// UpdateFields 只更新白名单内的单个字段（列表快捷设置：置顶、发布状态）。
// 与完整更新不同，它不要求客户端传 version，更新成功后 version 自增。
func (s *Service) UpdateFields(ctx context.Context, resource, id string, fields map[string]interface{}) (interface{}, error) {
	if !requiresVersion(resource) {
		return nil, ErrInvalid
	}
	if !request.ValidID(id) {
		return nil, request.Field("id", "记录 ID 不合法")
	}
	updates := map[string]interface{}{}
	if value, ok := fields["pinned"]; ok {
		if resource != "posts" && resource != "diaries" {
			return nil, request.Field("pinned", "置顶只支持文章与日记")
		}
		pinned, ok := value.(bool)
		if !ok {
			return nil, request.Field("pinned", "置顶取值必须是布尔值")
		}
		updates["pinned"] = pinned
	}
	if value, ok := fields["publishStatus"]; ok {
		status, ok := value.(string)
		if !ok || !request.ValidStatus(status) {
			return nil, request.Field("publishStatus", "发布状态只能是 draft / published / archived")
		}
		updates["publish_status"] = status
	}
	if len(updates) == 0 || len(updates) != len(fields) {
		return nil, request.Field("fields", "包含不支持的字段，快捷修改只支持 pinned 与 publishStatus")
	}
	updates["version"] = gorm.Expr("version + 1")
	db, err := s.database(ctx)
	if err != nil {
		return nil, err
	}
	item, _, _, err := modelFor(resource)
	if err != nil {
		return nil, err
	}
	var result interface{}
	err = db.Transaction(func(tx *gorm.DB) error {
		// 回收站中的记录不允许快捷设置
		if err := tx.Where("id = ?", id).First(item).Error; err != nil {
			return err
		}
		res := tx.Model(item).Where("id = ?", id).Updates(updates)
		if res.Error != nil {
			return mapDBError(res.Error)
		}
		if res.RowsAffected == 0 {
			return ErrNotFound
		}
		loaded, err := s.getWithDB(tx, resource, id)
		if err != nil {
			return err
		}
		result, err = enrichItemWithTags(tx, resource, loaded)
		return err
	})
	return result, err
}

type tagRelation struct {
	OwnerID string `gorm:"column:owner_id"`
	TagID   string `gorm:"column:tag_id"`
}

func loadTagIDs(db *gorm.DB, relationTable, ownerColumn string, ownerIDs []string) (map[string][]string, error) {
	result := make(map[string][]string, len(ownerIDs))
	for _, ownerID := range ownerIDs {
		result[ownerID] = []string{}
	}
	if len(ownerIDs) == 0 {
		return result, nil
	}
	var relations []tagRelation
	err := db.Table(relationTable).
		Select(ownerColumn+" AS owner_id", "tag_id").
		Where(ownerColumn+" IN ?", ownerIDs).
		Order(ownerColumn + " ASC, created_at ASC, tag_id ASC").
		Scan(&relations).Error
	if err != nil {
		return nil, err
	}
	for _, relation := range relations {
		result[relation.OwnerID] = append(result[relation.OwnerID], relation.TagID)
	}
	return result, nil
}

func enrichListWithTags(db *gorm.DB, resource string, list interface{}) (interface{}, error) {
	switch values := list.(type) {
	case *[]blogmodel.Post:
		ids := make([]string, 0, len(*values))
		for _, value := range *values {
			ids = append(ids, value.ID)
		}
		tags, err := loadTagIDs(db, "td_blog_post_tag", "post_id", ids)
		if err != nil {
			return nil, err
		}
		result := make([]PostWithTags, 0, len(*values))
		for _, value := range *values {
			result = append(result, PostWithTags{Post: value, TagIDs: tags[value.ID]})
		}
		return result, nil
	case *[]blogmodel.Diary:
		ids := make([]string, 0, len(*values))
		for _, value := range *values {
			ids = append(ids, value.ID)
		}
		tags, err := loadTagIDs(db, "td_blog_diary_tag", "diary_id", ids)
		if err != nil {
			return nil, err
		}
		result := make([]DiaryWithTags, 0, len(*values))
		for _, value := range *values {
			result = append(result, DiaryWithTags{Diary: value, TagIDs: tags[value.ID]})
		}
		return result, nil
	case *[]blogmodel.Bookmark:
		ids := make([]string, 0, len(*values))
		for _, value := range *values {
			ids = append(ids, value.ID)
		}
		tags, err := loadTagIDs(db, "td_blog_bookmark_tag", "bookmark_id", ids)
		if err != nil {
			return nil, err
		}
		result := make([]BookmarkWithTags, 0, len(*values))
		for _, value := range *values {
			result = append(result, BookmarkWithTags{Bookmark: value, TagIDs: tags[value.ID]})
		}
		return result, nil
	default:
		return list, nil
	}
}

func enrichItemWithTags(db *gorm.DB, resource string, item interface{}) (interface{}, error) {
	switch value := item.(type) {
	case *blogmodel.Post:
		tags, err := loadTagIDs(db, "td_blog_post_tag", "post_id", []string{value.ID})
		return &PostWithTags{Post: *value, TagIDs: tags[value.ID]}, err
	case *blogmodel.Diary:
		tags, err := loadTagIDs(db, "td_blog_diary_tag", "diary_id", []string{value.ID})
		return &DiaryWithTags{Diary: *value, TagIDs: tags[value.ID]}, err
	case *blogmodel.Bookmark:
		tags, err := loadTagIDs(db, "td_blog_bookmark_tag", "bookmark_id", []string{value.ID})
		return &BookmarkWithTags{Bookmark: *value, TagIDs: tags[value.ID]}, err
	default:
		return item, nil
	}
}

func (s *Service) Delete(ctx context.Context, resource, id string) error {
	db, err := s.database(ctx)
	if err != nil {
		return err
	}
	item, _, _, err := modelFor(resource)
	if resource == "categories" {
		if err := ensureCategoriesUnreferenced(db, []string{id}); err != nil {
			return err
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

// DeleteMany 批量软删除，返回实际删除条数；重复或已删除的 ID 会被自动忽略。
// 语义与单条删除一致：分类被内容引用时整体拒绝（ErrConflict），不会部分删除。
func (s *Service) DeleteMany(ctx context.Context, resource string, ids []string) (int64, error) {
	if len(ids) == 0 {
		return 0, ErrInvalid
	}
	db, err := s.database(ctx)
	if err != nil {
		return 0, err
	}
	item, _, _, err := modelFor(resource)
	if err != nil {
		return 0, err
	}
	if resource == "categories" {
		if err := ensureCategoriesUnreferenced(db, ids); err != nil {
			return 0, err
		}
	}
	result := db.Where("id IN ?", ids).Delete(item)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// ensureCategoriesUnreferenced 校验分类未被其归属内容引用。
// 分类通过 slug 被内容引用，归属表由 scope 决定。
func ensureCategoriesUnreferenced(db *gorm.DB, ids []string) error {
	var categories []blogmodel.Category
	if err := db.Where("id IN ?", ids).Find(&categories).Error; err != nil {
		return err
	}
	// 默认分类不可删除：无论是否被引用都拒绝
	for _, category := range categories {
		if category.Slug == defaultCategorySlug {
			return ErrConflict
		}
	}
	tables := map[string]string{
		blogmodel.CategoryPost:      "td_blog_post",
		blogmodel.CategoryPortfolio: "td_blog_portfolio_item",
		blogmodel.CategoryBookmark:  "td_blog_bookmark",
	}
	for _, category := range categories {
		table := tables[category.Scope]
		if table == "" {
			continue
		}
		var count int64
		if err := db.Table(table).Where("category_id = ? AND deleted_at IS NULL", category.Slug).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return ErrConflict
		}
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

// mapDBError 把数据库层的唯一键冲突翻译成业务错误。
// 必须容忍 nil：调用方可能直接把更新结果交进来（成功时 Error 为 nil），
// 早期写法漏了这一层，导致 err.Error() 空指针 panic。
func mapDBError(err error) error {
	if err == nil {
		return nil
	}
	if strings.Contains(err.Error(), "SQLSTATE 23505") || strings.Contains(strings.ToLower(err.Error()), "duplicate key") || strings.Contains(strings.ToLower(err.Error()), "unique constraint failed") {
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
	if category == "" && scope != "" {
		// 未指定分类：自动落到该 scope 的默认分类（不存在则创建），保证内容始终有归属
		if err := ensureDefaultCategory(tx, scope); err != nil {
			return err
		}
		assignDefaultCategory(item)
	} else if category != "" {
		var count int64
		if err := tx.Model(&blogmodel.Category{}).Where("scope = ? AND slug = ?", scope, category).Count(&count).Error; err != nil {
			return err
		}
		if count != 1 {
			return ErrReferenceNotFound
		}
	}
	if len(tagIDs) > 0 {
		var count int64
		if err := tx.Model(&blogmodel.Tag{}).Where("id IN ?", tagIDs).Count(&count).Error; err != nil {
			return err
		}
		if count != int64(len(tagIDs)) {
			return ErrReferenceNotFound
		}
	}
	_ = resource
	return nil
}

// defaultCategorySlug / defaultCategoryName 各 scope 的默认（未分类）分类，保留名、不可删除。
const (
	defaultCategoryName = "默认分类"
	defaultCategorySlug = "uncategorized"
)

// assignDefaultCategory 把内容的分类落到该 scope 的默认分类。
func assignDefaultCategory(item interface{}) {
	switch value := item.(type) {
	case *blogmodel.Post:
		value.CategoryID = defaultCategorySlug
	case *blogmodel.PortfolioItem:
		value.CategoryID = defaultCategorySlug
	case *blogmodel.Bookmark:
		value.CategoryID = defaultCategorySlug
	}
}

// ensureDefaultCategory 确保该 scope 的默认分类存在；并发重复创建时容忍唯一冲突。
func ensureDefaultCategory(tx *gorm.DB, scope string) error {
	var count int64
	if err := tx.Model(&blogmodel.Category{}).Where("scope = ? AND slug = ?", scope, defaultCategorySlug).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	cat := blogmodel.Category{Scope: scope, Slug: defaultCategorySlug, Name: defaultCategoryName, SortOrder: 9999, Enabled: true}
	if err := tx.Create(&cat).Error; err != nil {
		// 并发下已由其他请求创建，视为成功
		var again int64
		if qerr := tx.Model(&blogmodel.Category{}).Where("scope = ? AND slug = ?", scope, defaultCategorySlug).Count(&again).Error; qerr != nil {
			return qerr
		}
		if again > 0 {
			return nil
		}
		return err
	}
	return nil
}

// identifierOf 返回模型当前的公开标识（post/portfolio/tool 的 slug、diary 的 publicId）；
// 不参与自动生成标识的资源（分类/标签/收藏集等）返回 ok=false。
func identifierOf(item interface{}) (string, bool) {
	switch v := item.(type) {
	case *blogmodel.Post:
		return v.Slug, true
	case *blogmodel.Diary:
		return v.PublicID, true
	case *blogmodel.PortfolioItem:
		return v.Slug, true
	case *blogmodel.Tool:
		return v.Slug, true
	default:
		return "", false
	}
}

// resolveIdentifier 为 posts / portfolio-items / tools 的 slug 与 diaries 的 publicId 保证非空：
// 创建时为空 → 按标题/名称自动生成唯一标识；更新时为空 → 沿用既有记录的值，避免改写公开 URL。
func resolveIdentifier(tx *gorm.DB, resource, id string, item interface{}, create bool) error {
	value, ok := identifierOf(item)
	if !ok || value != "" {
		return nil
	}

	var table, field, label, title string
	switch v := item.(type) {
	case *blogmodel.Post:
		table, field, label, title = "td_blog_post", "slug", "post", v.Title
	case *blogmodel.Diary:
		table, field, label, title = "td_blog_diary", "public_id", "diary", v.Title
	case *blogmodel.PortfolioItem:
		table, field, label, title = "td_blog_portfolio_item", "slug", "portfolio", v.Title
	case *blogmodel.Tool:
		table, field, label, title = "td_blog_tool", "slug", "tool", v.Name
	}

	if !create {
		// 更新：沿用现有值，避免 URL 变化
		var existing string
		if err := tx.Table(table).Where("id = ?", id).Pluck(field, &existing).Error; err != nil {
			return err
		}
		if existing != "" {
			setIdentifier(item, existing)
		}
		return nil
	}

	// 创建：按标题生成，冲突时追加 -2 / -3 …
	base := slugify(title)
	if base == "" {
		base = label
	}
	candidate := base
	for suffix := 2; ; suffix++ {
		var count int64
		if err := tx.Table(table).Where(field+" = ? AND deleted_at IS NULL", candidate).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			break
		}
		candidate = fmt.Sprintf("%s-%d", base, suffix)
	}
	setIdentifier(item, candidate)
	return nil
}

// setIdentifier 写回生成的 slug / publicId 到模型。
func setIdentifier(item interface{}, value string) {
	switch v := item.(type) {
	case *blogmodel.Post:
		v.Slug = value
	case *blogmodel.Diary:
		v.PublicID = value
	case *blogmodel.PortfolioItem:
		v.Slug = value
	case *blogmodel.Tool:
		v.Slug = value
	}
}

// slugify 把标题转成 URL 友好的标识：保留字母/数字/中文，其余转连字符并折叠。
func slugify(s string) string {
	b := strings.Builder{}
	prevDash := false
	for _, r := range strings.ToLower(strings.TrimSpace(s)) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || (r >= '\u4e00' && r <= '\u9fff') {
			b.WriteRune(r)
			prevDash = false
			continue
		}
		if !prevDash && b.Len() > 0 {
			b.WriteByte('-')
			prevDash = true
		}
	}
	return strings.Trim(b.String(), "-")
}

func buildModel(resource string, payload interface{}) (interface{}, []string, string, error) {
	switch value := payload.(type) {
	case request.Category:
		if !request.ValidScope(value.Scope) {
			return nil, nil, "", request.Field("scope", "分类归属只能是 post / portfolio / bookmark")
		}
		if !request.ValidID(value.Slug) {
			return nil, nil, "", request.Field("slug", "Slug 不能为空，且长度不超过 128")
		}
		if strings.TrimSpace(value.Name) == "" {
			return nil, nil, "", request.Field("name", "名称不能为空")
		}
		enabled := true
		if value.Enabled != nil {
			enabled = *value.Enabled
		}
		return &blogmodel.Category{Scope: value.Scope, Slug: value.Slug, Name: value.Name, SortOrder: value.SortOrder, Enabled: enabled}, nil, "", nil
	case request.Tag:
		name := strings.TrimSpace(value.Name)
		if name == "" {
			return nil, nil, "", request.Field("name", "名称不能为空")
		}
		return &blogmodel.Tag{Name: name, NormalizedName: strings.ToLower(name)}, nil, "", nil
	case request.Post:
		if value.Slug != "" && !request.ValidID(value.Slug) {
			return nil, nil, "", request.Field("slug", "Slug 不能为空，且长度不超过 128")
		}
		if strings.TrimSpace(value.Title) == "" {
			return nil, nil, "", request.Field("title", "标题不能为空")
		}
		if !request.ValidStatus(value.PublishStatus) {
			return nil, nil, "", request.Field("publishStatus", "发布状态只能是 draft / published / archived")
		}
		on, at, err := publishedTimes(value.PublishStatus, value.PublishedOn, value.PublishedAt)
		if err != nil {
			return nil, nil, "", err
		}
		ids, err2 := request.NormalizeIDs(value.TagIDs)
		if err2 != nil {
			return nil, nil, "", request.Field("tagIds", "标签 ID 不合法")
		}
		return &blogmodel.Post{Slug: value.Slug, Title: value.Title, Summary: value.Summary, CategoryID: value.CategoryID, Author: value.Author, Content: value.Content, PublishStatus: value.PublishStatus, PublishedOn: on, PublishedAt: at, Pinned: value.Pinned, Version: max(value.Version, 1), PublishSource: value.PublishSource, SourceID: value.SourceID}, ids, "td_blog_post_tag", nil
	case request.Diary:
		if value.PublicID != "" && !request.ValidID(value.PublicID) {
			return nil, nil, "", request.Field("publicId", "公开 ID 不能为空，且长度不超过 128")
		}
		if strings.TrimSpace(value.Title) == "" {
			return nil, nil, "", request.Field("title", "标题不能为空")
		}
		if !request.ValidStatus(value.PublishStatus) {
			return nil, nil, "", request.Field("publishStatus", "发布状态只能是 draft / published / archived")
		}
		on, at, err := publishedTimes(value.PublishStatus, value.PublishedOn, value.PublishedAt)
		if err != nil {
			return nil, nil, "", err
		}
		ids, err2 := request.NormalizeIDs(value.TagIDs)
		if err2 != nil {
			return nil, nil, "", request.Field("tagIds", "标签 ID 不合法")
		}
		return &blogmodel.Diary{PublicID: value.PublicID, Title: value.Title, Summary: value.Summary, Content: value.Content, Mood: value.Mood, Weather: value.Weather, PublishStatus: value.PublishStatus, PublishedOn: on, PublishedAt: at, Pinned: value.Pinned, Version: max(value.Version, 1), PublishSource: value.PublishSource, SourceID: value.SourceID}, ids, "td_blog_diary_tag", nil
	case request.Portfolio:
		if value.Slug != "" && !request.ValidID(value.Slug) {
			return nil, nil, "", request.Field("slug", "Slug 长度不能超过 128")
		}
		if strings.TrimSpace(value.Title) == "" {
			return nil, nil, "", request.Field("title", "标题不能为空")
		}
		if !request.ValidStatus(value.PublishStatus) {
			return nil, nil, "", request.Field("publishStatus", "发布状态只能是 draft / published / archived")
		}
		if value.DemoURL != "" {
			demoURL, ok := request.NormalizeLinkURL(value.DemoURL)
			if !ok {
				return nil, nil, "", request.Field("demoUrl", "演示地址必须是合法链接（支持 http/https 等协议），或已上传的 /files 路径")
			}
			value.DemoURL = demoURL
		}
		if value.RepoURL != "" {
			repoURL, ok := request.NormalizeLinkURL(value.RepoURL)
			if !ok {
				return nil, nil, "", request.Field("repoUrl", "仓库地址必须是合法链接（支持 http/https 等协议），或已上传的 /files 路径")
			}
			value.RepoURL = repoURL
		}
		for _, mediaURL := range value.Gallery {
			if !validMediaURL(mediaURL) {
				return nil, nil, "", request.Field("gallery", "图集中的地址必须是以 https:// 开头的地址，或已上传的 /files 路径")
			}
		}
		on, at, err := publishedTimes(value.PublishStatus, value.PublishedOn, value.PublishedAt)
		if err != nil {
			return nil, nil, "", err
		}
		tech, e1 := marshal(value.TechStack)
		links, e2 := marshal(value.Links)
		gallery, e3 := marshal(value.Gallery)
		metrics, e4 := marshal(value.Metrics)
		if e1 != nil || e2 != nil || e3 != nil || e4 != nil {
			return nil, nil, "", ErrInvalid
		}
		return &blogmodel.PortfolioItem{Slug: value.Slug, Title: value.Title, Summary: value.Summary, CategoryID: value.CategoryID, Cover: value.Cover, TechStack: tech, Links: links, Gallery: gallery, Metrics: metrics, DemoURL: value.DemoURL, RepoURL: value.RepoURL, Status: value.Status, Role: value.Role, Content: value.Content, PublishStatus: value.PublishStatus, PublishedOn: on, PublishedAt: at, Version: max(value.Version, 1)}, nil, "", nil
	case request.Tool:
		if err := request.ValidateTool(value); err != nil {
			return nil, nil, "", mapToolValidationError(err)
		}
		_, at, _ := publishedTimes(value.PublishStatus, "", value.PublishedAt)
		if value.Kind == "link" {
			// 校验已通过，这里只做归一化（缺协议时补 https://）
			value.URL, _ = request.NormalizeLinkURL(value.URL)
			value.Cover = ""
			value.DevelopmentStatus = ""
			value.Content = ""
		} else {
			value.URL = ""
		}
		return &blogmodel.Tool{Slug: value.Slug, Kind: value.Kind, Name: value.Name, Description: value.Description, URL: value.URL, Cover: value.Cover, DevelopmentStatus: value.DevelopmentStatus, Content: value.Content, PublishStatus: value.PublishStatus, PublishedAt: at, SortOrder: value.SortOrder, Version: max(value.Version, 1)}, nil, "", nil
	case request.Bookmark:
		if strings.TrimSpace(value.Title) == "" {
			return nil, nil, "", request.Field("title", "标题不能为空")
		}
		url, ok := request.NormalizeLinkURL(value.URL)
		if !ok {
			return nil, nil, "", request.Field("url", "地址不能为空，且必须是合法链接（支持 http/https 等协议，不支持 javascript: 这类地址）")
		}
		value.URL = url
		if !request.ValidStatus(value.PublishStatus) {
			return nil, nil, "", request.Field("publishStatus", "发布状态只能是 draft / published / archived")
		}
		_, at, _ := publishedTimes(value.PublishStatus, "", value.PublishedAt)
		ids, err := request.NormalizeIDs(value.TagIDs)
		if err != nil {
			return nil, nil, "", request.Field("tagIds", "标签 ID 不合法")
		}
		return &blogmodel.Bookmark{Title: value.Title, URL: value.URL, Description: value.Description, CategoryID: value.CategoryID, Icon: value.Icon, PublishStatus: value.PublishStatus, PublishedAt: at, SortOrder: value.SortOrder, OpenInNewTab: value.OpenInNewTab, Version: max(value.Version, 1)}, ids, "td_blog_bookmark_tag", nil
	default:
		return nil, nil, "", fmt.Errorf("%w: %s", ErrInvalid, resource)
	}
}

func updateMap(item interface{}) map[string]interface{} {
	switch value := item.(type) {
	case *blogmodel.Tag:
		return map[string]interface{}{"name": value.Name, "normalized_name": value.NormalizedName}
	case *blogmodel.Post:
		return map[string]interface{}{"slug": value.Slug, "title": value.Title, "summary": value.Summary, "category_id": value.CategoryID, "author": value.Author, "content": value.Content, "publish_status": value.PublishStatus, "published_on": value.PublishedOn, "published_at": value.PublishedAt, "pinned": value.Pinned, "publish_source": value.PublishSource, "source_id": value.SourceID}
	case *blogmodel.Diary:
		return map[string]interface{}{"public_id": value.PublicID, "title": value.Title, "summary": value.Summary, "content": value.Content, "mood": value.Mood, "weather": value.Weather, "publish_status": value.PublishStatus, "published_on": value.PublishedOn, "published_at": value.PublishedAt, "pinned": value.Pinned, "publish_source": value.PublishSource, "source_id": value.SourceID}
	case *blogmodel.PortfolioItem:
		return map[string]interface{}{"slug": value.Slug, "title": value.Title, "summary": value.Summary, "category_id": value.CategoryID, "cover": value.Cover, "tech_stack": value.TechStack, "links": value.Links, "gallery": value.Gallery, "metrics": value.Metrics, "demo_url": value.DemoURL, "repo_url": value.RepoURL, "status": value.Status, "role": value.Role, "content": value.Content, "publish_status": value.PublishStatus, "published_on": value.PublishedOn, "published_at": value.PublishedAt}
	case *blogmodel.Tool:
		return map[string]interface{}{"slug": value.Slug, "kind": value.Kind, "name": value.Name, "description": value.Description, "url": value.URL, "cover": value.Cover, "development_status": value.DevelopmentStatus, "content": value.Content, "publish_status": value.PublishStatus, "published_at": value.PublishedAt, "sort_order": value.SortOrder}
	case *blogmodel.Bookmark:
		return map[string]interface{}{"title": value.Title, "url": value.URL, "description": value.Description, "category_id": value.CategoryID, "icon": value.Icon, "publish_status": value.PublishStatus, "published_at": value.PublishedAt, "sort_order": value.SortOrder, "open_in_new_tab": value.OpenInNewTab}
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
	if err := tx.Model(&current).Updates(map[string]interface{}{"slug": next.Slug, "name": next.Name, "sort_order": next.SortOrder, "enabled": next.Enabled}).Error; err != nil {
		return mapDBError(err)
	}
	return nil
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
			// Name 是必填（validateProfile），默认带上与 seed 一致的站长名，避免默认对象存不回去
			return blogresponse.Profile{Name: defaultProfileName, Links: []blogresponse.ProfileLink{}, Skills: []blogresponse.ProfileSkill{}}, nil
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
		return defaultSite(), nil
	}
	if err != nil {
		return nil, err
	}
	return siteFromModel(value)
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
	modules, err := normalizeSiteModules(value.Modules, true)
	if err != nil {
		return nil, err
	}
	value.Modules = modules
	if value.TechStack == nil {
		value.TechStack = []string{}
	}
	if value.Milestones == nil {
		value.Milestones = []blogresponse.SiteMilestone{}
	}
	tech, e1 := marshal(value.TechStack)
	modulesJSON, e2 := marshal(value.Modules)
	milestones, e3 := marshal(value.Milestones)
	home, e4 := marshal(value.Home)
	footer, e5 := marshal(value.Footer)
	banner, e6 := marshal(value.Banner)
	if e1 != nil || e2 != nil || e3 != nil || e4 != nil || e5 != nil || e6 != nil {
		return nil, ErrInvalid
	}
	item := &blogmodel.Site{SiteKey: "default", Name: value.Name, Slogan: value.Slogan, Intro: value.Intro, TechStack: tech, Modules: modulesJSON, Milestones: milestones, Home: home, Footer: footer, Banner: banner, MaintenanceMode: value.MaintenanceMode, MemoPublicEnabled: value.MemoPublicEnabled}
	// 必须显式指定列：GORM 用结构体更新会跳过零值字段，导致关闭维护模式（false）写不进库
	if err = upsertSettingColumns(db, &blogmodel.Site{}, "site_key", "default", item, siteUpdateColumns); err != nil {
		return nil, err
	}
	return value, nil
}

func validHexColor(value string) bool {
	if len(value) != 7 || value[0] != '#' {
		return false
	}
	for _, char := range value[1:] {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F')) {
			return false
		}
	}
	return true
}

// mapToolValidationError 把 request 层的工具校验失败细化为可提示的具体原因。
// 字段级错误（*request.FieldError）原样上抛，由 API 层带出字段名与原因。
func mapToolValidationError(err error) error {
	switch {
	case errors.Is(err, request.ErrToolURLRequired):
		return ErrToolURLRequired
	case errors.Is(err, request.ErrToolStatusRequired):
		return ErrToolStatusRequired
	}
	var fieldErr *request.FieldError
	if errors.As(err, &fieldErr) {
		return err
	}
	return ErrInvalid
}

// maxProfileContactLength 限制单条联系方式的长度：字段允许任意文本，只做长度保护。
const maxProfileContactLength = 500

func validateProfile(value blogresponse.Profile) error {
	if strings.TrimSpace(value.Name) == "" {
		return ErrInvalid
	}
	seen := map[string]struct{}{}
	for _, link := range value.Links {
		if strings.TrimSpace(link.Label) == "" {
			return ErrInvalid
		}
		if _, ok := seen[link.Label]; ok {
			return ErrInvalid
		}
		seen[link.Label] = struct{}{}
		// 联系方式不再限制协议：手机号、QQ 号、微信号等纯文本同样合法
		if utf8.RuneCountInString(link.URL) > maxProfileContactLength {
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

// siteFromModel 把站点记录转换为响应对象，并按固定模块约束补齐历史数据。
func siteFromModel(value *blogmodel.Site) (blogresponse.Site, error) {
	tech := []string{}
	modules := []blogresponse.SiteModule{}
	milestones := []blogresponse.SiteMilestone{}
	if json.Unmarshal(value.TechStack, &tech) != nil || json.Unmarshal(value.Modules, &modules) != nil || json.Unmarshal(value.Milestones, &milestones) != nil {
		return blogresponse.Site{}, ErrInvalid
	}
	normalized, err := normalizeSiteModules(modules, false)
	if err != nil {
		return blogresponse.Site{}, err
	}
	home := blogresponse.DefaultSiteHome()
	footer := blogresponse.DefaultSiteFooter()
	if len(value.Home) > 0 && string(value.Home) != "{}" && json.Unmarshal(value.Home, &home) != nil {
		return blogresponse.Site{}, ErrInvalid
	}
	home.Subtitle = strings.TrimSpace(home.Subtitle)
	if utf8.RuneCountInString(home.Subtitle) > blogresponse.MaxSiteHomeSubtitleLength {
		home.Subtitle = blogresponse.DefaultSiteHome().Subtitle
	}
	if len(value.Footer) > 0 && string(value.Footer) != "{}" && json.Unmarshal(value.Footer, &footer) != nil {
		return blogresponse.Site{}, ErrInvalid
	}
	banner := blogresponse.DefaultSiteBanner()
	if len(value.Banner) > 0 && string(value.Banner) != "{}" && json.Unmarshal(value.Banner, &banner) != nil {
		return blogresponse.Site{}, ErrInvalid
	}
	return blogresponse.Site{Name: value.Name, Slogan: value.Slogan, Intro: value.Intro, TechStack: tech, Modules: normalized, Milestones: milestones, Home: home, Footer: footer, Banner: banner, MaintenanceMode: value.MaintenanceMode, MemoPublicEnabled: value.MemoPublicEnabled}, nil
}

func validateSite(value blogresponse.Site) error {
	if strings.TrimSpace(value.Name) == "" {
		return ErrInvalid
	}
	if utf8.RuneCountInString(strings.TrimSpace(value.Banner.Text)) > 100 || !validHexColor(value.Banner.BackgroundColor) || !validHexColor(value.Banner.TextColor) {
		return ErrInvalid
	}
	for _, milestone := range value.Milestones {
		if !request.ValidDate(milestone.Date) || strings.TrimSpace(milestone.Title) == "" {
			return ErrInvalid
		}
	}
	if strings.TrimSpace(value.Home.Title) == "" || strings.TrimSpace(value.Home.AI.Title) == "" || strings.TrimSpace(value.Home.AI.Description) == "" || utf8.RuneCountInString(strings.TrimSpace(value.Home.Subtitle)) > blogresponse.MaxSiteHomeSubtitleLength {
		return ErrInvalid
	}
	if err := validateHomeTerminal(value.Home.Terminal); err != nil {
		return ErrInvalid
	}
	if _, err := normalizeSiteModules(value.Modules, true); err != nil {
		return ErrInvalid
	}
	// 链接类字段：不限协议（http/https/mailto/自定义均可），只拦可执行脚本的协议
	for _, url := range []string{value.Home.AI.LinkURL, value.Footer.LinkURL, value.Footer.ICPURL, value.Footer.PoliceURL} {
		if url != "" && !validSiteLink(url) {
			return ErrInvalid
		}
	}
	// 图片类字段：浏览器以子资源加载，http 在 https 站点会被混合内容拦截，仍只收 /files 路径或 https
	for _, url := range []string{value.Home.AI.ImageURL, value.Home.PortfolioImageURL, value.Home.BookmarkImageURL} {
		if url != "" && !validSiteImage(url) {
			return ErrInvalid
		}
	}
	return nil
}

// validSiteLink 站点设置里的链接：允许站内路径与未限协议的绝对地址（javascript: 等由 NormalizeLinkURL 拒绝）。
func validSiteLink(value string) bool {
	if strings.HasPrefix(value, "/files/") || value == "/Blog" || strings.HasPrefix(value, "/Blog/") {
		return true
	}
	_, ok := request.NormalizeLinkURL(value)
	return ok
}

// validSiteImage 站点设置里的图片：只接受已上传的 /files 路径或 https 绝对地址。
func validSiteImage(value string) bool {
	return strings.HasPrefix(value, "/files/") || request.ValidURL(value, false)
}

// validMediaURL 允许留空、已上传的 /files 路径，或 https 绝对地址。
func validMediaURL(value string) bool {
	return value == "" || strings.HasPrefix(value, "/files/") || request.ValidURL(value, false)
}

// validateHomeTerminal 校验主页终端卡片：命令必填，标题可选，输出行非空且数量受限。
func validateHomeTerminal(terminal blogresponse.SiteHomeTerminal) error {
	command := strings.TrimSpace(terminal.Command)
	if command == "" || utf8.RuneCountInString(command) > blogresponse.MaxSiteHomeTerminalCommandLength {
		return ErrInvalid
	}
	if utf8.RuneCountInString(strings.TrimSpace(terminal.Title)) > blogresponse.MaxSiteHomeTerminalTitleLength {
		return ErrInvalid
	}
	if len(terminal.Lines) > blogresponse.MaxSiteHomeTerminalLines {
		return ErrInvalid
	}
	for _, line := range terminal.Lines {
		line = strings.TrimSpace(line)
		if line == "" || utf8.RuneCountInString(line) > blogresponse.MaxSiteHomeTerminalLineLength {
			return ErrInvalid
		}
	}
	return nil
}

// siteUpdateColumns 站点配置允许写入的列，包含开关类字段（零值也需要写入）
var siteUpdateColumns = []string{"name", "slogan", "intro", "tech_stack", "modules", "milestones", "home", "footer", "banner", "maintenance_mode", "memo_public_enabled"}

// upsertSettingColumns 与 upsertSetting 相同，但只更新指定列且允许写入零值。
func upsertSettingColumns(db *gorm.DB, existing interface{}, keyColumn, keyValue string, values interface{}, columns []string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		result := tx.Unscoped().Where(keyColumn+" = ?", keyValue).First(existing)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return tx.Create(values).Error
		}
		if result.Error != nil {
			return result.Error
		}
		if err := tx.Unscoped().Model(existing).Select(columns).Updates(values).Error; err != nil {
			return err
		}
		// 已软删除的设置记录重新启用
		return tx.Unscoped().Model(existing).Update("deleted_at", nil).Error
	})
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
