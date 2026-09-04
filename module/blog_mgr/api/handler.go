package api

import (
	"context"
	"errors"
	"net/http"

	blogresponse "fastduck/treasure-doc/module/blog/data/response"
	"fastduck/treasure-doc/module/blog_mgr/data/request"
	"fastduck/treasure-doc/module/blog_mgr/data/response"
	"fastduck/treasure-doc/module/blog_mgr/internal/service"

	"github.com/gin-gonic/gin"
)

type Manager interface {
	List(context.Context, string, request.List) (service.Page, error)
	Get(context.Context, string, string) (interface{}, error)
	Create(context.Context, string, interface{}) (interface{}, error)
	Update(context.Context, string, string, interface{}) (interface{}, error)
	Delete(context.Context, string, string) error
	Restore(context.Context, string, string) error
	GetSetting(context.Context, string) (interface{}, error)
	PutSetting(context.Context, string, interface{}) (interface{}, error)
}

type Handler struct{ service Manager }

func New(manager Manager) *Handler { return &Handler{service: manager} }

func (h *Handler) List(c *gin.Context) {
	var query request.List
	if c.ShouldBindQuery(&query) != nil || query.Normalize() != nil {
		badRequest(c)
		return
	}
	data, err := h.service.List(c.Request.Context(), c.Param("resource"), query)
	h.write(c, data, err, false)
}
func (h *Handler) Detail(c *gin.Context) {
	if !request.ValidID(c.Param("id")) {
		badRequest(c)
		return
	}
	data, err := h.service.Get(c.Request.Context(), c.Param("resource"), c.Param("id"))
	h.write(c, data, err, false)
}
func (h *Handler) Create(c *gin.Context) {
	payload, ok := bindResource(c, c.Param("resource"))
	if !ok {
		return
	}
	data, err := h.service.Create(c.Request.Context(), c.Param("resource"), payload)
	h.write(c, data, err, true)
}
func (h *Handler) Update(c *gin.Context) {
	if !request.ValidID(c.Param("id")) {
		badRequest(c)
		return
	}
	payload, ok := bindResource(c, c.Param("resource"))
	if !ok {
		return
	}
	data, err := h.service.Update(c.Request.Context(), c.Param("resource"), c.Param("id"), payload)
	h.write(c, data, err, false)
}
func (h *Handler) Delete(c *gin.Context) {
	if !request.ValidID(c.Param("id")) {
		badRequest(c)
		return
	}
	err := h.service.Delete(c.Request.Context(), c.Param("resource"), c.Param("id"))
	h.write(c, map[string]bool{"deleted": true}, err, false)
}
func (h *Handler) Restore(c *gin.Context) {
	if !request.ValidID(c.Param("id")) {
		badRequest(c)
		return
	}
	err := h.service.Restore(c.Request.Context(), c.Param("resource"), c.Param("id"))
	h.write(c, map[string]bool{"restored": true}, err, false)
}
func (h *Handler) GetSetting(c *gin.Context) {
	data, err := h.service.GetSetting(c.Request.Context(), c.Param("setting"))
	h.write(c, data, err, false)
}
func (h *Handler) PutSetting(c *gin.Context) {
	var payload interface{}
	if c.Param("setting") == "profile" {
		payload = &blogresponse.Profile{}
	} else {
		payload = &blogresponse.Site{}
	}
	if c.ShouldBindJSON(payload) != nil {
		badRequest(c)
		return
	}
	switch value := payload.(type) {
	case *blogresponse.Profile:
		payload = *value
	case *blogresponse.Site:
		payload = *value
	}
	data, err := h.service.PutSetting(c.Request.Context(), c.Param("setting"), payload)
	h.write(c, data, err, false)
}

func bindResource(c *gin.Context, resource string) (interface{}, bool) {
	var payload interface{}
	switch resource {
	case "categories":
		payload = &request.Category{}
	case "tags":
		payload = &request.Tag{}
	case "posts":
		payload = &request.Post{}
	case "diaries":
		payload = &request.Diary{}
	case "portfolio-items":
		payload = &request.Portfolio{}
	case "tools":
		payload = &request.Tool{}
	case "bookmarks":
		payload = &request.Bookmark{}
	default:
		badRequest(c)
		return nil, false
	}
	if c.ShouldBindJSON(payload) != nil {
		badRequest(c)
		return nil, false
	}
	switch value := payload.(type) {
	case *request.Category:
		return *value, true
	case *request.Tag:
		return *value, true
	case *request.Post:
		return *value, true
	case *request.Diary:
		return *value, true
	case *request.Portfolio:
		return *value, true
	case *request.Tool:
		return *value, true
	case *request.Bookmark:
		return *value, true
	}
	return nil, false
}
func (h *Handler) write(c *gin.Context, data interface{}, err error, created bool) {
	if errors.Is(err, service.ErrInvalid) {
		badRequest(c)
		return
	}
	if errors.Is(err, service.ErrNotFound) {
		response.Error(c, http.StatusNotFound, 40410, "资源不存在")
		return
	}
	if errors.Is(err, service.ErrConflict) {
		response.Error(c, http.StatusConflict, 40900, "数据已变更或标识重复")
		return
	}
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50000, "服务内部错误")
		return
	}
	if created {
		response.Created(c, data)
	} else {
		response.OK(c, data)
	}
}
func badRequest(c *gin.Context) {
	response.Error(c, http.StatusBadRequest, 40001, "请求参数格式错误")
}
