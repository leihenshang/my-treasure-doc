package router_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"fastduck/treasure-doc/module/blog_mgr/internal/service"
	blogmgrrouter "fastduck/treasure-doc/module/blog_mgr/router"
	"fastduck/treasure-doc/module/user/global"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 集成测试：真实 SQLite + 真实 service + 真实 gin 路由，通过 HTTP 断言接口契约，
// 与 router_test.go（假 Manager，只验证路由注册）互补。
//
// 不覆盖的部分：鉴权/强制改密中间件（已有 module/user/router/middleware 下的独立测试）。

type testServer struct {
	t      *testing.T
	engine *gin.Engine
	db     *gorm.DB
}

func newTestServer(t *testing.T) *testServer {
	t.Helper()
	gin.SetMode(gin.TestMode)
	// 上传目录是常量 files/（相对进程工作目录），切到临时目录，避免测试往仓库里写文件
	t.Chdir(t.TempDir())

	// 测试库关掉同步写（_synchronous=OFF）纯为提速：真机配置仍需 WAL + 同步写
	dsn := "file:" + filepath.Join(t.TempDir(), "integration.db") + "?_busy_timeout=5000&_journal_mode=WAL&_foreign_keys=on&_synchronous=OFF"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("打开 SQLite 失败: %v", err)
	}
	tables := make([]interface{}, 0, len(global.TableMigrate))
	for _, table := range global.TableMigrate {
		tables = append(tables, table)
	}
	if err := db.AutoMigrate(tables...); err != nil {
		t.Fatalf("AutoMigrate 失败: %v", err)
	}

	global.Db = db
	t.Cleanup(func() {
		global.Db = nil
		if sqlDB, err := db.DB(); err == nil {
			_ = sqlDB.Close()
		}
	})

	engine := gin.New()
	// 生产环境用 gin.Default()（自带 Recovery）；测试里同样兜住 panic，
	// 否则任何一个 handler 的空指针都会把整个测试进程带崩，后面的用例全都看不到结果。
	engine.Use(gin.Recovery())
	blogmgrrouter.RegisterService(engine.Group("/api/blog-mgr"), service.New())
	return &testServer{t: t, engine: engine, db: db}
}

type envelope struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

func (s *testServer) do(method, path, body string) (int, envelope) {
	s.t.Helper()
	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, path, reader)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	s.engine.ServeHTTP(rec, req)

	var env envelope
	if rec.Body.Len() > 0 {
		if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
			s.t.Fatalf("%s %s 响应不是 JSON：%s", method, path, rec.Body.String())
		}
	}
	return rec.Code, env
}

func dataObject(t *testing.T, env envelope) map[string]any {
	t.Helper()
	var value map[string]any
	if err := json.Unmarshal(env.Data, &value); err != nil {
		t.Fatalf("data 不是对象：%s", string(env.Data))
	}
	return value
}

func dataList(t *testing.T, env envelope) []map[string]any {
	t.Helper()
	var page struct {
		List []map[string]any `json:"list"`
	}
	if err := json.Unmarshal(env.Data, &page); err != nil {
		t.Fatalf("data 不是分页对象：%s", string(env.Data))
	}
	return page.List
}

func paginationTotal(t *testing.T, env envelope) int {
	t.Helper()
	var page struct {
		Pagination struct {
			Total int `json:"total"`
		} `json:"pagination"`
	}
	if err := json.Unmarshal(env.Data, &page); err != nil {
		t.Fatalf("data 不是分页对象：%s", string(env.Data))
	}
	return page.Pagination.Total
}

func mustJSON(t *testing.T, value any) string {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("序列化失败: %v", err)
	}
	return string(data)
}

