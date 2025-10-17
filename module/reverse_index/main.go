package main

import (
	"fastduck/treasure-doc/module/reverse_index/index"
	"fastduck/treasure-doc/module/reverse_index/router"

	"github.com/gin-gonic/gin"
)

func main() {
	index.InitDict()
	gServer := gin.New()
	router.Register(gServer)

	_ = gServer.Run(":20251")
}
