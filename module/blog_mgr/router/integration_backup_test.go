package router_test

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	userapi "fastduck/treasure-doc/module/user/api"
	"fastduck/treasure-doc/module/user/config"
	"fastduck/treasure-doc/module/user/global"
	"fastduck/treasure-doc/module/user/router/middleware"
)

// setupBackupConf 把测试全局配置切到 sqlite + 临时备份目录/独立恢复库。
// 测试写到了 global.Db / global.Conf，必须逐用例还原，避免污染其他用例。
func setupBackupConf(t *testing.T) {
	t.Helper()
	previous := global.GetConf()
	dbFile := filepath.Join(t.TempDir(), "restored.db")
	global.Conf = &config.Config{
		Database: config.Database{
			Driver: config.DriverSQLite,
			Dsn:    "file:" + dbFile + "?_busy_timeout=5000&_journal_mode=WAL&_foreign_keys=on",
		},
		Backup: config.Backup{Dir: t.TempDir()},
	}
	t.Cleanup(func() {
		// 恢复会替换 global.Db（重开的新连接），这里一并关闭并复位，避免测试结束仍有连接占用文件
		if global.Db != nil {
			if sqlDb, err := global.Db.DB(); err == nil {
				_ = sqlDb.Close()
			}
			global.Db = nil
		}
		global.Conf = previous
	})
}

// TestFullBackupRoundTrip M-BAK-02：导出完整备份包 → 后台导入 → 文章/日志/图片恢复一致。
func TestFullBackupRoundTrip(t *testing.T) {
	s := newTestServer(t)
	setupBackupConf(t)

	// 数据准备：一篇文章、一篇日志、一张图片
	postID := s.createPost("backup-post", "备份文章", `,"publishStatus":"published"`)
	s.create("diaries", `{"publicId":"backup-diary","title":"备份日记","publishStatus":"published"}`)

	imgName := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa.png"
	imgDir := filepath.Join("files", "blog")
	if err := os.MkdirAll(imgDir, 0o755); err != nil {
		t.Fatalf("创建图片目录失败: %v", err)
	}
	if err := os.WriteFile(filepath.Join(imgDir, imgName), []byte("backup-image"), 0o644); err != nil {
		t.Fatalf("写图片失败: %v", err)
	}

	// 导出完整备份包
	archivePath, err := global.BuildBackupArchive(t.TempDir())
	if err != nil {
		t.Fatalf("导出完整备份包失败: %v", err)
	}
	archiveBytes, err := os.ReadFile(archivePath)
	if err != nil {
		t.Fatalf("读取备份包失败: %v", err)
	}

	// 删除图片，验证导入时一并恢复
	if err := os.Remove(filepath.Join(imgDir, imgName)); err != nil {
		t.Fatalf("删除图片失败: %v", err)
	}

	// 后台导入
	status, env := s.doUpload("/api/blog-mgr/backups/restore", []uploadFile{{field: "file", name: "backup.tar.gz", body: archiveBytes}})
	if status != http.StatusOK || env.Code != 0 {
		t.Fatalf("导入 = %d/%d %s", status, env.Code, env.Msg)
	}
	if !strings.Contains(string(env.Data), "restored") {
		t.Fatalf("导入响应缺少 restored 标记：%s", string(env.Data))
	}

	// 恢复后图片回到磁盘
	if _, err := os.Stat(filepath.Join(imgDir, imgName)); err != nil {
		t.Fatalf("图片未被恢复：%v", err)
	}

	// 恢复后数据库内容可读且完整（global.Db 已切换为恢复出的库）
	if status, env := s.do(http.MethodGet, "/api/blog-mgr/posts/"+postID, ""); status != http.StatusOK || env.Code != 0 {
		t.Fatalf("恢复后文章不可读：%d/%d %s", status, env.Code, env.Msg)
	}
	_, env = s.do(http.MethodGet, "/api/blog-mgr/diaries", "")
	if diaries := dataList(t, env); len(diaries) != 1 {
		t.Fatalf("恢复后日志数 = %d，想要 1", len(diaries))
	}
}

