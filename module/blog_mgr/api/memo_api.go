package api

import (
	"errors"
	"net/http"

	"fastduck/treasure-doc/module/blog_mgr/data/request"
	"fastduck/treasure-doc/module/blog_mgr/data/response"
	"fastduck/treasure-doc/module/blog_mgr/internal/service"

	"github.com/gin-gonic/gin"
)

// RegisterMemoRoutes 注册前台速记本管理端点（/api/memo/*，仅 Auth()）。
func RegisterMemoRoutes(group *gin.RouterGroup, manager *service.Service) {
	h := &MemoHandler{service: manager}
	group.GET("/memos", h.List)
	group.POST("/memos", h.Create)
	group.POST("/memos/batch-delete", h.DeleteMany)
	group.GET("/memos/:id", h.Detail)
	group.PUT("/memos/:id", h.Update)
	group.PATCH("/memos/:id/fields", h.UpdateFields)
	group.DELETE("/memos/:id", h.Delete)
}

type MemoHandler struct{ service *service.Service }

func (h *MemoHandler) List(c *gin.Context) {
	var query request.List
	if c.ShouldBindQuery(&query) != nil || query.Normalize() != nil {
		badRequest(c)
		return
	}
	data, err := h.service.ListMemo(c.Request.Context(), query)
	h.write(c, data, err, false)
}

func (h *MemoHandler) Detail(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	data, err := h.service.GetMemo(c.Request.Context(), id)
	h.write(c, data, err, false)
}

func (h *MemoHandler) Create(c *gin.Context) {
	var payload request.Memo
	if c.ShouldBindJSON(&payload) != nil {
		badRequest(c)
		return
	}
	data, err := h.service.CreateMemo(c.Request.Context(), payload)
	h.write(c, data, err, true)
}

func (h *MemoHandler) Update(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	var payload request.Memo
	if c.ShouldBindJSON(&payload) != nil {
		badRequest(c)
		return
	}
	data, err := h.service.UpdateMemo(c.Request.Context(), id, payload)
	h.write(c, data, err, false)
}

func (h *MemoHandler) UpdateFields(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	var payload map[string]interface{}
	if c.ShouldBindJSON(&payload) != nil || len(payload) == 0 {
		badRequest(c)
		return
	}
	data, err := h.service.UpdateMemoFields(c.Request.Context(), id, payload)
	h.write(c, data, err, false)
}

func (h *MemoHandler) Delete(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	err := h.service.DeleteMemo(c.Request.Context(), id)
	h.write(c, map[string]bool{"deleted": true}, err, false)
}

// DeleteMany 批量软删除速记，请求体 {"ids":["..."]}，返回 {"deleted": 实际删除条数}。
func (h *MemoHandler) DeleteMany(c *gin.Context) {
	var payload request.BatchIDs
	if c.ShouldBindJSON(&payload) != nil || payload.Normalize() != nil {
		badRequest(c)
		return
	}
	deleted, err := h.service.DeleteMemoBatch(c.Request.Context(), payload.IDs)
	h.write(c, gin.H{"deleted": deleted}, err, false)
}

func (h *MemoHandler) write(c *gin.Context, data interface{}, err error, created bool) {
	if err != nil {
		var fieldErr *request.FieldError
		switch {
		case errors.As(err, &fieldErr):
			response.ErrorWithData(c, http.StatusBadRequest, codeInvalidRequest, fieldErr.Reason, gin.H{"field": fieldErr.Field})
		case errors.Is(err, service.ErrInvalid):
			badRequest(c)
		case errors.Is(err, service.ErrNotFound):
			response.Error(c, http.StatusNotFound, codeNotFound, "速记不存在")
		case errors.Is(err, service.ErrConflict):
			response.Error(c, http.StatusConflict, codeConflict, "速记已变更，请刷新后重试")
		default:
			response.Error(c, http.StatusInternalServerError, codeInternal, "服务内部错误")
		}
		return
	}
	if created {
		response.Created(c, data)
		return
	}
	response.OK(c, data)
}
