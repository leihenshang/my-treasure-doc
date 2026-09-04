package router

import (
	"fastduck/treasure-doc/module/blog/api"
	"fastduck/treasure-doc/module/blog/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(apiGroup *gin.RouterGroup, dbProvider func() *gorm.DB) {
	RegisterService(apiGroup, service.New(dbProvider))
}

func RegisterService(apiGroup *gin.RouterGroup, publicService api.PublicService) {
	handler := api.NewHandler(publicService)
	blog := apiGroup.Group("blog")
	blog.GET("/categories", handler.BlogCategories)
	blog.GET("/tags", handler.BlogTags)
	blog.GET("/posts", handler.BlogPosts)
	blog.GET("/posts/:id", handler.BlogPost)

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
