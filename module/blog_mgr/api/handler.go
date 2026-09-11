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

// 管理端业务码，与前端约定保持一致。
const (
	codeInvalidRequest     = 40001
	codeReferenceAbsent    = 40002
	codeToolURLRequired    = 40003
	codeToolStatusRequired = 40004
	codeNotFound           = 40410
	codeConflict           = 40900
	codeInternal           = 50000
)

// Manager 是管理端业务入口，由 internal/service 实现，测试可注入替身。
type Manager interface {
	List(context.Context, string, request.List) (service.Page, error)
	Get(context.Context, string, string) (interface{}, error)
	Create(context.Context, string, interface{}) (interface{}, error)
	Update(context.Context, string, string, interface{}) (interface{}, error)
	UpdateFields(context.Context, string, string, map[string]interface{}) (interface{}, error)
	Delete(context.Context, string, string) error
	DeleteMany(context.Context, string, []string) (int64, error)
	Restore(context.Context, string, string) error
	GetSetting(context.Context, string) (interface{}, error)
	PutSetting(context.Context, string, interface{}) (interface{}, error)
	Stats(context.Context) (response.Stats, error)
}

// 资源标识与请求体类型集中定义，供路由注册与请求解析共用。
var (
	resourceNames = []string{"categories", "tags", "posts", "diaries", "portfolio-items", "tools", "bookmarks"}

	resourceBinders = map[string]binder{
		"categories":      jsonBinder[request.Category]{},
		"tags":            jsonBinder[request.Tag]{},
		"posts":           jsonBinder[request.Post]{},
		"diaries":         jsonBinder[request.Diary]{},
		"portfolio-items": jsonBinder[request.Portfolio]{},
		"tools":           jsonBinder[request.Tool]{},
		"bookmarks":       jsonBinder[request.Bookmark]{},
	}
)

// ResourceNames 返回管理端支持的资源标识，顺序固定。
func ResourceNames() []string { return resourceNames }

// binder 把请求体解析为某个资源的入参值。
type binder interface {
	bind(*gin.Context) (interface{}, bool)
}

type jsonBinder[T any] struct{}

func (jsonBinder[T]) bind(c *gin.Context) (interface{}, bool) {
	var payload T
	if c.ShouldBindJSON(&payload) != nil {
		return nil, false
	}
	return payload, true
}

type Handler struct{ service Manager }

func New(manager Manager) *Handler { return &Handler{service: manager} }

func (h *Handler) List(resource string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var query request.List
		if c.ShouldBindQuery(&query) != nil || query.Normalize() != nil {
			badRequest(c)
			return
		}
		data, err := h.service.List(c.Request.Context(), resource, query)
		h.write(c, data, err, false)
	}
}

func (h *Handler) Detail(resource string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := pathID(c)
		if !ok {
			return
		}
		data, err := h.service.Get(c.Request.Context(), resource, id)
		h.write(c, data, err, false)
	}
}

func (h *Handler) Create(resource string) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload, ok := bindResource(c, resource)
		if !ok {
			return
		}
		data, err := h.service.Create(c.Request.Context(), resource, payload)
		h.write(c, data, err, true)
	}
}

func (h *Handler) Update(resource string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := pathID(c)
		if !ok {
			return
		}
		payload, ok := bindResource(c, resource)
		if !ok {
			return
		}
		data, err := h.service.Update(c.Request.Context(), resource, id, payload)
		h.write(c, data, err, false)
	}
}

// UpdateFields 列表快捷设置，仅提交需要修改的字段，例如 {"pinned": true}。
func (h *Handler) UpdateFields(resource string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := pathID(c)
		if !ok {
			return
		}
		var payload map[string]interface{}
		if c.ShouldBindJSON(&payload) != nil || len(payload) == 0 {
			badRequest(c)
			return
		}
		data, err := h.service.UpdateFields(c.Request.Context(), resource, id, payload)
		h.write(c, data, err, false)
	}
}