// createPost 建一篇文章并返回 id
func (s *testServer) createPost(slug, title, extra string) string {
	s.t.Helper()
	body := fmt.Sprintf(`{"slug":%q,"title":%q,"publishStatus":"draft"%s}`, slug, title, extra)
	status, env := s.do(http.MethodPost, "/api/blog-mgr/posts", body)
	if status != http.StatusCreated {
		s.t.Fatalf("创建文章 %s = %d：%s", slug, status, env.Msg)
	}
	id, _ := dataObject(s.t, env)["id"].(string)
	if id == "" {
		s.t.Fatalf("创建文章 %s 未返回 id", slug)
	}
	return id
}

// TestPostLifecycle 文章完整生命周期：创建 → 列表 → 详情 → 更新（乐观锁）→ 删除 → 回收站 → 恢复
func TestPostLifecycle(t *testing.T) {
	s := newTestServer(t)
	id := s.createPost("lifecycle", "第一版", "")

	// 列表默认只看未删除
	_, env := s.do(http.MethodGet, "/api/blog-mgr/posts?page=1&pageSize=20&deleted=exclude&sort=desc", "")
	if list := dataList(t, env); len(list) != 1 || list[0]["id"] != id {
		t.Fatalf("列表内容不符合预期：%v", list)
	}
	if total := paginationTotal(t, env); total != 1 {
		t.Fatalf("分页 total = %d，想要 1", total)
	}

	// 详情必须带 tagIds 数组：前端据此判断"能否安全编辑"（缺了会直接报错不让编辑）
	_, env = s.do(http.MethodGet, "/api/blog-mgr/posts/"+id, "")
	detail := dataObject(t, env)
	if _, ok := detail["tagIds"].([]any); !ok {
		t.Fatalf("详情缺少 tagIds 数组：%#v", detail["tagIds"])
	}
	version, _ := detail["version"].(float64)
	if version < 1 {
		t.Fatalf("详情 version = %v，想要 >= 1", detail["version"])
	}

	// 更新：带正确 version
	body := fmt.Sprintf(`{"slug":"lifecycle","title":"第二版","publishStatus":"published","version":%d}`, int(version))
	status, env := s.do(http.MethodPatch, "/api/blog-mgr/posts/"+id, body)
	if status != http.StatusOK {
		t.Fatalf("更新 = %d：%s", status, env.Msg)
	}
	if updated := dataObject(t, env); updated["title"] != "第二版" {
		t.Fatalf("更新响应 title = %v", updated["title"])
	}
	_, env = s.do(http.MethodGet, "/api/blog-mgr/posts/"+id, "")
	detail = dataObject(t, env)
	if detail["title"] != "第二版" {
		t.Fatalf("更新未落库，title = %v", detail["title"])
	}
	if next, _ := detail["version"].(float64); int(next) != int(version)+1 {
		t.Fatalf("version 未自增：%v → %v", version, next)
	}

	// 用过期 version 再更新应冲突（409 + 业务码 40900）
	stale := fmt.Sprintf(`{"slug":"lifecycle","title":"过期写入","publishStatus":"draft","version":%d}`, int(version))
	status, env = s.do(http.MethodPatch, "/api/blog-mgr/posts/"+id, stale)
	if status != http.StatusConflict || env.Code != 40900 {
		t.Fatalf("过期 version 更新 = %d/%d，想要 409/40900", status, env.Code)
	}

	// 软删除：默认列表看不到，回收站能看到
	if status, env = s.do(http.MethodDelete, "/api/blog-mgr/posts/"+id, ""); status != http.StatusOK {
		t.Fatalf("删除 = %d：%s", status, env.Msg)
	}
	_, env = s.do(http.MethodGet, "/api/blog-mgr/posts?deleted=exclude", "")
	if list := dataList(t, env); len(list) != 0 {
		t.Fatalf("删除后默认列表仍有 %d 条", len(list))
	}
	_, env = s.do(http.MethodGet, "/api/blog-mgr/posts?deleted=only", "")
	if list := dataList(t, env); len(list) != 1 {
		t.Fatalf("回收站应有 1 条，实际 %d", len(list))
	}

	// 恢复
	if status, env = s.do(http.MethodPost, "/api/blog-mgr/posts/"+id+"/restore", ""); status != http.StatusOK {
		t.Fatalf("恢复 = %d：%s", status, env.Msg)
	}
	_, env = s.do(http.MethodGet, "/api/blog-mgr/posts?deleted=exclude", "")
	if list := dataList(t, env); len(list) != 1 {
		t.Fatalf("恢复后应回到默认列表，实际 %d 条", len(list))
	}
}

