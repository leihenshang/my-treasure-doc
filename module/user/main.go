package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	ginzap "github.com/gin-contrib/zap"

	"fastduck/treasure-doc/module/user/config"
	"fastduck/treasure-doc/module/user/global"
	"fastduck/treasure-doc/module/user/internal/service"
	"fastduck/treasure-doc/module/user/router"

	"github.com/gin-gonic/gin"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "c", config.DefaultConfig, "config file path")
	flag.Parse()
}

func main() {
	// 子命令：treasure-doc resetpwd <新密码>，用于重置默认管理员密码
	if args := flag.Args(); len(args) > 0 && args[0] == "resetpwd" {
		runResetPwd(args[1:])
		return
	}

	if destructFunc, err := global.InitModule(configFile); err != nil {
		fmt.Printf("failed to init modules:%v\n", err)
		os.Exit(1)
	} else {
		defer destructFunc()
	}

	if global.GetConf().App.IsRelease() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	//记录全部的访问日志
	//把gin致命错误写入日志
	r.Use(ginzap.Ginzap(global.Zap, time.RFC3339, true)).Use(ginzap.RecoveryWithZap(global.Zap, true))
	router.InitRouter(r)
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", global.GetConf().App.Port),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	global.Log.Info("service is started!", "address", s.Addr)
	global.Log.Error(s.ListenAndServe().Error())
}

// runResetPwd 重置默认管理员密码，新密码需满足 8-16 位规则
func runResetPwd(args []string) {
	if len(args) != 1 {
		fmt.Printf("用法: %s resetpwd <新密码>\n", os.Args[0])
		os.Exit(2)
	}

	destructFunc, err := global.InitModule(configFile)
	if err != nil {
		fmt.Printf("初始化失败:%v\n", err)
		os.Exit(1)
	}
	defer destructFunc()

	if err := service.NewUserService().ResetDefaultAdminPassword(args[0]); err != nil {
		fmt.Printf("重置密码失败:%v\n", err)
		os.Exit(1)
	}
	fmt.Printf("账号 [%s] 的密码已重置\n", service.DefaultAdminAccount)
}
