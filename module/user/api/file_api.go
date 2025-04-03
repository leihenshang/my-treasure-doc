package api

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"fastduck/treasure-doc/module/user/config"
	"fastduck/treasure-doc/module/user/data/response"
	"fastduck/treasure-doc/module/user/global"
	"fastduck/treasure-doc/module/user/utils"
)

type FileApi struct {
}

func NewFileApi() *FileApi {
	return &FileApi{}
}

// FileUpload 文件上传
func (fa *FileApi) FileUpload(c *gin.Context) {
	f, fHandler, fErr := c.Request.FormFile("file")
	if fErr != nil {
		response.Fail(c)
		return
	}
	defer f.Close()

	extension, err := fileCheck(fHandler)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}

	staticPath, err := utils.GenDir(getFilePath())
	if err != nil {
		global.Log.Errorf("failed to generate directory,error:%v", err)
		response.FailWithMessage(c, "创建目录失败")
		return
	}

	md5Str, err := utils.FileMd5(f)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	targetFileName := md5Str + "." + extension

	//重置文件指针，解决io.Copy一次后再次Copy文件为空的问题
	_, err = f.Seek(0, 0)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	targetFile, err := os.Create(filepath.Join(staticPath, targetFileName))
	if err != nil {
		response.FailWithMessage(c, "保存文件失败")
		return
	}
	defer targetFile.Close()
	_, err = io.Copy(targetFile, f)
	if err != nil {
		response.FailWithMessage(c, "复制文件到目标失败")
		return
	}

	pathMap := make(map[string]string)
	pathMap["path"] = filepath.Join("/", getFilePath(), targetFileName)
	response.OkWithData(c, pathMap)
}

func getFilePath() string {
	return filepath.Join(config.FilesPath, "uploads")
}

func fileCheck(fHandler *multipart.FileHeader) (extension string, err error) {
	// 文件大小限制
	var size int64 = 10
	if fHandler.Size > (1024 * 1024 * size) {
		return extension, fmt.Errorf("文件大小超出限制: %d MB", size)
	}

	// 文件后缀名限制
	extension = utils.GetFileExtension(fHandler.Filename)
	if extension == "" {
		return extension, fmt.Errorf("获取文件后缀失败")
	}
	allowExtensions := map[string]struct{}{"jpg": {}, "png": {}, "bmp": {}, "gif": {}}
	if _, ok := allowExtensions[strings.ToLower(extension)]; !ok {
		return extension, fmt.Errorf("后缀不符合规则: %s", extension)
	}

	return extension, nil
}
