package api

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"fastduck/treasure-doc/service/user/data/response"
)

// FileUpload 文件上传
func FileUpload(c *gin.Context) {
	f, fHandler, fErr := c.Request.FormFile("file")
	if fErr != nil {
		response.Fail(c)
		return
	}
	defer f.Close()

	dir, err := os.Getwd()
	if err != nil {
		response.FailWithMessage("获取目录失败", c)
		return
	}

	// 文件大小限制
	var size int64 = 10
	if fHandler.Size > (1024 * 1024 * size) {
		response.FailWithMessage(fmt.Sprintf("文件大小超出限制: %d MB", size), c)
		return
	}

	// 文件后缀名限制
	extension := getFileExtension(fHandler.Filename)
	if extension == "" {
		response.FailWithMessage("获取文件后缀失败", c)
		return
	}
	allowExtensions := map[string]struct{}{"jpg": {}, "png": {}, "bmp": {}, "gif": {}}
	if _, ok := allowExtensions[strings.ToLower(extension)]; !ok {
		response.FailWithMessage(fmt.Sprintf("后缀不符合规则: %s", extension), c)
		return
	}

	// 文件名字重新编码
	hash := md5.New()
	if _, err = io.Copy(hash, f); err != nil {
		response.FailWithMessage("计算文件hash失败", c)
		return
	}
	md5Value := hash.Sum(nil)
	md5Str := fmt.Sprintf("%x", md5Value)
	targetFileName := md5Str + "." + extension

	//重置文件指针，解决io.Copy一次后再次Copy文件为空的问题
	f.Seek(0, 0)

	// todo 检查静态文件保存目录是否存在

	fName := filepath.Join(dir, "static", targetFileName)
	targetFile, err := os.Create(fName)
	if err != nil {
		response.FailWithMessage("保存文件失败", c)
		return
	}
	defer targetFile.Close()
	_, err = io.Copy(targetFile, f)
	if err != nil {
		response.FailWithMessage("复制文件到目标失败", c)
		return
	}

	pathMap := make(map[string]string)
	pathMap["path"] = "/static/" + targetFileName
	response.OkWithData(pathMap, c)
}

func getFileExtension(fName string) string {
	index := strings.LastIndex(fName, ".")
	if index == -1 {
		return ""
	}

	return fName[index+1:]
}
