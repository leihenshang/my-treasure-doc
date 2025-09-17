package main

import (
	"fastduck/treasure-doc/module/hot_top/hot"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	hot.Start()
	r := gin.New()
	InitRouter(r)
	r.Use(Cors())
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", 2025),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe().Error())
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
		c.Header("Access-Control-Allow-Methods", "*")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "3600")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func InitRouter(r *gin.Engine) {
	route := r.Group("/").Use(Cors())

	route.GET("/baidu", func(c *gin.Context) {
		resp, _ := hot.GetHotCache().Get(hot.SourceBaidu)
		c.JSON(http.StatusOK, resp)
	})

	route.GET("/ithome", func(c *gin.Context) {
		resp, _ := hot.GetHotCache().Get(hot.SourceITHome)
		c.JSON(http.StatusOK, resp)
	})

	route.GET("/weibo", func(c *gin.Context) {
		resp, _ := hot.GetHotCache().Get(hot.SourceWeibo)
		c.JSON(http.StatusOK, resp)
	})

	route.GET("/36kr", func(c *gin.Context) {
		resp, _ := hot.GetHotCache().Get(hot.Source36Kr)
		c.JSON(http.StatusOK, resp)
	})

	route.GET("/douyin", func(c *gin.Context) {
		resp, _ := hot.GetHotCache().Get(hot.SourceDouyin)
		c.JSON(http.StatusOK, resp)
	})

	route.GET("/bilibili", func(c *gin.Context) {
		resp, _ := hot.GetHotCache().Get(hot.SourceBilibili)
		c.JSON(http.StatusOK, resp)
	})

	route.GET("/sspai", func(c *gin.Context) {
		resp, _ := hot.GetHotCache().Get(hot.SourceSspai)
		c.JSON(http.StatusOK, resp)
	})

	route.GET("/zhihu", func(c *gin.Context) {
		resp, _ := hot.GetHotCache().Get(hot.SourceZhihu)
		c.JSON(http.StatusOK, resp)
	})

}
