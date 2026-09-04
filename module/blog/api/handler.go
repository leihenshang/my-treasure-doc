package api

import (
	"context"
	"errors"
	"net/http"

	"fastduck/treasure-doc/module/blog/data/model"
	"fastduck/treasure-doc/module/blog/data/request"
	"fastduck/treasure-doc/module/blog/data/response"
	"fastduck/treasure-doc/module/blog/internal/service"

	"github.com/gin-gonic/gin"
)

type PublicService interface {
	Categories(context.Context, string) ([]response.Category, error)
	PostTags(context.Context) ([]string, error)
	ListPosts(context.Context, request.PostQuery) (response.Page, error)
	GetPost(context.Context, string) (response.Post, error)
	DiaryTags(context.Context) ([]string, error)
	ListDiaries(context.Context, request.DiaryQuery) (response.Page, error)
	GetDiary(context.Context, string) (response.Diary, error)
	ListPortfolio(context.Context, request.PortfolioQuery) ([]response.PortfolioSummary, error)
	GetPortfolio(context.Context, string) (response.PortfolioItem, error)
	ListTools(context.Context) ([]response.Tool, error)
	GetTool(context.Context, string) (response.Tool, error)
	ListBookmarks(context.Context, request.BookmarkQuery) ([]response.Bookmark, error)
	Profile(context.Context) (response.Profile, error)
	Site(context.Context) (response.Site, error)
	Stats(context.Context) (response.Stats, error)
}

type Handler struct {
	service PublicService
}

func NewHandler(publicService PublicService) *Handler {
	return &Handler{service: publicService}
}

func (h *Handler) BlogCategories(c *gin.Context) {
	h.categories(c, model.CategoryPost)
}

func (h *Handler) BlogTags(c *gin.Context) {
	data, err := h.service.PostTags(c.Request.Context())
	h.write(c, data, err)
}

func (h *Handler) BlogPosts(c *gin.Context) {
	var query request.PostQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		invalidQuery(c)
		return
	}
	if err := query.Normalize(); err != nil {
		writeQueryError(c, err)
		return
	}
	data, err := h.service.ListPosts(c.Request.Context(), query)
	h.write(c, data, err)
}

func (h *Handler) BlogPost(c *gin.Context) {
	h.detail(c, service.ErrPostNotFound, response.CodePostNotFound, "文章不存在或已被删除", func(id string) (interface{}, error) {
		return h.service.GetPost(c.Request.Context(), id)
	})
}

func (h *Handler) DiaryTags(c *gin.Context) {
	data, err := h.service.DiaryTags(c.Request.Context())
	h.write(c, data, err)
}

func (h *Handler) Diaries(c *gin.Context) {
	var query request.DiaryQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		invalidQuery(c)
		return
	}
	if err := query.Normalize(); err != nil {
		writeQueryError(c, err)
		return
	}
	data, err := h.service.ListDiaries(c.Request.Context(), query)
	h.write(c, data, err)
}

func (h *Handler) Diary(c *gin.Context) {
	h.detail(c, service.ErrDiaryNotFound, response.CodeDiaryNotFound, "日记不存在或已被删除", func(id string) (interface{}, error) {
		return h.service.GetDiary(c.Request.Context(), id)
	})
}

func (h *Handler) PortfolioCategories(c *gin.Context) {
	h.categories(c, model.CategoryPortfolio)
}

func (h *Handler) PortfolioItems(c *gin.Context) {
	var query request.PortfolioQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		invalidQuery(c)
		return
	}
	query.Normalize()
	data, err := h.service.ListPortfolio(c.Request.Context(), query)
	h.write(c, data, err)
}

func (h *Handler) PortfolioItem(c *gin.Context) {
	h.detail(c, service.ErrPortfolioNotFound, response.CodePortfolioMissing, "作品不存在或已移除", func(id string) (interface{}, error) {
		return h.service.GetPortfolio(c.Request.Context(), id)
	})
}

func (h *Handler) Tools(c *gin.Context) {
	data, err := h.service.ListTools(c.Request.Context())
	h.write(c, data, err)
}

func (h *Handler) Tool(c *gin.Context) {
	h.detail(c, service.ErrToolNotFound, response.CodeToolNotFound, "工具不存在，或该条目不是自研工具", func(id string) (interface{}, error) {
		return h.service.GetTool(c.Request.Context(), id)
	})
}

func (h *Handler) BookmarkCategories(c *gin.Context) {
	h.categories(c, model.CategoryBookmark)
}

func (h *Handler) Bookmarks(c *gin.Context) {
	var query request.BookmarkQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		invalidQuery(c)
		return
	}
	query.Normalize()
	data, err := h.service.ListBookmarks(c.Request.Context(), query)
	h.write(c, data, err)
}

func (h *Handler) Profile(c *gin.Context) {
	data, err := h.service.Profile(c.Request.Context())
	h.write(c, data, err)
}

func (h *Handler) Site(c *gin.Context) {
	data, err := h.service.Site(c.Request.Context())
	h.write(c, data, err)
}

func (h *Handler) Stats(c *gin.Context) {
	data, err := h.service.Stats(c.Request.Context())
	h.write(c, data, err)
}

func (h *Handler) categories(c *gin.Context, scope string) {
	data, err := h.service.Categories(c.Request.Context(), scope)
	h.write(c, data, err)
}

func (h *Handler) detail(c *gin.Context, notFound error, code int, message string, load func(string) (interface{}, error)) {
	id := c.Param("id")
	if !request.ValidPublicID(id) {
		invalidQuery(c)
		return
	}
	data, err := load(id)
	if errors.Is(err, notFound) {
		response.Error(c, http.StatusNotFound, code, message)
		return
	}
	h.write(c, data, err)
}

func (h *Handler) write(c *gin.Context, data interface{}, err error) {
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, "服务内部错误")
		return
	}
	response.OK(c, data)
}

func writeQueryError(c *gin.Context, err error) {
	if errors.Is(err, request.ErrInvalidSort) {
		response.Error(c, http.StatusBadRequest, response.CodeUnsupportedSort, "不支持的排序方式")
		return
	}
	invalidQuery(c)
}

func invalidQuery(c *gin.Context) {
	response.Error(c, http.StatusBadRequest, response.CodeInvalidQuery, "查询参数格式错误")
}
