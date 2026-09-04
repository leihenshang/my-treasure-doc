package request

import (
	"errors"
	"strings"
)

const (
	DefaultPage     = 1
	DefaultPageSize = 10
	MaxPageSize     = 100
)

var (
	ErrInvalidQuery = errors.New("invalid query parameters")
	ErrInvalidSort  = errors.New("unsupported sort order")
)

type PageQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
	Sort     string `form:"sort"`
}

func (q *PageQuery) Normalize() error {
	if q.Page == 0 {
		q.Page = DefaultPage
	}
	if q.PageSize == 0 {
		q.PageSize = DefaultPageSize
	}
	if q.Sort == "" {
		q.Sort = "desc"
	}

	if q.Page < 1 || q.PageSize < 1 || q.PageSize > MaxPageSize {
		return ErrInvalidQuery
	}
	if q.Sort != "asc" && q.Sort != "desc" {
		return ErrInvalidSort
	}
	return nil
}

func (q PageQuery) Offset() int {
	return (q.Page - 1) * q.PageSize
}

type PostQuery struct {
	PageQuery
	CategoryID string `form:"categoryId"`
	Tag        string `form:"tag"`
	Keyword    string `form:"keyword"`
}

func (q *PostQuery) Normalize() error {
	q.CategoryID = strings.TrimSpace(q.CategoryID)
	q.Tag = strings.TrimSpace(q.Tag)
	q.Keyword = strings.TrimSpace(q.Keyword)
	return q.PageQuery.Normalize()
}

type DiaryQuery struct {
	PageQuery
	Tag     string `form:"tag"`
	Keyword string `form:"keyword"`
}

func (q *DiaryQuery) Normalize() error {
	q.Tag = strings.TrimSpace(q.Tag)
	q.Keyword = strings.TrimSpace(q.Keyword)
	return q.PageQuery.Normalize()
}

type PortfolioQuery struct {
	CategoryID string `form:"categoryId"`
}

func (q *PortfolioQuery) Normalize() {
	q.CategoryID = strings.TrimSpace(q.CategoryID)
}

type BookmarkQuery struct {
	CategoryID string `form:"categoryId"`
	Keyword    string `form:"keyword"`
}

func (q *BookmarkQuery) Normalize() {
	q.CategoryID = strings.TrimSpace(q.CategoryID)
	q.Keyword = strings.TrimSpace(q.Keyword)
}

func EscapeLike(value string) string {
	value = strings.ReplaceAll(value, `\`, `\\`)
	value = strings.ReplaceAll(value, `%`, `\%`)
	return strings.ReplaceAll(value, `_`, `\_`)
}

func ValidPublicID(value string) bool {
	length := len(value)
	return length >= 1 && length <= 128
}
