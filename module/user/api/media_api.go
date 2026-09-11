package api

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"fastduck/treasure-doc/module/user/config"
	"fastduck/treasure-doc/module/user/data/response"
	"fastduck/treasure-doc/module/user/global"

	"github.com/gin-gonic/gin"
)

type mediaFile struct {
	Name string `json:"name"`
	// Path 为可直接引用的公开路径
	Path    string `json:"path"`
	Size    int64  `json:"size"`
	ModTime string `json:"modTime"`
}

// MediaApi 管理博客上传目录（/files/blog）下的文件。
type MediaApi struct{}

func NewMediaApi() *MediaApi { return &MediaApi{} }

// List 列出已上传的媒体文件，按修改时间倒序。
func (m *MediaApi) List(c *gin.Context) {
	dir := blogMediaDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			response.OkWithData(c, []mediaFile{})
			return
		}
		response.FailWithMessage(c, "读取上传目录失败")
		return
	}

	files := make([]mediaFile, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		files = append(files, mediaFile{
			Name:    entry.Name(),
			Path:    "/files/blog/" + entry.Name(),
			Size:    info.Size(),
			ModTime: info.ModTime().Format(time.RFC3339),
		})
	}
	sort.Slice(files, func(i, j int) bool { return files[i].ModTime > files[j].ModTime })
	response.OkWithData(c, files)
}

// References 统计指定媒体文件被哪些内容引用，供删除前提示。
func (m *MediaApi) References(c *gin.Context) {
	name, ok := safeMediaName(c)
	if !ok {
		return
	}
	references, total, err := countMediaReferences(name)
	if err != nil {
		response.FailWithMessage(c, "查询引用失败")
		return
	}
	response.OkWithData(c, gin.H{"total": total, "references": references})
}

// Delete 删除指定媒体文件；仅允许上传目录下的普通文件名。
func (m *MediaApi) Delete(c *gin.Context) {
	name, ok := safeMediaName(c)
	if !ok {
		return
	}
	path := filepath.Join(blogMediaDir(), name)
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		response.FailWithMessage(c, "文件不存在")
		return
	}
	if err := os.Remove(path); err != nil {
		response.FailWithMessage(c, "删除失败")
		return
	}
	response.OkWithData(c, gin.H{"deleted": true})
}

func blogMediaDir() string {
	return filepath.Join(config.FilesPath, "blog")
}

// maxMediaBatchNames 限制单次批量删除的文件数。
const maxMediaBatchNames = 200

// mediaBatchRequest 媒体库批量删除请求体，例如 {"names": ["a.png", "b.mp4"]}。
type mediaBatchRequest struct {
	Names []string `json:"names"`
}

// DeleteMany 批量删除媒体文件；不存在的文件会被跳过，返回实际删除数量。
func (m *MediaApi) DeleteMany(c *gin.Context) {
	var payload mediaBatchRequest
	if c.ShouldBindJSON(&payload) != nil {
		response.FailWithMessage(c, "请求参数格式错误")
		return
	}
	names, ok := normalizeMediaNames(payload.Names)
	if !ok {
		response.FailWithMessage(c, "文件名不合法")
		return
	}
	dir := blogMediaDir()
	var deleted int64
	for _, name := range names {
		path := filepath.Join(dir, name)
		info, err := os.Stat(path)
		if err != nil || info.IsDir() {
			// 已被删除或不是普通文件，跳过即可，保证批量删除幂等
			continue
		}
		if err := os.Remove(path); err != nil {
			response.FailWithMessage(c, fmt.Sprintf("删除 %s 失败", name))
			return
		}
		deleted++
	}
	response.OkWithData(c, gin.H{"deleted": deleted})
}

// safeMediaName 校验路径参数是上传目录下的普通文件名，避免越权访问。
func safeMediaName(c *gin.Context) (string, bool) {
	name := c.Param("name")
	if !validMediaName(name) {
		response.FailWithMessage(c, "文件名不合法")
		return "", false
	}
	return name, true
}

