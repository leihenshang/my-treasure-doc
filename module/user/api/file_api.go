package api

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"fastduck/treasure-doc/module/user/config"
	"fastduck/treasure-doc/module/user/data/response"

	"github.com/gin-gonic/gin"
)

const (
	// 单图场景（封面、主页宣传图）的大小上限
	maxBlogImageSize = 8 << 20
	// Markdown 编辑器内图片/视频的大小上限与单次数量上限
	maxBlogMediaSize  = 50 << 20
	maxBlogMediaCount = 10
	// 用于内容嗅探的头部长度，与 http.DetectContentType 的要求一致
	detectBufferSize = 512
)

// 允许的图片类型：内容类型（由文件内容嗅探得到）到存储后缀的映射
var imageMediaExtensions = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/gif":  ".gif",
	"image/webp": ".webp",
	"image/bmp":  ".bmp",
}

var videoMediaExtensions = map[string]string{
	"video/mp4":        ".mp4",
	"video/webm":       ".webm",
	"video/quicktime":  ".mov",
	"video/avi":        ".avi",
	"video/x-msvideo":  ".avi",
	"video/x-matroska": ".mkv",
	"video/mpeg":       ".mpeg",
}

var blogMediaExtensions = mergeMediaExtensions(imageMediaExtensions, videoMediaExtensions)

// 内容无法被嗅探（application/octet-stream）时按文件名后缀兜底，
// 兜底结果仍必须命中白名单内容类型；能被识别为其他类型的文件不会被兜底放行。
var contentTypeBySuffix = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".gif":  "image/gif",
	".webp": "image/webp",
	".bmp":  "image/bmp",
	".mp4":  "video/mp4",
	".m4v":  "video/mp4",
	".webm": "video/webm",
	".mov":  "video/quicktime",
	".avi":  "video/avi",
	".mkv":  "video/x-matroska",
	".mpeg": "video/mpeg",
	".mpg":  "video/mpeg",
}

func mergeMediaExtensions(sources ...map[string]string) map[string]string {
	merged := make(map[string]string)
	for _, source := range sources {
		for contentType, extension := range source {
			merged[contentType] = extension
		}
	}
	return merged
}

// blogMedia 描述一个已保存（或命中去重）的媒体文件
type blogMedia struct {
	// Path 公开访问路径，例如 /files/blog/<sha256>.png
	Path string
	// Name 上传时的原始文件名，便于调用方展示
	Name string
	Size int64
	// Existed 为 true 表示服务端已存在相同内容，未重复写入磁盘
	Existed bool
}

// UploadBlogImage 保存博客后台上传的图片（封面、主页宣传图等单图场景）。
func UploadBlogImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.FailWithMessage(c, "请选择图片文件")
		return
	}
	media, err := saveBlogMedia(file, imageMediaExtensions, maxBlogImageSize)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.OkWithData(c, gin.H{"path": media.Path})
}

// UploadBlogMedias 保存 Markdown 编辑器上传的图片/视频，支持一次提交多个文件。
// 文件按内容 sha256 命名，相同内容只保留一份；返回顺序与提交顺序一致。
func UploadBlogMedias(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		response.FailWithMessage(c, "解析上传内容失败")
		return
	}
	files := form.File["files"]
	if len(files) == 0 {
		files = form.File["file"]
	}
	if len(files) == 0 {
		response.FailWithMessage(c, "请选择要上传的文件")
		return
	}
	if len(files) > maxBlogMediaCount {
		response.FailWithMessage(c, fmt.Sprintf("一次最多上传 %d 个文件", maxBlogMediaCount))
		return
	}

	list := make([]gin.H, 0, len(files))
	for _, file := range files {
		media, err := saveBlogMedia(file, blogMediaExtensions, maxBlogMediaSize)
		if err != nil {
			response.FailWithMessage(c, fmt.Sprintf("[%s]%s", file.Filename, err.Error()))
			return
		}
		list = append(list, gin.H{
			"path":    media.Path,
			"name":    media.Name,
			"size":    media.Size,
			"existed": media.Existed,
		})
	}
	response.OkWithData(c, gin.H{"list": list})
}

// saveBlogMedia 校验并保存单个上传文件。文件以内容 sha256 加后缀命名，
// 因此相同内容重复上传时直接复用已有文件，不再写盘。
func saveBlogMedia(file *multipart.FileHeader, allowed map[string]string, maxSize int64) (blogMedia, error) {
	media := blogMedia{Name: file.Filename, Size: file.Size}
	if file.Size <= 0 {
		return media, errors.New("文件内容为空")
	}
	if file.Size > maxSize {
		return media, fmt.Errorf("文件大小不能超过 %dMB", maxSize>>20)
	}

	source, err := file.Open()
	if err != nil {
		return media, errors.New("读取文件失败")
	}
	defer source.Close()

	extension, err := detectMediaExtension(source, file.Filename, allowed)
	if err != nil {
		return media, err
	}
	if _, err = source.Seek(0, io.SeekStart); err != nil {
		return media, errors.New("读取文件失败")
	}

	hasher := sha256.New()
	if _, err = io.Copy(hasher, source); err != nil {
		return media, errors.New("读取文件失败")
	}
	name := hex.EncodeToString(hasher.Sum(nil)) + extension

	directory := filepath.Join(config.FilesPath, "blog")
	if err = os.MkdirAll(directory, 0o755); err != nil {
		return media, errors.New("创建文件目录失败")
	}
	media.Path = "/files/blog/" + name
	target := filepath.Join(directory, name)
	if info, statErr := os.Stat(target); statErr == nil && !info.IsDir() {
		media.Existed = true
		return media, nil
	}

	if _, err = source.Seek(0, io.SeekStart); err != nil {
		return media, errors.New("读取文件失败")
	}
	output, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		// 相同内容并发上传时可能已被其他请求写入，直接复用
		if os.IsExist(err) {
			media.Existed = true
			return media, nil
		}
		return media, errors.New("保存文件失败")
	}
	defer output.Close()
	if _, err = io.Copy(output, source); err != nil {
		_ = os.Remove(output.Name())
		return media, errors.New("保存文件失败")
	}
	return media, nil
}

// detectMediaExtension 通过文件内容嗅探类型，无法嗅探时按后缀兜底，返回存储后缀。
func detectMediaExtension(source io.ReadSeeker, filename string, allowed map[string]string) (string, error) {
	header := make([]byte, detectBufferSize)
	n, err := io.ReadFull(source, header)
	if err != nil && err != io.ErrUnexpectedEOF {
		return "", errors.New("读取文件失败")
	}
	contentType := http.DetectContentType(header[:n])
	if extension, ok := allowed[contentType]; ok {
		return extension, nil
	}
	// 只有内容完全无法识别时才信任后缀，避免被改名的文本/脚本文件蒙混过关
	if contentType == "application/octet-stream" {
		if fallback, ok := contentTypeBySuffix[strings.ToLower(filepath.Ext(filename))]; ok {
			if extension, exists := allowed[fallback]; exists {
				return extension, nil
			}
		}
	}
	return "", errors.New("不支持的文件格式")
}
