package router

import (
	"net/http"

	blogrouter "fastduck/treasure-doc/module/blog/router"
	blogmgrrouter "fastduck/treasure-doc/module/blog_mgr/router"
	"fastduck/treasure-doc/module/user/config"
	"fastduck/treasure-doc/module/user/router/middleware"

	"github.com/gin-gonic/gin"

	"fastduck/treasure-doc/module/user/api"
)

func InitRouter(r *gin.Engine) {
	r.Static("/web", config.WebPath)
	r.Static("/files", config.FilesPath)

	r.Any("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/web")
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "pong!",
		})
	})

	apiBase := r.Group("api")
	blogrouter.Register(apiBase)
	blogMgrRoute := apiBase.Group("blog-mgr")
	blogMgrRoute.Use(middleware.Cors(), middleware.Auth(), middleware.RequireAdmin())
	blogmgrrouter.Register(blogMgrRoute)

	// 博客管理员登录、退出登录与权限校验相关接口
	{
		userApi := api.NewUserApi()
		userRoute := apiBase.Group("user").Use(middleware.Cors())
		userRoute.POST("/captcha", userApi.UserCaptcha)
		userRoute.POST("/login", userApi.UserLogin)
		userRoute.Use(middleware.Auth())
		userRoute.POST("/logout", userApi.UserLogout)
	}
}
