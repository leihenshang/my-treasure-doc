package api

import (
	"crypto/md5"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"fastduck/treasure-doc/service/user/config"
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

	extension, err := fileCheck(fHandler)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	staticPath, err := checkSaveDir()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	md5Str, err := fileMd5(f)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	targetFileName := md5Str + "." + extension

	//重置文件指针，解决io.Copy一次后再次Copy文件为空的问题
	_, err = f.Seek(0, 0)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	targetFile, err := os.Create(filepath.Join(staticPath, targetFileName))
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
	pathMap["path"] = filepath.Join("/", getFilePath(), targetFileName)
	response.OkWithData(pathMap, c)
}

func getFilePath() string {
	return filepath.Join(config.FilePath, "statics")
}

func checkSaveDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("获取根目录失败")
	}

	staticPath := filepath.Join(dir, getFilePath())
	if _, err = os.Stat(staticPath); err != nil {
		if os.IsNotExist(err) {
			if err = os.Mkdir(staticPath, os.ModePerm); err != nil {
				return "", err
			}
		}
	}
	return staticPath, nil
}

func fileMd5(f io.Reader) (md5Str string, err error) {
	hash := md5.New()
	if _, err = io.Copy(hash, f); err != nil {
		err = fmt.Errorf("计算文件hash失败")
		return
	}
	md5Str = fmt.Sprintf("%x", hash.Sum(nil))
	return
}

func fileCheck(fHandler *multipart.FileHeader) (extension string, err error) {
	// 文件大小限制
	var size int64 = 10
	if fHandler.Size > (1024 * 1024 * size) {
		return extension, fmt.Errorf("文件大小超出限制: %d MB", size)
	}

	// 文件后缀名限制
	extension = getFileExtension(fHandler.Filename)
	if extension == "" {
		return extension, fmt.Errorf("获取文件后缀失败")
	}
	allowExtensions := map[string]struct{}{"jpg": {}, "png": {}, "bmp": {}, "gif": {}}
	if _, ok := allowExtensions[strings.ToLower(extension)]; !ok {
		return extension, fmt.Errorf("后缀不符合规则: %s", extension)
	}

	return extension, nil
}

func getFileExtension(fName string) string {
	index := strings.LastIndex(fName, ".")
	if index == -1 {
		return ""
	}

	return fName[index+1:]
}
