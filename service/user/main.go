package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"fastduck/treasure-doc/service/user/config"
	"fastduck/treasure-doc/service/user/global"
	"fastduck/treasure-doc/service/user/router"

	"github.com/gin-gonic/gin"
)

var cfgFile string

func init() {
	flag.StringVar(&cfgFile, "cfg", config.DefaultCfg, "config file path")
	flag.Parse()
}

func main() {
	if destructFunc, err := global.InitModule(cfgFile); err != nil {
		fmt.Printf("init module failed, err:%v\n", err)
		os.Exit(1)
	} else {
		defer destructFunc()
	}

	if global.CONFIG.App.IsRelease() {
		fmt.Println("设置模式为", gin.ReleaseMode)
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	//记录全部的访问日志
	// r.Use(ginzap.Ginzap(global.ZAP, time.RFC3339, true))
	//把gin致命错误写入日志
	//r.Use(ginzap.RecoveryWithZap(global.ZAP, true))
	router.InitRoute(r)
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", global.CONFIG.App.Port),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	global.ZAPSUGAR.Info("service is started!", "address", s.Addr)
	global.ZAPSUGAR.Error(s.ListenAndServe().Error())
}