// TestRestoreRejectInvalid M-BAK-03：损坏/空/伪造备份被拒且无副作用。
func TestRestoreRejectInvalid(t *testing.T) {
	s := newTestServer(t)
	setupBackupConf(t)
	postID := s.createPost("keep-post", "完好", "")

	status, env := s.doUpload("/api/blog-mgr/backups/restore", []uploadFile{{field: "file", name: "fake.tar.gz", body: []byte("this is not a tarball")}})
	if status != http.StatusOK || env.Code == 0 || !strings.Contains(env.Msg, "gzip") {
		t.Fatalf("伪备份未被拒绝：%d/%d %s", status, env.Code, env.Msg)
	}
	if status, env := s.do(http.MethodGet, "/api/blog-mgr/posts/"+postID, ""); status != http.StatusOK || env.Code != 0 {
		t.Fatalf("拒绝损坏备份后原数据受影响：%d/%d %s", status, env.Code, env.Msg)
	}
	if _, env = s.doUpload("/api/blog-mgr/backups/restore", []uploadFile{{field: "file", name: "empty.tar.gz", body: nil}}); env.Code == 0 || !strings.Contains(env.Msg, "空") {
		t.Fatalf("空文件未被拒绝：code=%d msg=%s", env.Code, env.Msg)
	}
}

// TestRestoreRejectsTraversal M-BAK-04：备份包内嵌路径穿越条目必须被拒。
func TestRestoreRejectsTraversal(t *testing.T) {
	s := newTestServer(t)
	setupBackupConf(t)

	var buf bytes.Buffer
	if err := buildTarGz(&buf, []tarEntry{
		{name: "manifest.json", body: []byte(`{"schemaVersion":1}`)},
		{name: "treasure_doc.db", body: []byte(strings.Repeat("X", 4096))},
		{name: "../bad.txt", body: []byte("boom")},
	}); err != nil {
		t.Fatalf("构造 tar.gz 失败: %v", err)
	}

	_, env := s.doUpload("/api/blog-mgr/backups/restore", []uploadFile{{field: "file", name: "evil.tar.gz", body: buf.Bytes()}})
	if env.Code == 0 {
		t.Fatalf("路径穿越未被拒绝：%s", env.Msg)
	}
	if _, err := os.Stat(filepath.Join("..", "bad.txt")); err == nil {
		t.Fatalf("路径穿越条目被落地，安全漏洞")
	}
}

// TestNasExportToken M-BAK-05：NAS 机器令牌接口——缺失/错误 401，正确返回完整备份流。
func TestNasExportToken(t *testing.T) {
	s := newTestServer(t)
	previous := global.GetConf()
	global.Conf = &config.Config{
		Database: config.Database{Driver: config.DriverSQLite, Dsn: "file:" + filepath.Join(t.TempDir(), "nas.db") + "?_journal_mode=WAL"},
		Backup:   config.Backup{Dir: t.TempDir(), ApiToken: "supersecret-nas-token"},
	}
	t.Cleanup(func() { global.Conf = previous })

	// NAS 路由属于 user 模块 API 树，测试引擎需手动补挂
	s.engine.GET("/api/backup/export", middleware.BackupToken(), userapi.NewBackupApi().NasExportArchive)

	// 无令牌
	rec := httptest.NewRecorder()
	s.engine.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/backup/export", nil))
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("无令牌 = %d，想要 401", rec.Code)
	}
	// 错误令牌
	rec = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/backup/export", nil)
	req.Header.Set(middleware.BackupTokenHeader, "wrong-token")
	s.engine.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("错误令牌 = %d，想要 401", rec.Code)
	}
	// 正确令牌
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/backup/export", nil)
	req.Header.Set(middleware.BackupTokenHeader, "supersecret-nas-token")
	s.engine.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("正确令牌 = %d，想要 200（body=%s）", rec.Code, rec.Body.String())
	}
	if len(rec.Body.Bytes()) < 2 || rec.Body.Bytes()[0] != 0x1f || rec.Body.Bytes()[1] != 0x8b {
		t.Fatalf("NAS 导出响应不是 gzip 备份流（前 2 字节：%x）", rec.Body.Bytes())
	}
	if disp := rec.Header().Get("Content-Disposition"); !strings.Contains(disp, ".tar.gz") {
		t.Fatalf("NAS 导出缺少 .tar.gz 附件头：%q", disp)
	}
}

