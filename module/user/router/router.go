package router

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	blogrouter "fastduck/treasure-doc/module/blog/router"
	blogmgrrouter "fastduck/treasure-doc/module/blog_mgr/router"
	"fastduck/treasure-doc/module/user/config"
	"fastduck/treasure-doc/module/user/router/middleware"

	"github.com/gin-gonic/gin"

	"fastduck/treasure-doc/module/user/api"
)

// serveSpaIndex 返回前端单页应用入口，用于 history 模式路由兜底。
func serveSpaIndex(c *gin.Context) {
	index := filepath.Join(config.WebPath, "index.html")
	if _, err := os.Stat(index); err != nil {
		c.String(http.StatusNotFound, "frontend not built")
		return
	}
	c.File(index)
}

func InitRouter(r *gin.Engine) {
	r.Static("/files", config.FilesPath)

	// 前端静态资源：命中真实文件则返回；未命中（前端路由深链）回退 index.html。
	r.GET("/web", serveSpaIndex)
	r.GET("/web/*filepath", func(c *gin.Context) {
		fp := c.Param("filepath")
		if fp == "" || fp == "/" {
			serveSpaIndex(c)
			return
		}
		file := filepath.Join(config.WebPath, fp)
		if info, err := os.Stat(file); err != nil || info.IsDir() {
			// 前端 history 模式深链（如 /web/article/1）刷新：回退到入口页。
			serveSpaIndex(c)
			return
		}
		c.File(file)
	})

	r.Any("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/web")
	})

	// 兜底路由：API 未匹配返回 JSON 404；其余未知路径交给前端路由处理。
	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "not found"})
			return
		}
		serveSpaIndex(c)
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
		userRoute.GET("/captcha", userApi.UserCaptcha)
		userRoute.POST("/login", userApi.UserLogin)
		userRoute.Use(middleware.Auth())
		userRoute.POST("/logout", userApi.UserLogout)
	}
}