// TestUpdateRequiresVersion 文章等资源更新必须带 version，缺失时按参数错误拒绝
func TestUpdateRequiresVersion(t *testing.T) {
	s := newTestServer(t)
	id := s.createPost("needs-version", "标题", "")

	status, env := s.do(http.MethodPatch, "/api/blog-mgr/posts/"+id, `{"slug":"needs-version","title":"无版本","publishStatus":"draft"}`)
	if status != http.StatusBadRequest || env.Code != 40001 {
		t.Fatalf("缺少 version 的更新 = %d/%d，想要 400/40001", status, env.Code)
	}
}

// TestBatchDeleteContract 批量删除：真实条数、忽略不存在/重复 ID、边界校验
func TestBatchDeleteContract(t *testing.T) {
	s := newTestServer(t)
	first := s.createPost("batch-1", "批量一", "")
	second := s.createPost("batch-2", "批量二", "")
	third := s.createPost("batch-3", "批量三", "")

	// 含不存在的 id：只删除真实存在的
	body := fmt.Sprintf(`{"ids":[%q,"not-exist-id"]}`, first)
	status, env := s.do(http.MethodPost, "/api/blog-mgr/posts/batch-delete", body)
	if status != http.StatusOK {
		t.Fatalf("批量删除 = %d：%s", status, env.Msg)
	}
	if deleted, _ := dataObject(t, env)["deleted"].(float64); int(deleted) != 1 {
		t.Fatalf("含不存在 id 时 deleted = %v，想要 1", deleted)
	}

	// 重复 id 去重：同一个 id 传两次只算一条
	body = fmt.Sprintf(`{"ids":[%q,%q]}`, second, second)
	_, env = s.do(http.MethodPost, "/api/blog-mgr/posts/batch-delete", body)
	if deleted, _ := dataObject(t, env)["deleted"].(float64); int(deleted) != 1 {
		t.Fatalf("重复 id 时 deleted = %v，想要 1", deleted)
	}

	// 已删除的 id 再删一次：不应再计入
	body = fmt.Sprintf(`{"ids":[%q]}`, second)
	_, env = s.do(http.MethodPost, "/api/blog-mgr/posts/batch-delete", body)
	if deleted, _ := dataObject(t, env)["deleted"].(float64); int(deleted) != 0 {
		t.Fatalf("重复删除已删记录 deleted = %v，想要 0", deleted)
	}

	// 空集合 / 超限 → 40001
	for name, payload := range map[string]string{
		"空集合": `{"ids":[]}`,
		"空 ID":  `{"ids":[""]}`,
	} {
		status, env = s.do(http.MethodPost, "/api/blog-mgr/posts/batch-delete", payload)
		if status != http.StatusBadRequest || env.Code != 40001 {
			t.Fatalf("%s = %d/%d，想要 400/40001", name, status, env.Code)
		}
	}
	tooMany := make([]string, 201)
	for i := range tooMany {
		tooMany[i] = fmt.Sprintf("id-%d", i)
	}
	status, env = s.do(http.MethodPost, "/api/blog-mgr/posts/batch-delete", mustJSON(t, map[string]any{"ids": tooMany}))
	if status != http.StatusBadRequest || env.Code != 40001 {
		t.Fatalf("超过 200 条 = %d/%d，想要 400/40001", status, env.Code)
	}

	// third 仍在（没有被前面的请求误删）
	_, env = s.do(http.MethodGet, "/api/blog-mgr/posts/"+third, "")
	if id := dataObject(t, env)["id"]; id != third {
		t.Fatalf("未被删除的记录查不到了：%v", id)
	}
}

