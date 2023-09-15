package api

import (
	"fastduck/treasure-doc/service/user/response"
	"github.com/gin-gonic/gin"
)

// FileUpload 文件上传
func FileUpload(c *gin.Context) {
	pathMap := make(map[string]string)
	response.OkWithData(pathMap, c)
}
