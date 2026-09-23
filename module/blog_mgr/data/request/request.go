package request

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
	"time"

	blogmodel "fastduck/treasure-doc/module/blog/data/model"
	blogresponse "fastduck/treasure-doc/module/blog/data/response"
)

var ErrInvalid = errors.New("invalid request")

// 工具校验的细分原因：让管理端能直接提示「缺哪个字段」，而不是笼统的参数格式错误。
var (
	ErrToolURLRequired    = errors.New("link tool requires a valid url")
	ErrToolStatusRequired = errors.New("own tool requires development status")
)

// FieldError 指向具体字段的校验失败。
//
// 响应里会同时带出字段名（data.field）与可直接展示给用户的原因（msg），
// 让前端能提示「哪个字段、该怎么改」，而不是笼统的「请求参数格式错误」。
type FieldError struct {
	Field  string
	Reason string
}

func (e *FieldError) Error() string { return e.Reason }

// Field 构造字段级校验错误；reason 会直接展示给用户，需写明字段与具体要求。
func Field(field, reason string) *FieldError {
	return &FieldError{Field: field, Reason: reason}
}

type List struct {
	Page       int    `form:"page"`
	PageSize   int    `form:"pageSize"`
	Keyword    string `form:"keyword"`
	Status     string `form:"status"`
	Deleted    string `form:"deleted"`
	Scope      string `form:"scope"`
	CategoryID string `form:"categoryId"`
	Sort       string `form:"sort"`
}

func (q *List) Normalize() error {
	if q.Page == 0 {
		q.Page = 1
	}
	if q.PageSize == 0 {
		q.PageSize = 20
	}
	if q.Deleted == "" {
		q.Deleted = "exclude"
	}
	if q.Sort == "" {
		q.Sort = "desc"
	}
	q.Keyword = strings.TrimSpace(q.Keyword)
	q.Status = strings.TrimSpace(q.Status)
	q.Scope = strings.TrimSpace(q.Scope)
	q.CategoryID = strings.TrimSpace(q.CategoryID)
	if q.Page < 1 || q.PageSize < 1 || q.PageSize > 100 {
		return ErrInvalid
	}
	if q.Sort != "asc" && q.Sort != "desc" {
		return ErrInvalid
	}
	if q.Deleted != "exclude" && q.Deleted != "only" && q.Deleted != "all" {
		return ErrInvalid
	}
	if q.Status != "" && !ValidStatus(q.Status) {
		return ErrInvalid
	}
	if q.Scope != "" && !ValidScope(q.Scope) {
		return ErrInvalid
	}
	return nil
}

func (q List) Offset() int { return (q.Page - 1) * q.PageSize }

type Category struct {
	Scope     string `json:"scope"`
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	SortOrder int    `json:"sortOrder"`
	Enabled   *bool  `json:"enabled"`
}
type Tag struct {
	Name string `json:"name"`
}
type Post struct {
	Slug          string     `json:"slug"`
	Title         string     `json:"title"`
	Summary       string     `json:"summary"`
	CategoryID    string     `json:"categoryId"`
	Author        string     `json:"author"`
	Content       string     `json:"content"`
	PublishStatus string     `json:"publishStatus"`
	PublishedOn   string     `json:"publishedOn"`
	PublishedAt   *time.Time `json:"publishedAt"`
	Pinned        bool       `json:"pinned"`
	Version       int        `json:"version"`
	TagIDs        []string   `json:"tagIds"`
	// Overwrite 仅发布端使用：true 时按标题推导 slug、服务端校验已发布后覆盖，否则新建。
	Overwrite bool `json:"overwrite"`
	// 发布来源 + 该来源内的稳定标识（思源端传 siyuan + 文档 id），用于覆盖更新精确定位
	PublishSource string `json:"publishSource"`
	SourceID      string `json:"sourceId"`
}
type Diary struct {
	PublicID      string     `json:"publicId"`
	Title         string     `json:"title"`
	Summary       string     `json:"summary"`
	Content       string     `json:"content"`
	Mood          string     `json:"mood"`
	Weather       string     `json:"weather"`
	PublishStatus string     `json:"publishStatus"`
	PublishedOn   string     `json:"publishedOn"`
	PublishedAt   *time.Time `json:"publishedAt"`
	Pinned        bool       `json:"pinned"`
	Version       int        `json:"version"`
	TagIDs        []string   `json:"tagIds"`
	// Overwrite 仅发布端使用：true 时按标题推导 slug、服务端校验已发布后覆盖，否则新建。
	Overwrite bool `json:"overwrite"`
	// 发布来源 + 该来源内的稳定标识（思源端传 siyuan + 文档 id），用于覆盖更新精确定位
	PublishSource string `json:"publishSource"`
	SourceID      string `json:"sourceId"`
}
type Portfolio struct {
	Slug          string                       `json:"slug"`
	Title         string                       `json:"title"`
	Summary       string                       `json:"summary"`
	CategoryID    string                       `json:"categoryId"`
	Cover         string                       `json:"cover"`
	TechStack     []string                     `json:"techStack"`
	Links         []blogresponse.PortfolioLink `json:"links"`
	Gallery       []string                     `json:"gallery"`
	Metrics       []string                     `json:"metrics"`
	DemoURL       string                       `json:"demoUrl"`
	RepoURL       string                       `json:"repoUrl"`
	Status        string                       `json:"status"`
	Role          string                       `json:"role"`
	Content       string                       `json:"content"`
	PublishStatus string                       `json:"publishStatus"`
	PublishedOn   string                       `json:"publishedOn"`
	PublishedAt   *time.Time                   `json:"publishedAt"`
	Version       int                          `json:"version"`
}
type Tool struct {
	Slug              string     `json:"slug"`
	Kind              string     `json:"kind"`
	Name              string     `json:"name"`
	Description       string     `json:"description"`
	URL               string     `json:"url"`
	Cover             string     `json:"cover"`
	DevelopmentStatus string     `json:"developmentStatus"`
	Content           string     `json:"content"`
	PublishStatus     string     `json:"publishStatus"`
	PublishedAt       *time.Time `json:"publishedAt"`
	SortOrder         int        `json:"sortOrder"`
	Version           int        `json:"version"`
}
type Bookmark struct {
	Title         string     `json:"title"`
	URL           string     `json:"url"`
	Description   string     `json:"description"`
	CategoryID    string     `json:"categoryId"`
	Icon          string     `json:"icon"`
	PublishStatus string     `json:"publishStatus"`
	PublishedAt   *time.Time `json:"publishedAt"`
	SortOrder     int        `json:"sortOrder"`
	OpenInNewTab  bool       `json:"openInNewTab"`
	Version       int        `json:"version"`
	TagIDs        []string   `json:"tagIds"`
}
type Profile = blogresponse.Profile
type Site = blogresponse.Site

