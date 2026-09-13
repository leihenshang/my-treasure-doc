package router_test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"fastduck/treasure-doc/module/user/config"
	"fastduck/treasure-doc/module/user/global"
)

// 上传 / 备份 / 个人资料校验的接口契约测试，对应 doc/feature-inventory.md 的：
//   S-06 单图上传、S-07 媒体上传（多文件 + 去重）、M-BAK-01 备份、M-SET-01 个人资料校验
//
// 说明：这些 handler 走 module/user 的响应约定 —— **HTTP 恒为 200**，成败看业务码
// （SUCCESS = 0 / ERROR = 1），所以断言都读 `env.Code` 而不是 HTTP 状态码。

type uploadFile struct {
	field string
	name  string
	body  []byte
}

// pngBytes 最小 PNG：前 8 字节签名足以让 http.DetectContentType 判定为 image/png
func pngBytes(payload string) []byte {
	return append([]byte{0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a}, []byte(payload)...)
}

// mp4Bytes 最小 MP4：ftyp box 头，DetectContentType 判为 video/mp4
func mp4Bytes(payload string) []byte {
	return append([]byte{0x00, 0x00, 0x00, 0x18, 'f', 't', 'y', 'p', 'm', 'p', '4', '2'}, []byte(payload)...)
}

// plainList 解析裸数组响应（备份列表等，没有分页包装；
// 套件里既有的 dataList 是「解分页 list」语义，不适用）
func plainList(t *testing.T, env envelope) []map[string]any {
	t.Helper()
	var value []map[string]any
	if err := json.Unmarshal(env.Data, &value); err != nil {
		t.Fatalf("data 不是数组：%s", string(env.Data))
	}
	return value
}

