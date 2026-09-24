package router

import (
	"fastduck/treasure-doc/module/blog_mgr/api"
	"fastduck/treasure-doc/module/blog_mgr/internal/service"

	"github.com/gin-gonic/gin"
)

// RegisterMemo 在前台速记本分组下注册管理端点（调用方已加 Auth 中间件，登录即博主本人）。
func RegisterMemo(group *gin.RouterGroup) {
	api.RegisterMemoRoutes(group, service.New())
}
