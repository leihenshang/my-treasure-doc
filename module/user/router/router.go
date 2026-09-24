package router

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"fastduck/treasure-doc/module/user/api"
	"fastduck/treasure-doc/module/user/config"
	"fastduck/treasure-doc/module/user/global"
	"fastduck/treasure-doc/module/user/router/middleware"

	blogrouter "fastduck/treasure-doc/module/blog/router"
	blogmgrrouter "fastduck/treasure-doc/module/blog_mgr/router"

	"github.com/gin-gonic/gin"
)

// webBaseHref 把前端静态资源基址固定到站点根，保证 / 与深链（如 /article/1）
// 刷新时相对资源都从根路径加载。
const webBaseHref = `<base href="/">`

// staticMountPrefixes 是静态资源挂载点前缀，这类路径命中不到文件时永远不是前端路由：
//   - /assets/ 由 Vite 产出，所有带哈希的 js/css/字体都在下面；
//   - /files/  由 r.Static 托管上传文件。
//
// Vue Router 的路径只会落在 /Blog 或 /BlogManage 前缀内，因此这些前缀可以安全地
// 排除在 SPA 兜底之外。
var staticMountPrefixes = []string{"/assets/", "/files/"}

// isStaticMountPath 判断请求是否指向静态资源挂载点。
func isStaticMountPath(urlPath string) bool {
	for _, prefix := range staticMountPrefixes {
		if strings.HasPrefix(urlPath, prefix) {
			return true
		}
	}
	return false
}

func InitRouter(r *gin.Engine) {
	// 全局中间件：静态缓存、安全响应头、gzip 压缩，以及针对登录/上传的限流
	r.Use(
		middleware.StaticCache(),
		middleware.SecureHeaders(),
		middleware.Gzip(),
		middleware.RateLimit([]middleware.RateRule{
			{Prefix: "/api/user/login", Rate: 0.1, Burst: 5},
			{Prefix: "/api/blog-mgr/uploads/", Rate: 0.5, Burst: 20},
			{Prefix: "/api/publish", Rate: 0.2, Burst: 20},
		}),
	)

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
			// 静态资源缺失时不能回退 index.html：浏览器对模块脚本 / 样式表做严格
			// MIME 校验，拿到 text/html 会直接报 "Expected a JavaScript-or-Wasm
			// module script but the server responded with a MIME type of text/html"，
			// 把「资源没上传 / index.html 与 assets 版本不一致」这类部署问题伪装成
			// 难懂的 MIME 错误。这里显式 404，并禁止缓存这个否定结果，避免修好部署
			// 后仍被浏览器缓存的 404 挡住。
			//
			// 判定范围：已知静态挂载前缀（/assets/、/files/）之外，凡是带扩展名的
			// 请求（js/css/字体/图片等）都视为静态资源。这样即使前端 index.html 把
			// 资源解析到了非 /assets/ 前缀（例如 /web/assets/...），缺失时也返回 404
			// 而非 HTML，错误表现更直观。
			if isStaticMountPath(urlPath) || filepath.Ext(urlPath) != "" {
				c.Header("Cache-Control", "no-store")
				c.String(http.StatusNotFound, "resource not found")
				return
			}
			serveSpaIndex(c)
		}
	})
}

// registerAPI 挂载公开博客、后台管理与用户接口。
func registerAPI(r *gin.Engine) {
	apiBase := r.Group("api")

	blogrouter.Register(apiBase)
	blogrouter.RegisterSiteFiles(r)

	blogMgr := apiBase.Group("blog-mgr")
	blogMgr.Use(middleware.Auth(), middleware.RequireAdmin())
	blogmgrrouter.Register(blogMgr)

	var backupAllowIPs, publishAllowIPs []string
	if cfg := global.GetConf(); cfg != nil {
		backupAllowIPs = cfg.Backup.AllowIPs
		publishAllowIPs = cfg.Publish.AllowIPs
	}

	// NAS 机器令牌下载完整备份（免登录），仅只读导出，不走管理端鉴权链。
	apiBase.GET("/backup/export", middleware.IPWhitelist(backupAllowIPs), middleware.BackupToken(), api.NewBackupApi().NasExportArchive)

	// 机器令牌发布接口：免登录、用 X-Publish-Token，创建即发布；可叠加 IP 白名单。
	publish := apiBase.Group("publish")
	publish.Use(middleware.IPWhitelist(publishAllowIPs), middleware.RequirePublishToken())
	blogmgrrouter.RegisterPublish(publish)
	// 发布时上传文章中的图片/附件：同样用发布令牌鉴权，返回可替换引用的 /files/blog/... 地址。
	publish.POST("/uploads", api.PublishUploadMedias)

	userAPI := api.NewUserApi()
	user := apiBase.Group("user")
	user.GET("/captcha", userAPI.UserCaptcha)
	user.POST("/login", userAPI.UserLogin)
	user.Use(middleware.Auth())
	user.POST("/logout", userAPI.UserLogout)
	user.POST("/change-pwd", userAPI.UserChangePwd)

	// 前台速记本管理：登录即博主本人，仅 Auth()（不叠加 RequireAdmin）。
	memo := apiBase.Group("memo")
	memo.Use(middleware.Auth())
	blogmgrrouter.RegisterMemo(memo)
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
	// 入口页不缓存，保证前端发版后用户能立即拿到新资源引用
	c.Header("Cache-Control", "no-cache")
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
