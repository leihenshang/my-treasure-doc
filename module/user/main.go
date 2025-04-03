package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	ginzap "github.com/gin-contrib/zap"

	"fastduck/treasure-doc/service/user/config"
	"fastduck/treasure-doc/service/user/global"
	"fastduck/treasure-doc/service/user/router"

	"github.com/gin-gonic/gin"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "c", config.DefaultConfig, "config file path")
	flag.Parse()
}

func main() {
	if destructFunc, err := global.InitModule(configFile); err != nil {
		fmt.Printf("failed to init modules:%v\n", err)
		os.Exit(1)
	} else {
		defer destructFunc()
	}

	if global.Conf.App.IsRelease() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	//记录全部的访问日志
	//把gin致命错误写入日志
	r.Use(ginzap.Ginzap(global.Zap, time.RFC3339, true)).Use(ginzap.RecoveryWithZap(global.Zap, true))
	router.InitRouter(r)
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", global.Conf.App.Port),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	global.Log.Info("service is started!", "address", s.Addr)
	global.Log.Error(s.ListenAndServe().Error())
}