func (s *testServer) doUpload(path string, files []uploadFile) (int, envelope) {
	s.t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	for _, file := range files {
		part, err := writer.CreateFormFile(file.field, file.name)
		if err != nil {
			s.t.Fatalf("构造 multipart 失败: %v", err)
		}
		if _, err := part.Write(file.body); err != nil {
			s.t.Fatalf("写入 multipart 失败: %v", err)
		}
	}
	if err := writer.Close(); err != nil {
		s.t.Fatalf("关闭 multipart 失败: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, path, &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()
	s.engine.ServeHTTP(rec, req)

	var env envelope
	if rec.Body.Len() > 0 {
		if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
			s.t.Fatalf("上传响应不是 JSON：%s", rec.Body.String())
		}
	}
	return rec.Code, env
}

// TestBlogImageUploadContract 对应 S-06
func TestBlogImageUploadContract(t *testing.T) {
	s := newTestServer(t)
	storedPath := regexp.MustCompile(`^/files/blog/[0-9a-f]{64}\.png$`)

	status, env := s.doUpload("/api/blog-mgr/uploads/images", []uploadFile{{field: "file", name: "封面.png", body: pngBytes("cover")}})
	if status != http.StatusOK || env.Code != 0 {
		t.Fatalf("上传 PNG = %d/%d %s", status, env.Code, env.Msg)
	}
	path, _ := dataObject(t, env)["path"].(string)
	if !storedPath.MatchString(path) {
		t.Fatalf("图片路径 = %q，想要 /files/blog/<sha256>.png", path)
	}
	name := strings.TrimPrefix(path, "/files/blog/")
	if _, err := os.Stat(filepath.Join("files", "blog", name)); err != nil {
		t.Fatalf("文件没有真正落盘：%v", err)
	}

	// 文本内容改成 .png 后缀：内容嗅探优先于后缀，必须拒绝（否则等于任何文件都能上传）
	if _, env = s.doUpload("/api/blog-mgr/uploads/images", []uploadFile{{field: "file", name: "fake.png", body: []byte("<html>not an image</html>")}}); env.Code == 0 {
		t.Fatalf("文本伪装成 .png 被接受：%s", env.Msg)
	}
	// 0 字节
	if _, env = s.doUpload("/api/blog-mgr/uploads/images", []uploadFile{{field: "file", name: "empty.png", body: nil}}); env.Code == 0 || !strings.Contains(env.Msg, "空") {
		t.Fatalf("空文件未被拒绝：code=%d msg=%s", env.Code, env.Msg)
	}
	// 超过 8MB
	oversize := append(pngBytes("big"), make([]byte, 8<<20)...)
	if _, env = s.doUpload("/api/blog-mgr/uploads/images", []uploadFile{{field: "file", name: "big.png", body: oversize}}); env.Code == 0 || !strings.Contains(env.Msg, "8MB") {
		t.Fatalf("超大图片未被拒绝：code=%d msg=%s", env.Code, env.Msg)
	}
	// 没有 file 字段
	if _, env = s.doUpload("/api/blog-mgr/uploads/images", nil); env.Code == 0 {
		t.Fatalf("缺少 file 字段却成功了")
	}
}

// TestBlogMediaUploadContract 对应 S-07
func TestBlogMediaUploadContract(t *testing.T) {
	s := newTestServer(t)

	status, env := s.doUpload("/api/blog-mgr/uploads/medias", []uploadFile{
		{field: "files", name: "图.png", body: pngBytes("image-a")},
		{field: "files", name: "片.mp4", body: mp4Bytes("video-a")},
	})
	if status != http.StatusOK || env.Code != 0 {
		t.Fatalf("上传媒体 = %d/%d %s", status, env.Code, env.Msg)
	}
	items, _ := dataObject(t, env)["list"].([]any)
	if len(items) != 2 {
		t.Fatalf("返回条目数 = %d，想要 2（且顺序与提交一致）", len(items))
	}
	first, _ := items[0].(map[string]any)
	second, _ := items[1].(map[string]any)
	if path, _ := first["path"].(string); !strings.HasSuffix(path, ".png") {
		t.Fatalf("第一个文件应为 png：%v", first["path"])
	}
	if path, _ := second["path"].(string); !strings.HasSuffix(path, ".mp4") {
		t.Fatalf("第二个文件应为 mp4：%v", second["path"])
	}

	// 同一份内容再传一次：命中 sha256 去重，标记 existed 且不重复写盘
	_, env = s.doUpload("/api/blog-mgr/uploads/medias", []uploadFile{{field: "files", name: "另一份.png", body: pngBytes("image-a")}})
	if env.Code != 0 {
		t.Fatalf("重复上传失败：%s", env.Msg)
	}
	repeat, _ := dataObject(t, env)["list"].([]any)
	entry, _ := repeat[0].(map[string]any)
	if entry["existed"] != true {
		t.Fatalf("重复内容未标记 existed：%v", entry)
	}
	if entry["path"] != first["path"] {
		t.Fatalf("重复内容路径不一致：%v vs %v", entry["path"], first["path"])
	}
	entries, err := os.ReadDir(filepath.Join("files", "blog"))
	if err != nil {
		t.Fatalf("读取上传目录失败: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("上传目录文件数 = %d，想要 2（去重后不应多出文件）", len(entries))
	}

	// 超出单次数量上限
	tooMany := make([]uploadFile, 0, 11)
	for index := 0; index < 11; index++ {
		tooMany = append(tooMany, uploadFile{field: "files", name: "f.png", body: pngBytes(string(rune('a' + index)))})
	}
	if _, env = s.doUpload("/api/blog-mgr/uploads/medias", tooMany); env.Code == 0 || !strings.Contains(env.Msg, "最多上传") {
		t.Fatalf("超量上传未被拒绝：code=%d msg=%s", env.Code, env.Msg)
	}

	// 文本文件（非白名单内容）必须拒绝
	if _, env = s.doUpload("/api/blog-mgr/uploads/medias", []uploadFile{{field: "files", name: "note.txt", body: []byte("plain text")}}); env.Code == 0 {
		t.Fatalf("文本文件被接受为媒体：%s", env.Msg)
	}
}

// TestProfileSettingValidation 对应 M-SET-01
func TestProfileSettingValidation(t *testing.T) {
	s := newTestServer(t)

	valid := `{"name":"Treasure","avatar":"T","role":"全栈","location":"中国","motto":"m","bio":"b","links":[{"id":"github","label":"GitHub","value":"github.com","url":"https://github.com/"}],"skills":[{"name":"Go","level":80,"group":"后端"}]}`
	if status, env := s.do(http.MethodPut, "/api/blog-mgr/profile", valid); status != http.StatusOK || env.Code != 0 {
		t.Fatalf("保存合法资料 = %d/%d %s", status, env.Code, env.Msg)
	}
	_, env := s.do(http.MethodGet, "/api/blog-mgr/profile", "")
	profile := dataObject(t, env)
	if profile["name"] != "Treasure" {
		t.Fatalf("资料未正确落库：%v", profile["name"])
	}
	if links, _ := profile["links"].([]any); len(links) != 1 {
		t.Fatalf("联系方式未回显：%v", profile["links"])
	}

	tests := map[string]string{
		"名称为空":     `{"name":"  ","links":[],"skills":[]}`,
		"联系方式缺ID":  `{"name":"T","links":[{"label":"GitHub","value":"v"}],"skills":[]}`,
		"联系方式缺label": `{"name":"T","links":[{"id":"github","value":"v"}],"skills":[]}`,
		"联系方式ID重复": `{"name":"T","links":[{"id":"a","label":"A","value":"v"},{"id":"a","label":"B","value":"v"}],"skills":[]}`,
		"联系方式URL超长": `{"name":"T","links":[{"id":"a","label":"A","value":"v","url":"` + strings.Repeat("x", 501) + `"}],"skills":[]}`,
	}
	for name, body := range tests {
		status, env := s.do(http.MethodPut, "/api/blog-mgr/profile", body)
		if status != http.StatusBadRequest || env.Code != 40001 {
			t.Fatalf("%s = %d/%d，想要 400/40001", name, status, env.Code)
		}
	}

	// 非法载荷不能污染已保存的数据
	_, env = s.do(http.MethodGet, "/api/blog-mgr/profile", "")
	if profile = dataObject(t, env); profile["name"] != "Treasure" {
		t.Fatalf("失败的写入污染了已有资料：%v", profile["name"])
	}
}

// TestBackupContract 对应 M-BAK-01
func TestBackupContract(t *testing.T) {
	s := newTestServer(t)

	// 未配置（GetConf 为空）时按「仅 SQLite 驱动支持备份」拒绝，而不是 panic
	_, env := s.do(http.MethodGet, "/api/blog-mgr/backups", "")
	if env.Code == 0 || !strings.Contains(env.Msg, "SQLite") {
		t.Fatalf("无配置时列出备份 = code %d msg %s，想要明确的驱动提示", env.Code, env.Msg)
	}

	// 配置为 sqlite + 临时备份目录后：创建 → 列表 → 下载
	backupDir := t.TempDir()
	previous := global.GetConf()
	global.Conf = &config.Config{
		Database: config.Database{Driver: config.DriverSQLite, Dsn: "file:" + filepath.Join(t.TempDir(), "treasure_doc.db") + "?_journal_mode=WAL"},
		Backup:   config.Backup{Dir: backupDir},
	}
	t.Cleanup(func() { global.Conf = previous })

	status, env := s.do(http.MethodPost, "/api/blog-mgr/backups", "")
	if status != http.StatusOK || env.Code != 0 {
		t.Fatalf("创建备份 = %d/%d %s", status, env.Code, env.Msg)
	}
	created := dataObject(t, env)
	name, _ := created["name"].(string)
	if !strings.HasPrefix(name, "treasure_doc_") || !strings.HasSuffix(name, ".db") {
		t.Fatalf("备份文件名 = %q，想要 treasure_doc_<时间戳>.db", name)
	}
	if size, _ := created["size"].(float64); size <= 0 {
		t.Fatalf("备份文件大小为 %v，想要 > 0", created["size"])
	}

	_, env = s.do(http.MethodGet, "/api/blog-mgr/backups", "")
	items := plainList(t, env)
	if len(items) != 1 || items[0]["name"] != name {
		t.Fatalf("备份列表 = %v，想要只含 %s", items, name)
	}

	// 下载：绕过 do()，因为响应体是文件而不是 JSON
	rec := httptest.NewRecorder()
	s.engine.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/blog-mgr/backups/"+name, nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("下载备份 = %d", rec.Code)
	}
	if disposition := rec.Header().Get("Content-Disposition"); !strings.Contains(disposition, name) {
		t.Fatalf("下载缺少 Content-Disposition：%q", disposition)
	}
	if !bytes.HasPrefix(rec.Body.Bytes(), []byte("SQLite format 3\x00")) {
		t.Fatalf("下载内容不是 SQLite 数据库文件（前 16 字节：%q）", rec.Body.Bytes()[:min(16, rec.Body.Len())])
	}

	// 目录穿越与不存在的文件都必须被拒绝
	if _, env = s.do(http.MethodGet, "/api/blog-mgr/backups/..", ""); env.Code == 0 {
		t.Fatalf("`..` 被当作合法备份名：%s", env.Msg)
	}
	if _, env = s.do(http.MethodGet, "/api/blog-mgr/backups/not-exist.db", ""); env.Code == 0 || !strings.Contains(env.Msg, "不存在") {
		t.Fatalf("不存在的备份未被拒绝：code=%d msg=%s", env.Code, env.Msg)
	}
}
