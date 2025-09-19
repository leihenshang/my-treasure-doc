package main

import (
	"fastduck/treasure-doc/module/hot_top/conf"
	"fastduck/treasure-doc/module/hot_top/hot"
	"fastduck/treasure-doc/module/hot_top/route"
	"flag"
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
	genEngine := gin.New()
	genEngine.Use(gin.Logger())
	route.InitRoute(genEngine).Use(route.MiddleWareCors())

	addr := conf.GetConf().App.GetAddr()

	s := &http.Server{
		Addr:         addr,
		Handler:      genEngine,
		ReadTimeout:  360 * time.Second,
		WriteTimeout: 360 * time.Second,
		// MaxHeaderBytes: 1 << 20,
	}
	log.Printf("service is started! address: [http://%s]\n", addr)
	log.Fatal(s.ListenAndServe().Error())
}
