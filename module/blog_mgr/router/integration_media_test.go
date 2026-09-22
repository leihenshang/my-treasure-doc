package router_test

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	blogmodel "fastduck/treasure-doc/module/blog/data/model"
	"fastduck/treasure-doc/module/user/global"
)

// mediaFileRow 用于反解媒体列表里单条 JSON。
type mediaFileRow struct {
	Name           string `json:"name"`
	Path           string `json:"path"`
	Referenced     bool   `json:"referenced"`
	ReferenceCount int64  `json:"referenceCount"`
	OriginName     string `json:"originName"`
	Ext            string `json:"ext"`
	Mime           string `json:"mime"`
	Width          int    `json:"width"`
	Height         int    `json:"height"`
}

func listMediaRows(t *testing.T, env envelope) []mediaFileRow {
	t.Helper()
	var rows []mediaFileRow
	if err := json.Unmarshal(env.Data, &rows); err != nil {
		t.Fatalf("媒体列表 data 不是数组：%s", string(env.Data))
	}
	return rows
}

// writeMediaFile 在测试上传目录写一个 64 位 hex 文件名（模拟按内容 sha256 命名的上传文件）。
func writeMediaFile(t *testing.T, hexName, ext string) string {
	t.Helper()
	name := hexName + ext
	dir := filepath.Join("files", "blog")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("创建上传目录失败: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, name), []byte("media-"+hexName), 0o644); err != nil {
		t.Fatalf("写媒体文件失败: %v", err)
	}
	return name
}

// TestMediaMetadataFromUpload M-MED-02：上传的文件应在列表中携带元数据（原始名/后缀/类型）。
func TestMediaMetadataFromUpload(t *testing.T) {
	s := newTestServer(t)

	status, env := s.doUpload("/api/blog-mgr/uploads/medias", []uploadFile{{field: "file", name: "原图.png", body: pngBytes("cover-meta")}})
	if status != http.StatusOK || env.Code != 0 {
		t.Fatalf("上传媒体 = %d/%d %s", status, env.Code, env.Msg)
	}
	items, _ := dataObject(t, env)["list"].([]any)
	first, _ := items[0].(map[string]any)
	path, _ := first["path"].(string)
	name := strings.TrimPrefix(path, "/files/blog/")

	_, env = s.do(http.MethodGet, "/api/blog-mgr/medias", "")
	for _, row := range listMediaRows(t, env) {
		if row.Name != name {
			continue
		}
		if row.OriginName != "原图.png" {
			t.Fatalf("原始文件名未登记：%+v", row)
		}
		if row.Ext != ".png" || row.Mime != "image/png" {
			t.Fatalf("元数据 ext/mime 错误：%+v", row)
		}
		return
	}
	t.Fatalf("上传的文件未出现在媒体列表：%s", name)
}

// TestMediaMetadataAndReferenceMarking M-MED-01：元数据同步 + 列表引用标注 + ref 筛选 + 删除清理元数据。
func TestMediaMetadataAndReferenceMarking(t *testing.T) {
	s := newTestServer(t)

	usedName := writeMediaFile(t, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", ".png")
	orphanName := writeMediaFile(t, "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", ".png")

	// 启动时或流程外文件，需要一次目录同步补齐元数据表
	if err := global.SyncMediaCatalog(); err != nil {
		t.Fatalf("同步媒体目录失败: %v", err)
	}

	// 元数据表应登记两个文件
	var mediaCount int64
	if err := s.db.Model(&blogmodel.Media{}).Count(&mediaCount).Error; err != nil {
		t.Fatalf("统计元数据表失败: %v", err)
	}
	if mediaCount != 2 {
		t.Fatalf("元数据表记录数 = %d，想要 2", mediaCount)
	}

	// 一篇文章引用 used 文件
	s.createPost("media-ref", "引用图文的文章", `,"content":"![a](/files/blog/`+usedName+`)"`)

	// 列表：used 被引用，orphan 未被引用
	_, env := s.do(http.MethodGet, "/api/blog-mgr/medias", "")
	rows := listMediaRows(t, env)
	byName := map[string]mediaFileRow{}
	for _, row := range rows {
		byName[row.Name] = row
	}
	if used := byName[usedName]; !used.Referenced || used.ReferenceCount < 1 {
		t.Fatalf("被引用文件标记错误：%+v", used)
	}
	if orphan := byName[orphanName]; orphan.Referenced {
		t.Fatalf("孤儿文件被误标为已引用：%+v", orphan)
	}

	// 筛选：仅未引用 → 只剩 orphan
	_, env = s.do(http.MethodGet, "/api/blog-mgr/medias?ref=unreferenced", "")
	if rows := listMediaRows(t, env); len(rows) != 1 || rows[0].Name != orphanName {
		t.Fatalf("unreferenced 筛选 = %+v，想要只剩 %s", rows, orphanName)
	}
	// 筛选：仅已引用 → 只剩 used
	_, env = s.do(http.MethodGet, "/api/blog-mgr/medias?ref=referenced", "")
	if rows := listMediaRows(t, env); len(rows) != 1 || rows[0].Name != usedName {
		t.Fatalf("referenced 筛选 = %+v，想要只剩 %s", rows, usedName)
	}

	// 删除 orphan（未引用，允许）后，元数据行同步移除
	_, env = s.do(http.MethodDelete, "/api/blog-mgr/medias/"+orphanName, "")
	if env.Code != 0 {
		t.Fatalf("删除孤儿文件失败：%s", env.Msg)
	}
	if err := s.db.Model(&blogmodel.Media{}).Where("name = ?", orphanName).First(&blogmodel.Media{}).Error; err == nil {
		t.Fatalf("删除文件后元数据行仍存在")
	}
}
