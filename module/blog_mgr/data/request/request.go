package request

import (
	"errors"
	"net/url"
	"strings"
	"time"

	blogmodel "fastduck/treasure-doc/module/blog/data/model"
	blogresponse "fastduck/treasure-doc/module/blog/data/response"
)

var ErrInvalid = errors.New("invalid request")

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
}
type Portfolio struct {
	Slug          string                       `json:"slug"`
	Title         string                       `json:"title"`
	Summary       string                       `json:"summary"`
	CategoryID    string                       `json:"categoryId"`
	Cover         string                       `json:"cover"`
	TechStack     []string                     `json:"techStack"`
	Links         []blogresponse.PortfolioLink `json:"links"`
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
	PublicID      string     `json:"publicId"`
	Title         string     `json:"title"`
	URL           string     `json:"url"`
	Description   string     `json:"description"`
	CategoryID    string     `json:"categoryId"`
	Icon          string     `json:"icon"`
	PublishStatus string     `json:"publishStatus"`
	PublishedAt   *time.Time `json:"publishedAt"`
	SortOrder     int        `json:"sortOrder"`
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

func ValidURL(value string, allowMail bool) bool {
	parsed, err := url.ParseRequestURI(value)
	if err != nil {
		return false
	}
	return parsed.Scheme == "https" || allowMail && parsed.Scheme == "mailto"
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
	if !ValidID(value.Slug) || strings.TrimSpace(value.Name) == "" || !ValidStatus(value.PublishStatus) {
		return ErrInvalid
	}
	if value.Kind == "link" {
		if !ValidURL(value.URL, false) {
			return ErrInvalid
		}
		return nil
	}
	if value.Kind != "own" || strings.TrimSpace(value.DevelopmentStatus) == "" {
		return ErrInvalid
	}
	return nil
}