// TestBackupTokenManage M-BAK-06：NAS 备份令牌在后台可读写、持久化到 DB，并由 NAS 接口即时生效。
func TestBackupTokenManageAndNas(t *testing.T) {
	s := newTestServer(t)
	// 配置无兜底令牌：验证完全由后台 DB 驱动
	previous := global.GetConf()
	global.Conf = &config.Config{
		Database: config.Database{Driver: config.DriverSQLite, Dsn: "file:" + filepath.Join(t.TempDir(), "token.db") + "?_journal_mode=WAL"},
		Backup:   config.Backup{Dir: t.TempDir()},
	}
	t.Cleanup(func() {
		if global.Db != nil {
			if sqlDb, err := global.Db.DB(); err == nil {
				_ = sqlDb.Close()
			}
			global.Db = nil
		}
		global.Conf = previous
	})
	s.engine.GET("/api/backup/export", middleware.BackupToken(), userapi.NewBackupApi().NasExportArchive)

	// 默认无令牌 → NAS 401
	rec := httptest.NewRecorder()
	s.engine.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/backup/export", nil))
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("未配置令牌时 NAS = %d，想要 401", rec.Code)
	}

	// 后台写入令牌
	status, env := s.do(http.MethodPut, "/api/blog-mgr/backups/token", `{"token":"db-token-123"}`)
	if status != http.StatusOK || env.Code != 0 {
		t.Fatalf("保存令牌 = %d/%d %s", status, env.Code, env.Msg)
	}
	// 读回
	_, env = s.do(http.MethodGet, "/api/blog-mgr/backups/token", "")
	if got, _ := dataObject(t, env)["token"].(string); got != "db-token-123" {
		t.Fatalf("读回令牌 = %q，想要 db-token-123", got)
	}

	// 用 DB 令牌调用 NAS 接口 → 200
	rec = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/backup/export", nil)
	req.Header.Set(middleware.BackupTokenHeader, "db-token-123")
	s.engine.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("DB 令牌调 NAS = %d，想要 200", rec.Code)
	}

	// 后台清除令牌 → NAS 恢复 401
	if status, env = s.do(http.MethodPut, "/api/blog-mgr/backups/token", `{"token":""}`); status != http.StatusOK || env.Code != 0 {
		t.Fatalf("清除令牌 = %d/%d %s", status, env.Code, env.Msg)
	}
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/backup/export", nil)
	req.Header.Set(middleware.BackupTokenHeader, "db-token-123")
	s.engine.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("清除令牌后 NAS = %d，想要 401", rec.Code)
	}
}

// --- 测试工具：用标准库构造 gzip tar（供构造非法包/路径穿越场景） ---

type tarEntry struct {
	name string
	body []byte
}

// buildTarGz 把一个简单条目列表打包成 gzip tar 写入 dst。
func buildTarGz(dst *bytes.Buffer, entries []tarEntry) error {
	gz := gzip.NewWriter(dst)
	tw := tar.NewWriter(gz)
	for _, entry := range entries {
		hdr := &tar.Header{Name: entry.name, Mode: 0o644, Size: int64(len(entry.body))}
		if len(entry.body) == 0 && strings.HasSuffix(entry.name, "/") {
			hdr.Typeflag = tar.TypeDir
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
		if len(entry.body) > 0 {
			if _, err := tw.Write(entry.body); err != nil {
				return err
			}
		}
	}
	if err := tw.Close(); err != nil {
		return err
	}
	return gz.Close()
}