// validMediaName 判断是否为上传目录下的普通文件名（不含路径分隔符，长度受限）。
func validMediaName(name string) bool {
	return name != "" && len(name) <= 255 && name == filepath.Base(name) && name != "." && name != ".."
}

// normalizeMediaNames 去重并校验批量删除的文件名；空集合或超过上限视为非法。
func normalizeMediaNames(names []string) ([]string, bool) {
	result := make([]string, 0, len(names))
	seen := make(map[string]struct{}, len(names))
	for _, name := range names {
		name = strings.TrimSpace(name)
		if !validMediaName(name) {
			return nil, false
		}
		if _, exists := seen[name]; exists {
			continue
		}
		seen[name] = struct{}{}
		result = append(result, name)
	}
	if len(result) == 0 || len(result) > maxMediaBatchNames {
		return nil, false
	}
	return result, true
}

// mediaReference 描述某资源字段对当前文件的引用条数。
type mediaReference struct {
	Resource string `json:"resource"`
	Field    string `json:"field"`
	Count    int64  `json:"count"`
}

// mediaReferenceTarget 描述一个可能保存了 /files/blog/<文件名> 的表字段。
type mediaReferenceTarget struct {
	Table    string
	Column   string
	Resource string
	Field    string
}

// 内容里保存的是公开路径字符串（如封面、图标、头像、图集与正文 Markdown 内联图片），
// 因此按文件名做包含匹配即可覆盖全部引用形式。
var mediaReferenceTargets = []mediaReferenceTarget{
	{"td_blog_post", "content", "文章", "正文"},
	{"td_blog_diary", "content", "日记", "正文"},
	{"td_blog_portfolio_item", "cover", "作品", "封面"},
	{"td_blog_portfolio_item", "gallery", "作品", "图集"},
	{"td_blog_portfolio_item", "content", "作品", "正文"},
	{"td_blog_tool", "cover", "利器", "封面"},
	{"td_blog_tool", "content", "利器", "正文"},
	{"td_blog_bookmark", "icon", "收藏集", "图标"},
	{"td_blog_profile", "avatar", "关于我", "头像"},
	{"td_blog_profile", "bio", "关于我", "简介"},
	{"td_blog_profile", "links", "关于我", "链接图标"},
	{"td_blog_site", "home", "站点设置", "主页"},
	{"td_blog_site", "modules", "站点设置", "模块"},
	{"td_blog_site", "banner", "站点设置", "横幅"},
	{"td_blog_site", "footer", "站点设置", "页脚"},
	{"td_blog_site", "intro", "站点设置", "简介"},
}

// countMediaReferences 统计各资源字段对该文件名的引用数量。
// 只统计未软删除的记录——软删除内容不会再渲染，不构成真实引用。
func countMediaReferences(name string) ([]mediaReference, int64, error) {
	if global.Db == nil {
		return nil, 0, errors.New("database is not initialized")
	}
	pattern := "%" + escapeLikeValue(name) + "%"
	references := make([]mediaReference, 0, len(mediaReferenceTargets))
	var total int64
	for _, target := range mediaReferenceTargets {
		var count int64
		// 统一 CAST 成文本：PostgreSQL 下这些列是 jsonb，无法直接 LIKE。
		condition := "CAST(" + target.Column + " AS TEXT) LIKE ? ESCAPE '\\'"
		err := global.Db.Table(target.Table).
			Where("deleted_at IS NULL").
			Where(condition, pattern).
			Count(&count).Error
		if err != nil {
			return nil, 0, err
		}
		if count == 0 {
			continue
		}
		references = append(references, mediaReference{
			Resource: target.Resource,
			Field:    target.Field,
			Count:    count,
		})
		total += count
	}
	return references, total, nil
}

// escapeLikeValue 转义 LIKE 通配符，避免文件名中的 % 或 _ 造成误匹配。
func escapeLikeValue(value string) string {
	return strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(value)
}
