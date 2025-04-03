package main

import (
	"flag"
	"fmt"
	"os"

	"fastduck/treasure-doc/module/user/config"
	"fastduck/treasure-doc/module/user/global"
	"fastduck/treasure-doc/module/user/internal/service"
)

const DefaultPwd = "12345678"

var user string
var pwd string
var cfgPath string

func init() {
	flag.StringVar(&user, "u", "", "user account")
	flag.StringVar(&pwd, "p", "", "user password")
	flag.StringVar(&cfgPath, "cfg", config.DefaultConfig, "config file path")
	flag.Parse()
}

func main() {
	if user == "" {
		fmt.Println("user cannot be empty")
		os.Exit(1)
	}
	if pwd == "" {
		pwd = DefaultPwd
	}

	if err := global.InitRestPwd(cfgPath); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := resetPwd(user, pwd); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("reset complete! user: %s, pwd: %s \n", user, pwd)
}

func resetPwd(user string, pwd string) error {
	return service.ResetPwd(user, pwd, pwd)
}
