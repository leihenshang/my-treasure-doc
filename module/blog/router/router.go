package router

import (
	"fastduck/treasure-doc/module/blog/api"
	"fastduck/treasure-doc/module/blog/internal/service"

	"github.com/gin-gonic/gin"
)

func Register(apiGroup *gin.RouterGroup) {
	RegisterService(apiGroup, service.New())
}

func RegisterSiteFiles(r *gin.Engine) {
	RegisterSiteFilesService(r, service.New())
}

// RegisterSiteFilesService 在站点根路径注册 robots.txt / sitemap.xml / rss.xml。
func RegisterSiteFilesService(r *gin.Engine, publicService api.PublicService) {
	handler := api.NewHandler(publicService)
	r.GET("/robots.txt", handler.Robots)
	r.GET("/sitemap.xml", handler.Sitemap)
	r.GET("/rss.xml", handler.RSS)
}

func RegisterService(apiGroup *gin.RouterGroup, publicService api.PublicService) {
	handler := api.NewHandler(publicService)
	blog := apiGroup.Group("blog")
	blog.GET("/categories", handler.BlogCategories)
	blog.GET("/tags", handler.BlogTags)
	blog.GET("/posts", handler.BlogPosts)
	blog.GET("/posts/:id", handler.BlogPost)
	blog.GET("/archive", handler.BlogArchive)

	blog.GET("/diary/tags", handler.DiaryTags)
	blog.GET("/diaries", handler.Diaries)
	blog.GET("/diaries/:id", handler.Diary)

	portfolio := blog.Group("portfolio")
	portfolio.GET("/categories", handler.PortfolioCategories)
	portfolio.GET("/items", handler.PortfolioItems)
	portfolio.GET("/items/:id", handler.PortfolioItem)

	blog.GET("/tools", handler.Tools)
	blog.GET("/tools/:id", handler.Tool)
	blog.GET("/bookmark/categories", handler.BookmarkCategories)
	blog.GET("/bookmarks", handler.Bookmarks)
	blog.GET("/profile", handler.Profile)
	blog.GET("/site", handler.Site)
	blog.GET("/stats", handler.Stats)
}
