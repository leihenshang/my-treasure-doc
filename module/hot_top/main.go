package main

import (
	"fastduck/treasure-doc/module/hot_top/conf"
	"fastduck/treasure-doc/module/hot_top/hot"
	"fastduck/treasure-doc/module/hot_top/route"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "c", conf.DefaultConfig, "config file path")
	flag.Parse()
}

func main() {
	if err := conf.InitConf(configFile); err != nil {
		log.Fatalf("init config failed, err:%v\n", err)
		return
	}

	hot.NewHot(time.Hour).Start()
	r := gin.New()
	route.InitRoute(r).Use(route.MiddleWareCors())

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", 2025),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe().Error())
}
