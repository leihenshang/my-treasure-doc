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

// webBaseHref 前端产物用相对路径引用资源（./assets/...），固定 base 到站点根，
// 保证 / 与深链（如 /article/1）刷新时资源都从根路径加载。
const webBaseHref = `<base href="/">`

// serveSpaIndex 返回前端单页应用入口，用于 history 模式路由兜底。
func serveSpaIndex(c *gin.Context) {
	index := filepath.Join(config.WebPath, "index.html")
	data, err := os.ReadFile(index)
	if err != nil {
		c.String(http.StatusNotFound, "frontend not built")
		return
	}
	html := strings.Replace(string(data), "<head>", "<head>"+webBaseHref, 1)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// serveWebFile 尝试返回 web 目录（前端构建产物）下命中的静态文件，未命中返回 false。
func serveWebFile(c *gin.Context, urlPath string) bool {
	webRoot := filepath.Clean(config.WebPath)
	file := filepath.Join(webRoot, filepath.Clean(urlPath))
	// 阻断 ../ 越权访问 web 目录之外的文件
	if !strings.HasPrefix(file, webRoot+string(os.PathSeparator)) {
		return false
	}
	info, err := os.Stat(file)
	if err != nil || info.IsDir() {
		return false
	}
	c.File(file)
	return true
}

func InitRouter(r *gin.Engine) {
	r.Static("/files", config.FilesPath)

	r.GET("/", serveSpaIndex)

	// 兜底路由：API 未匹配返回 JSON 404；命中前端静态文件则直接返回，
	// 其余未知路径（前端路由深链）交给 index.html 处理。
	r.NoRoute(func(c *gin.Context) {
		urlPath := c.Request.URL.Path
		if strings.HasPrefix(urlPath, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "not found"})
			return
		}
		// 兼容旧入口 /web/xxx，收敛到根路径
		if rest := strings.TrimPrefix(urlPath, "/web"); rest != urlPath && (rest == "" || rest[0] == '/') {
			c.Redirect(http.StatusMovedPermanently, "/"+strings.TrimPrefix(rest, "/"))
			return
		}
		if serveWebFile(c, urlPath) {
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
		userRoute.POST("/change-pwd", userApi.UserChangePwd)
	}
}