// TestCategoryReferenceConflict 分类被内容引用时整体拒绝，不允许"删一半"
func TestCategoryReferenceConflict(t *testing.T) {
	s := newTestServer(t)

	status, env := s.do(http.MethodPost, "/api/blog-mgr/categories", `{"scope":"post","slug":"tech","name":"技术"}`)
	if status != http.StatusCreated {
		t.Fatalf("创建分类 = %d：%s", status, env.Msg)
	}
	used, _ := dataObject(t, env)["id"].(string)

	status, env = s.do(http.MethodPost, "/api/blog-mgr/categories", `{"scope":"post","slug":"unused","name":"没人用"}`)
	if status != http.StatusCreated {
		t.Fatalf("创建分类 = %d：%s", status, env.Msg)
	}
	unused, _ := dataObject(t, env)["id"].(string)

	// 文章引用 tech 分类
	s.createPost("uses-category", "引用分类的文章", `,"categoryId":"tech"`)

	status, env = s.do(http.MethodDelete, "/api/blog-mgr/categories/"+used, "")
	if status != http.StatusConflict || env.Code != 40900 {
		t.Fatalf("删除被引用分类 = %d/%d，想要 409/40900", status, env.Code)
	}

	// 批量删除：只要有一个被引用，整批都不删
	body := mustJSON(t, map[string]any{"ids": []string{used, unused}})
	status, env = s.do(http.MethodPost, "/api/blog-mgr/categories/batch-delete", body)
	if status != http.StatusConflict || env.Code != 40900 {
		t.Fatalf("批量删除含被引用分类 = %d/%d，想要 409/40900", status, env.Code)
	}
	_, env = s.do(http.MethodGet, "/api/blog-mgr/categories/"+unused, "")
	if id := dataObject(t, env)["id"]; id != unused {
		t.Fatalf("批量删除被拒后未被引用的分类也被删了：%v", id)
	}
}

// TestMediaLibraryContract 媒体库：列表 / 引用统计 / 批量删除（含不存在的文件）
func TestMediaLibraryContract(t *testing.T) {
	s := newTestServer(t)

	mediaDir := filepath.Join("files", "blog")
	if err := os.MkdirAll(mediaDir, 0o755); err != nil {
		t.Fatalf("创建上传目录失败: %v", err)
	}
	used := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa.png"
	orphan := "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb.png"
	for _, name := range []string{used, orphan} {
		if err := os.WriteFile(filepath.Join(mediaDir, name), []byte("fake-image"), 0o644); err != nil {
			t.Fatalf("写入测试文件失败: %v", err)
		}
	}

	// 一篇文章的正文引用其中一张图
	s.createPost("with-media", "带图文章", fmt.Sprintf(`,"content":"![x](/files/blog/%s)"`, used))

	// 列表：两个文件都在，path 是可直接引用的公开路径
	status, env := s.do(http.MethodGet, "/api/blog-mgr/medias", "")
	if status != http.StatusOK || env.Code != 0 {
		t.Fatalf("媒体列表 = %d/%d：%s", status, env.Code, env.Msg)
	}
	var files []map[string]any
	if err := json.Unmarshal(env.Data, &files); err != nil {
		t.Fatalf("媒体列表 data 不是数组：%s", string(env.Data))
	}
	if len(files) != 2 {
		t.Fatalf("媒体列表应有 2 个文件，实际 %d", len(files))
	}
	for _, file := range files {
		if file["path"] != "/files/blog/"+file["name"].(string) {
			t.Fatalf("媒体路径不符合约定：%v", file)
		}
	}

	// 引用统计：被引用的 1 处，孤立的 0 处
	_, env = s.do(http.MethodGet, "/api/blog-mgr/medias/"+used+"/references", "")
	references := dataObject(t, env)
	if total, _ := references["total"].(float64); int(total) != 1 {
		t.Fatalf("被引用文件 total = %v，想要 1（%s）", references["total"], string(env.Data))
	}
	_, env = s.do(http.MethodGet, "/api/blog-mgr/medias/"+orphan+"/references", "")
	if total, _ := dataObject(t, env)["total"].(float64); int(total) != 0 {
		t.Fatalf("孤立文件 total = %v，想要 0", total)
	}

	// 批量删除：不存在的文件被跳过，返回实际删除数量
	payload := mustJSON(t, map[string]any{"names": []string{orphan, "not-exist.png"}})
	_, env = s.do(http.MethodPost, "/api/blog-mgr/medias/batch-delete", payload)
	if deleted, _ := dataObject(t, env)["deleted"].(float64); int(deleted) != 1 {
		t.Fatalf("批量删除 deleted = %v，想要 1", deleted)
	}
	if _, err := os.Stat(filepath.Join(mediaDir, orphan)); !os.IsNotExist(err) {
		t.Fatalf("孤立文件未被删除：%v", err)
	}

	// 非法文件名（带路径分隔符）必须被拒绝
	status, env = s.do(http.MethodPost, "/api/blog-mgr/medias/batch-delete", `{"names":["../escape.png"]}`)
	if status != http.StatusOK || env.Code == 0 {
		t.Fatalf("非法文件名未被拒绝：%d/%d %s", status, env.Code, env.Msg)
	}
}