func ValidStatus(value string) bool {
	return value == blogmodel.StatusDraft || value == blogmodel.StatusPublished || value == blogmodel.StatusArchived
}
func ValidScope(value string) bool {
	return value == blogmodel.CategoryPost || value == blogmodel.CategoryPortfolio || value == blogmodel.CategoryBookmark
}
func ValidID(value string) bool {
	value = strings.TrimSpace(value)
	return len(value) > 0 && len(value) <= 128
}
func ValidDate(value string) bool { _, err := time.Parse("2006-01-02", value); return err == nil }

// MaxBatchIDs 限制单次批量操作的 ID 数量，避免生成超长 IN 语句。
const MaxBatchIDs = 200

// BatchIDs 批量操作请求体，例如批量删除：{"ids": ["1", "2"]}。
type BatchIDs struct {
	IDs []string `json:"ids"`
}

// Normalize 去重并校验 ID；空集合或超过上限时返回 ErrInvalid。
func (b *BatchIDs) Normalize() error {
	ids := make([]string, 0, len(b.IDs))
	seen := make(map[string]struct{}, len(b.IDs))
	for _, id := range b.IDs {
		id = strings.TrimSpace(id)
		if !ValidID(id) {
			return ErrInvalid
		}
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	if len(ids) == 0 || len(ids) > MaxBatchIDs {
		return ErrInvalid
	}
	b.IDs = ids
	return nil
}

func ValidURL(value string, allowMail bool) bool {
	parsed, err := url.ParseRequestURI(value)
	if err != nil {
		return false
	}
	return parsed.Scheme == "https" || allowMail && parsed.Scheme == "mailto"
}

// 这些协议渲染成 <a href> 后会在访客浏览器里执行脚本，必须拒绝；
// 其余协议（http / https / mailto / ftp / 自定义等）不做限制。
var unsafeSchemes = map[string]struct{}{
	"javascript": {},
	"data":       {},
	"vbscript":   {},
	"file":       {},
}

var schemePrefix = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9+.-]*:`)

// NormalizeLinkURL 归一化「地址」类字段（收藏集/利器的地址、作品演示与仓库地址、站点链接等）：
// 不限制协议；没写协议时按 https 补全（`example.com` → `https://example.com`）；
// 拒绝 javascript: / data: 这类可执行脚本的协议，以及无法解析的地址。
func NormalizeLinkURL(value string) (string, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", false
	}
	if !schemePrefix.MatchString(value) && !strings.HasPrefix(value, "/") {
		value = "https://" + value
	}
	parsed, err := url.ParseRequestURI(value)
	if err != nil {
		return "", false
	}
	if _, unsafe := unsafeSchemes[strings.ToLower(parsed.Scheme)]; unsafe {
		return "", false
	}
	return value, true
}

func NormalizeIDs(values []string) ([]string, error) {
	result := make([]string, 0, len(values))
	seen := map[string]struct{}{}
	for _, value := range values {
		value = strings.TrimSpace(value)
		if !ValidID(value) {
			return nil, ErrInvalid
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result, nil
}

func ValidateTool(value Tool) error {
	if value.Slug != "" && !ValidID(value.Slug) {
		return Field("slug", "Slug 长度不能超过 128")
	}
	if strings.TrimSpace(value.Name) == "" {
		return Field("name", "名称不能为空")
	}
	if !ValidStatus(value.PublishStatus) {
		return Field("publishStatus", "发布状态只能是 draft / published / archived")
	}
	switch value.Kind {
	case "link":
		if _, ok := NormalizeLinkURL(value.URL); !ok {
			return ErrToolURLRequired
		}
		return nil
	case "own":
		if strings.TrimSpace(value.DevelopmentStatus) == "" {
			return ErrToolStatusRequired
		}
		return nil
	default:
		return Field("kind", "类型只能是「自研工具」或「外部链接」")
	}
}
