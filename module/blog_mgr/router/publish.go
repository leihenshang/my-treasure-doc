package router

import (
	blogmgrapi "fastduck/treasure-doc/module/blog_mgr/api"
	"fastduck/treasure-doc/module/blog_mgr/internal/service"

	"github.com/gin-gonic/gin"
)

// RegisterPublish 挂载机器令牌发布接口（/api/publish/?:resource）。
// https://xxx/api/publish/... 需要 X-Publish-Token，发布即强制 publishStatus=published。
// 鉴权中间件由调用方在该 group 上叠加（此处只负责路由与真实 service 的装配）。
func RegisterPublish(group *gin.RouterGroup) {
	blogmgrapi.RegisterPublishRoutes(group, service.New())
}