// TestSiteMaintenanceRoundTrip 站点设置：维护模式开关必须能写回 false（零值也要落库）
func TestSiteMaintenanceRoundTrip(t *testing.T) {
	s := newTestServer(t)

	_, env := s.do(http.MethodGet, "/api/blog-mgr/site", "")
	site := dataObject(t, env)
	if _, ok := site["maintenanceMode"]; !ok {
		t.Fatalf("站点设置缺少 maintenanceMode：%s", string(env.Data))
	}
	// 全新库的默认站点对象没有 name（已知问题 F3，见 doc/test-findings.md），
	// 这里补一个合法名称，让本用例聚焦"维护模式零值能否写回"。
	if name, _ := site["name"].(string); strings.TrimSpace(name) == "" {
		site["name"] = "Treasure Blog"
	}

	for _, want := range []bool{true, false, true} {
		site["maintenanceMode"] = want
		status, putEnv := s.do(http.MethodPut, "/api/blog-mgr/site", mustJSON(t, site))
		if status != http.StatusOK {
			t.Fatalf("保存站点设置 = %d：%s", status, putEnv.Msg)
		}
		_, env = s.do(http.MethodGet, "/api/blog-mgr/site", "")
		got, _ := dataObject(t, env)["maintenanceMode"].(bool)
		if got != want {
			t.Fatalf("maintenanceMode 写回不一致：写入 %v，读回 %v（零值被 GORM 跳过的典型症状）", want, got)
		}
	}
}

// TestSettingsDefaultIsRoundTrippable 对应 M-SET-03（F3 回归用例）：
// 全新库上 GET /site 与 GET /profile 返回的默认对象必须能原样 PUT 回去 ——
// 默认对象现在自带与 seed 一致的名称（validateSite / validateProfile 都要求 name 非空）。
func TestSettingsDefaultIsRoundTrippable(t *testing.T) {
	s := newTestServer(t)
	for _, setting := range []string{"site", "profile"} {
		_, env := s.do(http.MethodGet, "/api/blog-mgr/"+setting, "")
		status, putEnv := s.do(http.MethodPut, "/api/blog-mgr/"+setting, string(env.Data))
		if status != http.StatusOK {
			t.Fatalf("%s 默认对象无法原样保存：%d %s", setting, status, putEnv.Msg)
		}
	}
}

