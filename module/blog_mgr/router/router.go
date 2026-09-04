package router

import (
	"fastduck/treasure-doc/module/blog_mgr/api"
	"fastduck/treasure-doc/module/blog_mgr/internal/service"

	"github.com/gin-gonic/gin"
)

func Register(group *gin.RouterGroup) {
	RegisterService(group, service.New())
}
func RegisterService(group *gin.RouterGroup, manager api.Manager) {
	handler := api.New(manager)
	resources := []string{"categories", "tags", "posts", "diaries", "portfolio-items", "tools", "bookmarks"}
	for _, resource := range resources {
		route := group.Group("/" + resource)
		route.GET("", setParam("resource", resource), handler.List)
		route.POST("", setParam("resource", resource), handler.Create)
		route.GET("/:id", setParam("resource", resource), handler.Detail)
		route.PATCH("/:id", setParam("resource", resource), handler.Update)
		route.DELETE("/:id", setParam("resource", resource), handler.Delete)
		route.POST("/:id/restore", setParam("resource", resource), handler.Restore)
	}
	for _, setting := range []string{"profile", "site"} {
		route := group.Group("/" + setting)
		route.GET("", setParam("setting", setting), handler.GetSetting)
		route.PUT("", setParam("setting", setting), handler.PutSetting)
	}
}
func setParam(key, value string) gin.HandlerFunc {
	return func(c *gin.Context) { c.Params = append(c.Params, gin.Param{Key: key, Value: value}); c.Next() }
}
