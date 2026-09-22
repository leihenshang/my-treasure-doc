package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

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
	VisitorStats(context.Context, int) (response.VisitorStats, error)
	// 发布端（令牌）只读/覆盖
	ListCategoriesForPublish(context.Context, string) ([]service.PublishCategory, error)
	ListTagsForPublish(context.Context) ([]service.PublishTag, error)
	PublishLookup(context.Context, string, string) (service.PublishStatus, error)
	PublishUpdate(context.Context, string, string, interface{}) (interface{}, error)
	// 编辑历史
	ListEditHistory(context.Context, string, string) ([]service.EditHistoryMeta, error)
	GetEditHistory(context.Context, string, string, int) (service.EditHistoryDetail, error)
	RestoreEditHistory(context.Context, string, string, int) (interface{}, error)
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

// RegisterPublishRoutes 注册机器令牌发布接口（/api/publish/…）：
// 免管理员 X-Token、改用 X-Publish-Token 鉴权，创建即发布（强制 publishStatus=published）。
// 覆盖：发布创建、分类/标签列表、发布状态查询、按 slug 强制覆盖（仅 posts/diaries）。
func RegisterPublishRoutes(group *gin.RouterGroup, manager Manager) {
	handler := New(manager)
	for _, resource := range []string{"posts", "diaries", "portfolio-items", "tools", "bookmarks"} {
		group.POST("/"+resource, handler.PublishCreate(resource))
	}
	// 发布端元数据：插件据此渲染分类/标签下拉
	group.GET("/categories", handler.PublishCategories())
	group.GET("/tags", handler.PublishTags())
	// 发布状态查询 + 强制覆盖（仅文章/日记）
	for _, resource := range []string{"posts", "diaries"} {
		group.GET("/"+resource+"/by-slug/:slug", handler.PublishLookup(resource))
		group.PUT("/"+resource+"/:slug", handler.PublishForceUpdate(resource))
	}
}

// PublishCategories 发布端可选分类列表（?scope= 可选）。
func (h *Handler) PublishCategories() gin.HandlerFunc {
	return func(c *gin.Context) {
		scope := strings.TrimSpace(c.Query("scope"))
		if scope != "" && !request.ValidScope(scope) {
			badRequest(c)
			return
		}
		data, err := h.service.ListCategoriesForPublish(c.Request.Context(), scope)
		h.write(c, data, err, false)
	}
}

// PublishTags 发布端可选标签列表。
func (h *Handler) PublishTags() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := h.service.ListTagsForPublish(c.Request.Context())
		h.write(c, data, err, false)
	}
}

// PublishLookup 查询某文章/日记是否已发布（按 slug）。
func (h *Handler) PublishLookup(resource string) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := h.service.PublishLookup(c.Request.Context(), resource, c.Param("slug"))
		h.write(c, data, err, false)
	}
}

// PublishForceUpdate 按 slug 强制覆盖文章/日记并发布。
func (h *Handler) PublishForceUpdate(resource string) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload, ok := bindResource(c, resource)
		if !ok {
			return
		}
		data, err := h.service.PublishUpdate(c.Request.Context(), resource, c.Param("slug"), forcePublishStatus(payload))
		h.write(c, data, err, false)
	}
}

// HistoryList 某文章/日记的编辑历史列表。
func (h *Handler) HistoryList() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := h.service.ListEditHistory(c.Request.Context(), c.Param("resource"), c.Param("id"))
		h.write(c, data, err, false)
	}
}

// HistoryDetail 某历史版本的完整快照。
func (h *Handler) HistoryDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		seq, err := strconv.Atoi(c.Param("seq"))
		if err != nil {
			badRequest(c)
			return
		}
		data, err := h.service.GetEditHistory(c.Request.Context(), c.Param("resource"), c.Param("id"), seq)
		h.write(c, data, err, false)
	}
}

// HistoryRestore 恢复某历史版本为当前内容。
func (h *Handler) HistoryRestore() gin.HandlerFunc {
	return func(c *gin.Context) {
		seq, err := strconv.Atoi(c.Param("seq"))
		if err != nil {
			badRequest(c)
			return
		}
		data, err := h.service.RestoreEditHistory(c.Request.Context(), c.Param("resource"), c.Param("id"), seq)
		h.write(c, data, err, false)
	}
}

// PublishCreate 发布接口的创建处理：绑定后强制发布，再走与后台一致的创建链路。
func (h *Handler) PublishCreate(resource string) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload, ok := bindResource(c, resource)
		if !ok {
			return
		}
		data, err := h.service.Create(c.Request.Context(), resource, forcePublishStatus(payload))
		h.write(c, data, err, true)
	}
}

// forcePublishStatus 把待发布内容的 publishStatus 强制为 published（发布接口语义），返回改写后的入参。
func forcePublishStatus(payload interface{}) interface{} {
	switch value := payload.(type) {
	case request.Post:
		value.PublishStatus = "published"
		return value
	case request.Diary:
		value.PublishStatus = "published"
		return value
	case request.Portfolio:
		value.PublishStatus = "published"
		return value
	case request.Tool:
		value.PublishStatus = "published"
		return value
	case request.Bookmark:
		value.PublishStatus = "published"
		return value
	}
	return payload
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

// VisitorStats 返回时间范围内的访客 IP 统计，days 默认 7，限制在 [1, 365]。
func (h *Handler) VisitorStats() gin.HandlerFunc {
	return func(c *gin.Context) {
		days := 7
		if raw := c.Query("days"); raw != "" {
			n, err := strconv.Atoi(raw)
			if err != nil || n < 1 || n > 365 {
				badRequest(c)
				return
			}
			days = n
		}
		data, err := h.service.VisitorStats(c.Request.Context(), days)
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
	// 字段级校验失败：把「哪个字段 + 该怎么改」一并带出去，前端可直接定位到表单项
	var fieldErr *request.FieldError
	if errors.As(err, &fieldErr) {
		response.ErrorWithData(c, http.StatusBadRequest, codeInvalidRequest, fieldErr.Reason, gin.H{"field": fieldErr.Field})
		return
	}
	switch {
	case errors.Is(err, service.ErrReferenceNotFound):
		response.ErrorWithData(c, http.StatusBadRequest, codeReferenceAbsent, "关联的分类或标签不存在，请先创建", gin.H{"field": "categoryId 或 tagIds"})
	case errors.Is(err, service.ErrToolURLRequired):
		response.ErrorWithData(c, http.StatusBadRequest, codeToolURLRequired, "外链工具必须填写有效地址（支持 http/https 等协议，不支持 javascript: 这类地址）", gin.H{"field": "url"})
	case errors.Is(err, service.ErrToolStatusRequired):
		response.ErrorWithData(c, http.StatusBadRequest, codeToolStatusRequired, "自研工具必须填写开发状态", gin.H{"field": "developmentStatus"})
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