// TestListFiltersAndSort 列表筛选与排序
func TestListFiltersAndSort(t *testing.T) {
	s := newTestServer(t)

	first := s.createPost("filter-alpha", "Alpha 特殊词", `,"publishStatus":"published"`)
	time.Sleep(5 * time.Millisecond)
	second := s.createPost("filter-beta", "Beta 普通", "")
	time.Sleep(5 * time.Millisecond)
	third := s.createPost("filter-gamma", "Gamma 特殊词", "")

	// 关键词只命中带"特殊词"的两篇
	_, env := s.do(http.MethodGet, "/api/blog-mgr/posts?keyword=%E7%89%B9%E6%AE%8A%E8%AF%8D", "")
	if list := dataList(t, env); len(list) != 2 {
		t.Fatalf("关键词筛选命中 %d 条，想要 2", len(list))
	}

	// 状态筛选
	_, env = s.do(http.MethodGet, "/api/blog-mgr/posts?status=published", "")
	if list := dataList(t, env); len(list) != 1 || list[0]["id"] != first {
		t.Fatalf("状态筛选结果不符合预期：%v", list)
	}

	// 排序：desc 最新在前，asc 最早在前
	_, env = s.do(http.MethodGet, "/api/blog-mgr/posts?sort=desc", "")
	if list := dataList(t, env); len(list) != 3 || list[0]["id"] != third {
		t.Fatalf("desc 排序首条应为最新创建的：%v", list)
	}
	_, env = s.do(http.MethodGet, "/api/blog-mgr/posts?sort=asc", "")
	if list := dataList(t, env); len(list) != 3 || list[0]["id"] != first {
		t.Fatalf("asc 排序首条应为最早创建的：%v", list)
	}
	_ = second
}

// create 通用创建助手（夹具用），返回响应里的 id
func (s *testServer) create(resource, body string) string {
	s.t.Helper()
	status, env := s.do(http.MethodPost, "/api/blog-mgr/"+resource, body)
	if status != http.StatusCreated {
		s.t.Fatalf("创建 %s 失败：%d %s", resource, status, env.Msg)
	}
	id, _ := dataObject(s.t, env)["id"].(string)
	if id == "" {
		s.t.Fatalf("创建 %s 未返回 id：%s", resource, string(env.Data))
	}
	return id
}

// TestUpdateFieldsContract 对应 M-API-05：快捷改字段的白名单与版本行为
func TestUpdateFieldsContract(t *testing.T) {
	s := newTestServer(t)
	postID := s.createPost("fields-post", "快捷字段", "")

	// publishStatus：不要求 version，改完 version 自增
	status, env := s.do(http.MethodPatch, "/api/blog-mgr/posts/"+postID+"/fields", `{"publishStatus":"published"}`)
	if status != http.StatusOK {
		t.Fatalf("快捷改发布状态 = %d：%s", status, env.Msg)
	}
	_, env = s.do(http.MethodGet, "/api/blog-mgr/posts/"+postID, "")
	detail := dataObject(t, env)
	if detail["publishStatus"] != "published" {
		t.Fatalf("快捷改发布状态未生效：%v", detail["publishStatus"])
	}
	versionAfter, _ := detail["version"].(float64)

	// pinned 仅 posts/diaries 支持
	if status, env = s.do(http.MethodPatch, "/api/blog-mgr/posts/"+postID+"/fields", `{"pinned":true}`); status != http.StatusOK {
		t.Fatalf("快捷置顶 = %d：%s", status, env.Msg)
	}
	_, env = s.do(http.MethodGet, "/api/blog-mgr/posts/"+postID, "")
	detail = dataObject(t, env)
	if detail["pinned"] != true {
		t.Fatalf("快捷置顶未生效：%v", detail["pinned"])
	}
	if next, _ := detail["version"].(float64); next <= versionAfter {
		t.Fatalf("快捷改字段后 version 未自增：%v → %v", versionAfter, next)
	}

	toolID := s.create("tools", `{"slug":"fields-tool","kind":"own","name":"工具","developmentStatus":"可用","publishStatus":"draft"}`)
	for name, test := range map[string]struct{ resource, id, body string }{
		"非白名单字段":   {"posts", postID, `{"title":"改名"}`},
		"非法发布状态":   {"posts", postID, `{"publishStatus":"bogus"}`},
		"置顶只支持文章日记": {"tools", toolID, `{"pinned":true}`},
		"分类不走快捷接口":  {"categories", "cat-1", `{"publishStatus":"published"}`},
	} {
		status, env = s.do(http.MethodPatch, "/api/blog-mgr/"+test.resource+"/"+test.id+"/fields", test.body)
		if status != http.StatusBadRequest || env.Code != 40001 {
			t.Fatalf("%s = %d/%d，想要 400/40001", name, status, env.Code)
		}
	}
}

