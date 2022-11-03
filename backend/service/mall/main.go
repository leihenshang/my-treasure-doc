package main

import (
	"database/sql"
	"fastduck/treasure-doc/service/mall/global"
	"fastduck/treasure-doc/service/mall/router"
	"fmt"
	"net/http"
	"time"

	// ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
)

func main() {

	//全局初始化
	global.GlobalInit("api")

	//同步写入日志
	defer global.ZAP.Sync()
	defer global.ZAPSUGAR.Sync()

	//关闭mysql
	db, _ := global.DB.DB()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	//设置运行模式
	if global.CONFIG.App.IsRelease() {
		fmt.Println("设置模式为", gin.ReleaseMode)
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	//记录全部的访问日志
	// r.Use(ginzap.Ginzap(global.ZAP, time.RFC3339, true))

	//把gin致命错误写入日志
	// r.Use(ginzap.RecoveryWithZap(global.ZAP, true))

	//初始化路由
	router.InitRoute(r)

	addr := fmt.Sprintf("%s:%d", global.CONFIG.App.Host, global.CONFIG.App.Port)
	//设置服务
	s := &http.Server{
		Addr:           addr,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	global.ZAPSUGAR.Info("service is started!\n", "http://", addr)
	global.ZAPSUGAR.Error(s.ListenAndServe().Error())
}
