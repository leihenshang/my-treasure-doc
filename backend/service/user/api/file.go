package api

import (
	"fastduck/treasure-doc/service/user/response"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"path/filepath"
)

// FileUpload 文件上传
func FileUpload(c *gin.Context) {

	f, fHandler, fErr := c.Request.FormFile("file")
	if fErr != nil {
		response.Fail(c)
	}
	defer f.Close()

	dir, err := os.Getwd()
	if err != nil {
		response.FailWithMessage("获取目录失败", c)
		return
	}

	// todo 文件后缀名限制
	// todo 文件大小限制
	// todo 文件名字重新编码

	fName := filepath.Join(dir, "static", fHandler.Filename)
	targetFile, err := os.Create(fName)
	if err != nil {
		response.FailWithMessage("保存文件失败", c)
	}
	defer targetFile.Close()
	_, err = io.Copy(targetFile, f)
	if err != nil {
		response.FailWithMessage("复制文件到目标失败", c)
	}

	pathMap := make(map[string]string)
	pathMap["path"] = "/static/" + fHandler.Filename
	response.OkWithData(pathMap, c)
}
