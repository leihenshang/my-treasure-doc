package route

import (
	"fastduck/treasure-doc/module/hot_top/hot"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) *gin.Engine {
	route := r.Group("/").Use(MiddleWareCors())
	for k := range hot.UrlConfMap {
		route.GET(string(k), func(c *gin.Context) {
			resp, _ := hot.GetHotCache().Get(k)
			c.JSON(http.StatusOK, resp.HotData)
		})
	}

	route.GET("all", func(c *gin.Context) {
		c.JSON(http.StatusOK, hot.GetHotCache().GetAllMap())
	})

	
	return r
}
