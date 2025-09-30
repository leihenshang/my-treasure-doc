package route

import (
	"fastduck/treasure-doc/module/hot_top/conf"
	"fastduck/treasure-doc/module/hot_top/hot"
	hotcache "fastduck/treasure-doc/module/hot_top/hot/hot_cache"
	"fastduck/treasure-doc/module/hot_top/service"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func InitRoute(r *gin.Engine) *gin.Engine {
	route := r.Group("/").Use(MiddleWareCors())

	// register all the routes from conf.HotConfListMap
	for _, v := range conf.HotConfListMap {
		route.GET(string(v.Source), func(c *gin.Context) {
			resp, _ := hotcache.GetHotMemCache().Get(v.Source)
			c.JSON(http.StatusOK, resp)
		})
	}

	route.POST("refresh", func(c *gin.Context) {
		source := c.PostForm("source")
		res := hot.GetHot().RefreshHotCache(conf.Source(source))
		c.JSON(http.StatusOK, res)
	})

	// get all the hot data
	route.GET("all", func(c *gin.Context) {
		c.JSON(http.StatusOK, hotcache.GetHotMemCache().GetAllMap())
	})

	// use deepseek to analyze hot data
	route.GET("analysis-ds", func(c *gin.Context) {
		question := c.Query("question")
		answer, err := service.ThinkWithDeepSeek(question)
		if err != nil {
			c.JSON(http.StatusOK, err.Error())
			return
		}
		c.JSON(http.StatusOK, answer)
	})

	// frontend static files
	distDir := "./dist"
	if _, err := os.Stat(distDir); !os.IsNotExist(err) {
		r.StaticFS("/dist", gin.Dir(distDir, false))
		// 处理前端路由的SPA回退
		r.NoRoute(func(c *gin.Context) {
			c.File(filepath.Join(distDir, "index.html"))
		})
	}

	// if direct access the "/", redirect to /dist
	r.Any("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/dist")
	})

	return r
}
