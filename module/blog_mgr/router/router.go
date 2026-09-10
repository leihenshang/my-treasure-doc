package router

import (
	"fastduck/treasure-doc/module/blog_mgr/api"
	"fastduck/treasure-doc/module/blog_mgr/internal/service"
	userapi "fastduck/treasure-doc/module/user/api"

	"github.com/gin-gonic/gin"
)

func Register(group *gin.RouterGroup) {
	RegisterService(group, service.New())
}

func RegisterService(group *gin.RouterGroup, manager api.Manager) {
	handler := api.New(manager)

	for _, resource := range api.ResourceNames() {
		route := group.Group("/" + resource)
		route.GET("", handler.List(resource))
		route.POST("", handler.Create(resource))
		route.GET("/:id", handler.Detail(resource))
		route.PATCH("/:id", handler.Update(resource))
		route.PATCH("/:id/fields", handler.UpdateFields(resource))
		route.DELETE("/:id", handler.Delete(resource))
		route.POST("/:id/restore", handler.Restore(resource))
	}

	for _, setting := range []string{"profile", "site"} {
		group.GET("/"+setting, handler.GetSetting(setting))
		group.PUT("/"+setting, handler.PutSetting(setting))
	}

	group.POST("/uploads/images", userapi.UploadBlogImage)
	group.POST("/uploads/medias", userapi.UploadBlogMedias)
}
