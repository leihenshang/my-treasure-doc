package router

import (
	"fastduck/treasure-doc/module/reverse_index/api"

	"github.com/gin-gonic/gin"
)

func Register(g *gin.Engine) {
	rootPath := g.Group("/reverse_index")

	rootPath.POST("/index", api.Index)
	rootPath.GET("/search", api.Search)
	rootPath.GET("/list", api.List)
}