// TestCategorySlugRenameMigratesContent 对应 M-API-10：分类 scope 校验 + 改 slug 会同步内容
func TestCategorySlugRenameMigratesContent(t *testing.T) {
	s := newTestServer(t)

	// scope 是封闭枚举
	status, env := s.do(http.MethodPost, "/api/blog-mgr/categories", `{"scope":"food","slug":"hotpot","name":"火锅"}`)
	if status != http.StatusBadRequest || env.Code != 40001 {
		t.Fatalf("非法 scope = %d/%d，想要 400/40001", status, env.Code)
	}

	categoryID := s.create("categories", `{"scope":"post","slug":"tech","name":"技术"}`)
	postID := s.createPost("rename-post", "引用分类的文章", `,"categoryId":"tech"`)

	// 改 slug（对外 ID）后，已有内容的 categoryId 必须一起迁移，否则内容会指向不存在的分类
	body := mustJSON(t, map[string]any{"scope": "post", "slug": "tech-renamed", "name": "技术", "sortOrder": 0})
	status, env = s.do(http.MethodPatch, "/api/blog-mgr/categories/"+categoryID, body)
	if status != http.StatusOK {
		t.Fatalf("改分类 slug = %d：%s（曾因 mapDBError(nil) panic 变成 500，F5）", status, env.Msg)
	}
	_, env = s.do(http.MethodGet, "/api/blog-mgr/posts/"+postID, "")
	// 响应键名走 module/common/response.Normalize：只把开头的连续缩写小写，
	// 所以是 categoryID / publicID（前端 normalizeManageEntity 再转成 categoryId）
	if got := dataObject(t, env)["categoryID"]; got != "tech-renamed" {
		t.Fatalf("改 slug 后内容仍指向旧分类：%v", got)
	}
}

// TestTagRelationReplacement 对应 M-API-11：标签整体替换 + 引用不存在的标签细分错误
func TestTagRelationReplacement(t *testing.T) {
	s := newTestServer(t)
	tagA := s.create("tags", `{"name":"标签甲"}`)
	tagB := s.create("tags", `{"name":"标签乙"}`)
	postID := s.createPost("tag-post", "带标签的文章", fmt.Sprintf(`,"tagIds":[%q,%q]`, tagA, tagB))

	readTags := func() []any {
		_, env := s.do(http.MethodGet, "/api/blog-mgr/posts/"+postID, "")
		tags, _ := dataObject(t, env)["tagIds"].([]any)
		return tags
	}
	if tags := readTags(); len(tags) != 2 {
		t.Fatalf("创建后 tagIds = %v，想要 2 个", tags)
	}

	// 整体替换为 1 个
	version := func() int {
		_, env := s.do(http.MethodGet, "/api/blog-mgr/posts/"+postID, "")
		value, _ := dataObject(t, env)["version"].(float64)
		return int(value)
	}
	body := fmt.Sprintf(`{"slug":"tag-post","title":"带标签的文章","publishStatus":"draft","version":%d,"tagIds":[%q]}`, version(), tagA)
	if status, env := s.do(http.MethodPatch, "/api/blog-mgr/posts/"+postID, body); status != http.StatusOK {
		t.Fatalf("更新标签 = %d：%s", status, env.Msg)
	}
	if tags := readTags(); len(tags) != 1 || tags[0] != tagA {
		t.Fatalf("标签未被整体替换：%v", tags)
	}

	// 传空数组 = 清空
	body = fmt.Sprintf(`{"slug":"tag-post","title":"带标签的文章","publishStatus":"draft","version":%d,"tagIds":[]}`, version())
	if status, env := s.do(http.MethodPatch, "/api/blog-mgr/posts/"+postID, body); status != http.StatusOK {
		t.Fatalf("清空标签 = %d：%s", status, env.Msg)
	}
	if tags := readTags(); len(tags) != 0 {
		t.Fatalf("标签未被清空：%v", tags)
	}

	// 引用不存在的标签 → 细分业务码 40002（前端据此提示"请先创建标签"）
	body = fmt.Sprintf(`{"slug":"tag-post","title":"带标签的文章","publishStatus":"draft","version":%d,"tagIds":["not-exist-tag"]}`, version())
	status, env := s.do(http.MethodPatch, "/api/blog-mgr/posts/"+postID, body)
	if status != http.StatusBadRequest || env.Code != 40002 {
		t.Fatalf("引用不存在的标签 = %d/%d，想要 400/40002", status, env.Code)
	}
}

