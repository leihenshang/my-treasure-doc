package route

import (
	"fastduck/treasure-doc/module/hot_top/conf"
	"fastduck/treasure-doc/module/hot_top/hot"
	"fastduck/treasure-doc/module/hot_top/service"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func InitRoute(r *gin.Engine) *gin.Engine {
	route := r.Group("/").Use(MiddleWareCors())

	for _, v := range conf.UrlList {
		route.GET(string(v.Source), func(c *gin.Context) {
			resp, _ := hot.GetHotCache().Get(v.Source)
			c.JSON(http.StatusOK, resp)
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

	distDir := "./dist"
	if _, err := os.Stat(distDir); !os.IsNotExist(err) {
		r.StaticFS("/dist", gin.Dir(distDir, false))
		// 处理前端路由的SPA回退
		r.NoRoute(func(c *gin.Context) {
			c.File(filepath.Join(distDir, "index.html"))
		})
	}

	// 移除原来的重定向逻辑
	r.Any("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/dist")
	})

	return r
}