func (h *Handler) Delete(resource string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := pathID(c)
		if !ok {
			return
		}
		err := h.service.Delete(c.Request.Context(), resource, id)
		h.write(c, map[string]bool{"deleted": true}, err, false)
	}
}

// DeleteMany 批量软删除，请求体 {"ids": ["..."]}；返回 {"deleted": 实际删除条数}。
func (h *Handler) DeleteMany(resource string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload request.BatchIDs
		if c.ShouldBindJSON(&payload) != nil || payload.Normalize() != nil {
			badRequest(c)
			return
		}
		deleted, err := h.service.DeleteMany(c.Request.Context(), resource, payload.IDs)
		h.write(c, gin.H{"deleted": deleted}, err, false)
	}
}

func (h *Handler) Restore(resource string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := pathID(c)
		if !ok {
			return
		}
		err := h.service.Restore(c.Request.Context(), resource, id)
		h.write(c, map[string]bool{"restored": true}, err, false)
	}
}

func (h *Handler) GetSetting(setting string) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := h.service.GetSetting(c.Request.Context(), setting)
		h.write(c, data, err, false)
	}
}

func (h *Handler) PutSetting(setting string) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload, ok := bindSetting(c, setting)
		if !ok {
			return
		}
		data, err := h.service.PutSetting(c.Request.Context(), setting, payload)
		h.write(c, data, err, false)
	}
}

// Stats 返回后台仪表盘的总览数据。
func (h *Handler) Stats() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := h.service.Stats(c.Request.Context())
		h.write(c, data, err, false)
	}
}

// pathID 读取并校验路径上的资源 ID。
func pathID(c *gin.Context) (string, bool) {
	id := c.Param("id")
	if !request.ValidID(id) {
		badRequest(c)
		return "", false
	}
	return id, true
}

func bindResource(c *gin.Context, resource string) (interface{}, bool) {
	b, ok := resourceBinders[resource]
	if !ok {
		badRequest(c)
		return nil, false
	}
	payload, ok := b.bind(c)
	if !ok {
		badRequest(c)
		return nil, false
	}
	return payload, true
}

func bindSetting(c *gin.Context, setting string) (interface{}, bool) {
	if setting == "profile" {
		var payload blogresponse.Profile
		if c.ShouldBindJSON(&payload) != nil {
			return nil, false
		}
		return payload, true
	}
	var payload blogresponse.Site
	if c.ShouldBindJSON(&payload) != nil {
		return nil, false
	}
	return payload, true
}

func (h *Handler) write(c *gin.Context, data interface{}, err error, created bool) {
	switch {
	case errors.Is(err, service.ErrReferenceNotFound):
		response.Error(c, http.StatusBadRequest, codeReferenceAbsent, "关联的分类或标签不存在，请先创建")
	case errors.Is(err, service.ErrToolURLRequired):
		response.Error(c, http.StatusBadRequest, codeToolURLRequired, "外链工具必须填写 HTTPS 地址")
	case errors.Is(err, service.ErrToolStatusRequired):
		response.Error(c, http.StatusBadRequest, codeToolStatusRequired, "自研工具必须填写开发状态")
	case errors.Is(err, service.ErrInvalid):
		badRequest(c)
	case errors.Is(err, service.ErrNotFound):
		response.Error(c, http.StatusNotFound, codeNotFound, "资源不存在")
	case errors.Is(err, service.ErrConflict):
		response.Error(c, http.StatusConflict, codeConflict, "数据已变更或标识重复")
	case err != nil:
		response.Error(c, http.StatusInternalServerError, codeInternal, "服务内部错误")
	case created:
		response.Created(c, data)
	default:
		response.OK(c, data)
	}
}

func badRequest(c *gin.Context) {
	response.Error(c, http.StatusBadRequest, codeInvalidRequest, "请求参数格式错误")
}