// TestToolValidationCodes 对应 M-API-12：利器两类形态的必填分支与细分业务码
func TestToolValidationCodes(t *testing.T) {
	s := newTestServer(t)
	tests := map[string]struct {
		body string
		code int
	}{
		"外链缺地址":   {`{"slug":"t-link","kind":"link","name":"外链","publishStatus":"draft"}`, 40003},
		"外链非 HTTPS": {`{"slug":"t-http","kind":"link","name":"外链","url":"http://example.com","publishStatus":"draft"}`, 40003},
		"自研缺开发状态": {`{"slug":"t-own","kind":"own","name":"自研","publishStatus":"draft"}`, 40004},
		"非法 kind":  {`{"slug":"t-bad","kind":"other","name":"X","publishStatus":"draft"}`, 40001},
	}
	for name, test := range tests {
		status, env := s.do(http.MethodPost, "/api/blog-mgr/tools", test.body)
		if status != http.StatusBadRequest || env.Code != test.code {
			t.Fatalf("%s = %d/%d，想要 400/%d", name, status, env.Code, test.code)
		}
	}
}

// TestStatsContract 对应 M-API-13：仪表盘统计按发布状态分类
func TestStatsContract(t *testing.T) {
	s := newTestServer(t)
	s.createPost("stats-draft", "草稿", "")
	postID := s.createPost("stats-published", "已发布", "")
	if status, env := s.do(http.MethodPatch, "/api/blog-mgr/posts/"+postID+"/fields", `{"publishStatus":"published"}`); status != http.StatusOK {
		t.Fatalf("设置发布状态 = %d：%s", status, env.Msg)
	}
	s.create("diaries", `{"publicId":"stats-diary","title":"日记","publishStatus":"draft"}`)
	s.create("categories", `{"scope":"post","slug":"stats-cat","name":"统计"}`)
	s.create("tags", `{"name":"统计标签"}`)

	status, env := s.do(http.MethodGet, "/api/blog-mgr/stats", "")
	if status != http.StatusOK {
		t.Fatalf("统计接口 = %d：%s", status, env.Msg)
	}
	stats := dataObject(t, env)
	posts, _ := stats["posts"].(map[string]any)
	if posts == nil {
		t.Fatalf("统计缺少 posts 分组：%s", string(env.Data))
	}
	if total, _ := posts["total"].(float64); int(total) != 2 {
		t.Fatalf("posts.total = %v，想要 2", posts["total"])
	}
	if published, _ := posts["published"].(float64); int(published) != 1 {
		t.Fatalf("posts.published = %v，想要 1", posts["published"])
	}
	if draft, _ := posts["draft"].(float64); int(draft) != 1 {
		t.Fatalf("posts.draft = %v，想要 1", posts["draft"])
	}
	diaries, _ := stats["diaries"].(map[string]any)
	if total, _ := diaries["total"].(float64); int(total) != 1 {
		t.Fatalf("diaries.total = %v，想要 1", diaries["total"])
	}
	if categories, _ := stats["categories"].(float64); int(categories) != 1 {
		t.Fatalf("categories = %v，想要 1", stats["categories"])
	}
}
