package router

import "github.com/gin-gonic/gin"

func Register(g *gin.Engine) {
	rootPath := g.Group("/reverse_index")

	rootPath.POST("/index", nil)
	rootPath.GET("/search", nil)
}
