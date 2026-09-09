package api

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"fastduck/treasure-doc/module/user/config"
	"fastduck/treasure-doc/module/user/data/response"

	"github.com/gin-gonic/gin"
)

const maxBlogImageSize = 8 << 20

var imageExtensions = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/gif":  ".gif",
	"image/webp": ".webp",
}

// UploadBlogImage 保存博客后台上传的图片，并返回公开访问路径。
func UploadBlogImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.FailWithMessage(c, "请选择图片文件")
		return
	}
	if file.Size <= 0 || file.Size > maxBlogImageSize {
		response.FailWithMessage(c, "图片大小需在 8MB 以内")
		return
	}

	source, err := file.Open()
	if err != nil {
		response.FailWithMessage(c, "读取图片失败")
		return
	}
	defer source.Close()

	header := make([]byte, 512)
	n, err := io.ReadFull(source, header)
	if err != nil && err != io.ErrUnexpectedEOF {
		response.FailWithMessage(c, "读取图片失败")
		return
	}
	contentType := http.DetectContentType(header[:n])
	extension, ok := imageExtensions[contentType]
	if !ok && strings.EqualFold(strings.TrimSpace(mime.TypeByExtension(filepath.Ext(file.Filename))), "image/webp") {
		contentType, extension, ok = "image/webp", ".webp", true
	}
	if !ok || contentType == "" {
		response.FailWithMessage(c, "仅支持 JPG、PNG、GIF 和 WebP 图片")
		return
	}
	if _, err := source.Seek(0, io.SeekStart); err != nil {
		response.FailWithMessage(c, "读取图片失败")
		return
	}

	directory := filepath.Join(config.FilesPath, "blog")
	if err := os.MkdirAll(directory, 0o755); err != nil {
		response.FailWithMessage(c, "创建图片目录失败")
		return
	}
	name, err := randomFileName(extension)
	if err != nil {
		response.FailWithMessage(c, "生成图片名称失败")
		return
	}
	output, err := os.OpenFile(filepath.Join(directory, name), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		response.FailWithMessage(c, "保存图片失败")
		return
	}
	defer output.Close()
	if _, err := io.Copy(output, source); err != nil {
		_ = os.Remove(output.Name())
		response.FailWithMessage(c, "保存图片失败")
		return
	}
	response.OkWithData(c, gin.H{"path": "/files/blog/" + name})
}

func randomFileName(extension string) (string, error) {
	buffer := make([]byte, 16)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	return strings.ToLower(hex.EncodeToString(buffer)) + extension, nil
}
