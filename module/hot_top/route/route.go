package route

import (
	"fastduck/treasure-doc/module/hot_top/hot"
	"fastduck/treasure-doc/module/hot_top/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRoute(r *gin.Engine) *gin.Engine {
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

	route.GET("analysis-ds", func(c *gin.Context) {
		question := c.Query("question")
		answer, err := service.ThinkWithDeepSeek(question)
		if err != nil {
			c.JSON(http.StatusOK, err.Error())
			return
		}
		c.JSON(http.StatusOK, answer)
	})

	return r
}
