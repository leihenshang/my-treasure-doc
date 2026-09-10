package router

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"fastduck/treasure-doc/module/user/api"
	"fastduck/treasure-doc/module/user/config"
	"fastduck/treasure-doc/module/user/router/middleware"

	blogrouter "fastduck/treasure-doc/module/blog/router"
	blogmgrrouter "fastduck/treasure-doc/module/blog_mgr/router"

	"github.com/gin-gonic/gin"
)

// webBaseHref 把前端静态资源基址固定到站点根，保证 / 与深链（如 /article/1）
// 刷新时相对资源都从根路径加载。
const webBaseHref = `<base href="/">`

func InitRouter(r *gin.Engine) {
	registerFrontend(r)
	registerAPI(r)
}

// registerFrontend 托管前端构建产物：/files 静态资源、SPA 入口，以及 history 模式兜底。
func registerFrontend(r *gin.Engine) {
	r.Static("/files", config.FilesPath)
	r.GET("/", serveSpaIndex)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "pong!"})
	})

	// API 未匹配返回 JSON 404；命中前端静态文件则直接返回，其余交给 index.html。
	r.NoRoute(func(c *gin.Context) {
		urlPath := c.Request.URL.Path
		if strings.HasPrefix(urlPath, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "msg": "not found"})
			return
		}
		// 兼容旧入口 /web/xxx，收敛到根路径
		if rest := strings.TrimPrefix(urlPath, "/web"); rest != urlPath && (rest == "" || rest[0] == '/') {
			c.Redirect(http.StatusMovedPermanently, "/"+strings.TrimPrefix(rest, "/"))
			return
		}
		if !serveWebFile(c, urlPath) {
			serveSpaIndex(c)
		}
	})
}

// registerAPI 挂载公开博客、后台管理与用户接口。
func registerAPI(r *gin.Engine) {
	apiBase := r.Group("api")

	blogrouter.Register(apiBase)

	blogMgr := apiBase.Group("blog-mgr")
	blogMgr.Use(middleware.Cors(), middleware.Auth(), middleware.RequireAdmin())
	blogmgrrouter.Register(blogMgr)

	userAPI := api.NewUserApi()
	user := apiBase.Group("user").Use(middleware.Cors())
	user.GET("/captcha", userAPI.UserCaptcha)
	user.POST("/login", userAPI.UserLogin)
	user.Use(middleware.Auth())
	user.POST("/logout", userAPI.UserLogout)
	user.POST("/change-pwd", userAPI.UserChangePwd)
}

// serveSpaIndex 返回前端单页应用入口。
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

// serveWebFile 返回 web 目录下命中的静态文件，未命中或越权返回 false。
func serveWebFile(c *gin.Context, urlPath string) bool {
	webRoot := filepath.Clean(config.WebPath)
	file := filepath.Join(webRoot, filepath.Clean(urlPath))
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
