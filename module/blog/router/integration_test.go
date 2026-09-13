package router_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	blogservice "fastduck/treasure-doc/module/blog/internal/service"
	blogrouter "fastduck/treasure-doc/module/blog/router"
	mgrrouter "fastduck/treasure-doc/module/blog_mgr/router"
	"fastduck/treasure-doc/module/user/global"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 公开接口集成测试（对应 doc/feature-inventory.md 的 P-API-* / X-*）。
//
// 数据统一通过**管理端接口**写入（真实校验 + 真实落库），再断言公开接口的输出，
// 覆盖「后台写 → 前台读」这条最关键的契约：可见性、筛选、分页、降级结构。
// 鉴权中间件不在本用例范围内（见 module/user/router/middleware 下的独立测试）。

const (
	pastTime   = "2026-01-15T10:00:00Z"
	futureTime = "2030-01-15T10:00:00Z"
	pastDate   = "2026-01-15"
)

type server struct {
	t      *testing.T
	engine *gin.Engine
	db     *gorm.DB
}

func newServer(t *testing.T) *server {
	t.Helper()
	gin.SetMode(gin.TestMode)
	t.Chdir(t.TempDir())

	dsn := "file:" + filepath.Join(t.TempDir(), "public.db") + "?_busy_timeout=5000&_journal_mode=WAL&_foreign_keys=on&_synchronous=OFF"
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
	// 与生产一致地兜住 panic：单个 handler 崩溃不应该带崩整个测试进程
	engine.Use(gin.Recovery())
	blogrouter.RegisterSiteFiles(engine)
	blogrouter.RegisterService(engine.Group("/api"), blogservice.New())
	// 这里只借用管理端路由做数据夹具：Register 内部自行构造 service.New()，
	// 因此不需要（也不能）从 module/blog 导入 blog_mgr 的 internal 包。
	mgrrouter.Register(engine.Group("/api/blog-mgr"))
	return &server{t: t, engine: engine, db: db}
}

type apiEnvelope struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

func (s *server) do(method, path, body string) (int, apiEnvelope) {
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

	var env apiEnvelope
	if rec.Body.Len() > 0 {
		if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
			s.t.Fatalf("%s %s 响应不是 JSON：%s", method, path, rec.Body.String())
		}
	}
	return rec.Code, env
}

func (s *server) getRaw(path string) (int, string) {
	s.t.Helper()
	rec := httptest.NewRecorder()
	s.engine.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	return rec.Code, rec.Body.String()
}

func object(t *testing.T, env apiEnvelope) map[string]any {
	t.Helper()
	var value map[string]any
	if err := json.Unmarshal(env.Data, &value); err != nil {
		t.Fatalf("data 不是对象：%s", string(env.Data))
	}
	return value
}

func list(t *testing.T, env apiEnvelope) []map[string]any {
	t.Helper()
	var value []map[string]any
	if err := json.Unmarshal(env.Data, &value); err != nil {
		t.Fatalf("data 不是数组：%s", string(env.Data))
	}
	return value
}

func pageList(t *testing.T, env apiEnvelope) []map[string]any {
	t.Helper()
	var page struct {
		List []map[string]any `json:"list"`
	}
	if err := json.Unmarshal(env.Data, &page); err != nil {
		t.Fatalf("data 不是分页对象：%s", string(env.Data))
	}
	return page.List
}

func jsonBody(t *testing.T, value any) string {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("序列化失败: %v", err)
	}
	return string(data)
}

// create 通过管理端接口创建资源，返回 id；非 2xx 直接失败（说明夹具或接口契约有问题）
func (s *server) create(resource, body string) string {
	s.t.Helper()
	status, env := s.do(http.MethodPost, "/api/blog-mgr/"+resource, body)
	if status != http.StatusCreated {
		s.t.Fatalf("创建 %s 失败：%d %s", resource, status, env.Msg)
	}
	id, _ := object(s.t, env)["id"].(string)
	if id == "" {
		s.t.Fatalf("创建 %s 未返回 id：%s", resource, string(env.Data))
	}
	return id
}

func (s *server) createCategory(scope, slug, name string) string {
	s.t.Helper()
	status, env := s.do(http.MethodPost, "/api/blog-mgr/categories", fmt.Sprintf(`{"scope":%q,"slug":%q,"name":%q}`, scope, slug, name))
	if status != http.StatusCreated {
		s.t.Fatalf("创建分类失败：%d %s", status, env.Msg)
	}
	return slug // 分类对外的 ID 就是 slug
}

func (s *server) createTag(name string) string {
	s.t.Helper()
	return s.create("tags", fmt.Sprintf(`{"name":%q}`, name))
}

// TestPublicPostVisibilityAndFilters 对应 P-API-02 / P-API-03 / P-API-04
func TestPublicPostVisibilityAndFilters(t *testing.T) {
	s := newServer(t)
	s.createCategory("post", "tech", "技术")
	tagID := s.createTag("Go")

	s.create("posts", fmt.Sprintf(`{"slug":"published","title":"已发布文章","summary":"摘要","publishStatus":"published","publishedOn":%q,"publishedAt":%q,"categoryId":"tech","tagIds":[%q]}`, pastDate, pastTime, tagID))
	s.create("posts", `{"slug":"draft","title":"草稿文章","publishStatus":"draft"}`)
	s.create("posts", fmt.Sprintf(`{"slug":"future","title":"未来文章","publishStatus":"published","publishedOn":%q,"publishedAt":%q}`, pastDate, futureTime))
	// 注意：公开接口的 ID 是 slug（posts/portfolio-items/tools）或 publicId（diaries/bookmarks），
	// 与管理端详情返回的雪花 ID 不是同一套命名空间。
	const published = "published"

	// 列表只包含「已发布 + 发布时间已到」
	status, env := s.do(http.MethodGet, "/api/blog/posts", "")
	if status != http.StatusOK {
		t.Fatalf("公开列表 = %d %s", status, env.Msg)
	}
	items := pageList(t, env)
	if len(items) != 1 {
		t.Fatalf("公开列表应有 1 篇（草稿与未来文章不可见），实际 %d：%s", len(items), string(env.Data))
	}
	first := items[0]
	if first["id"] != published || first["title"] != "已发布文章" {
		t.Fatalf("列表首项不符合预期：%v", first)
	}
	if first["category"] != "tech" {
		t.Fatalf("列表 category = %v，想要 tech", first["category"])
	}
	if tags, ok := first["tags"].([]any); !ok || len(tags) != 1 || tags[0] != "Go" {
		t.Fatalf("列表 tags = %#v，想要 [Go]", first["tags"])
	}
	if first["date"] != pastDate {
		t.Fatalf("列表 date = %v，想要 %s", first["date"], pastDate)
	}

	// 详情：已发布 200，草稿/未来 404
	if status, env = s.do(http.MethodGet, "/api/blog/posts/"+published, ""); status != http.StatusOK {
		t.Fatalf("已发布详情 = %d %s", status, env.Msg)
	}
	detail := object(t, env)
	if detail["content"] == nil {
		t.Fatalf("详情缺少 content 字段：%s", string(env.Data))
	}
	beforeViews, _ := detail["views"].(float64)
	_, env = s.do(http.MethodGet, "/api/blog/posts/"+published, "")
	afterViews, _ := object(t, env)["views"].(float64)
	if afterViews <= beforeViews {
		t.Fatalf("浏览量未自增：%v → %v", beforeViews, afterViews)
	}

	// 草稿与未来文章的详情不可访问（公开详情按 slug 查询）
	for _, slug := range []string{"draft", "future"} {
		status, _ := s.do(http.MethodGet, "/api/blog/posts/"+slug, "")
		if status != http.StatusNotFound {
			t.Fatalf("%s 的公开详情 = %d，想要 404（不可见内容不能泄露）", slug, status)
		}
	}

	// 归档：按月份分组，只含已发布
	_, env = s.do(http.MethodGet, "/api/blog/archive", "")
	var groups []map[string]any
	if err := json.Unmarshal(env.Data, &groups); err != nil {
		t.Fatalf("归档 data 不是数组：%s", string(env.Data))
	}
	if len(groups) != 1 || groups[0]["month"] != "2026-01" {
		t.Fatalf("归档分组 = %v，想要 1 组 2026-01", groups)
	}
	if posts, ok := groups[0]["posts"].([]any); !ok || len(posts) != 1 {
		t.Fatalf("归档分组内文章数 = %#v，想要 1", groups[0]["posts"])
	}

	// 筛选：关键词 / 分类 / 标签
	_, env = s.do(http.MethodGet, "/api/blog/posts?keyword=%E5%B7%B2%E5%8F%91%E5%B8%83", "") // 已发布
	if items = pageList(t, env); len(items) != 1 {
		t.Fatalf("关键词筛选命中 %d，想要 1", len(items))
	}
	_, env = s.do(http.MethodGet, "/api/blog/posts?keyword=%25", "") // %
	if items = pageList(t, env); len(items) != 0 {
		t.Fatalf("关键词 %%%% 命中了 %d 条：LIKE 通配符未转义", len(items))
	}
	_, env = s.do(http.MethodGet, "/api/blog/posts?keyword=_", "")
	if items = pageList(t, env); len(items) != 0 {
		t.Fatalf("关键词 _ 命中了 %d 条：下划线通配符未转义", len(items))
	}
	_, env = s.do(http.MethodGet, "/api/blog/posts?categoryId=tech", "")
	if items = pageList(t, env); len(items) != 1 {
		t.Fatalf("分类筛选命中 %d，想要 1", len(items))
	}
	_, env = s.do(http.MethodGet, "/api/blog/posts?categoryId=not-exist", "")
	if items = pageList(t, env); len(items) != 0 {
		t.Fatalf("不存在的分类命中了 %d 条", len(items))
	}
	_, env = s.do(http.MethodGet, "/api/blog/posts?tag=Go", "")
	if items = pageList(t, env); len(items) != 1 {
		t.Fatalf("标签筛选命中 %d，想要 1", len(items))
	}

	// 分页上限：pageSize 超过 100 应被拒绝，而不是默默截断
	status, env = s.do(http.MethodGet, "/api/blog/posts?pageSize=101", "")
	if status == http.StatusOK {
		t.Fatalf("pageSize=101 被接受：应受 MaxPageSize=100 限制（%s）", string(env.Data))
	}
	status, env = s.do(http.MethodGet, "/api/blog/posts?sort=random", "")
	if status == http.StatusOK {
		t.Fatalf("非法 sort 被接受：应返回错误（%s）", string(env.Data))
	}
}

// TestPublicPostSorting 对应 P-API-02 的排序契约
func TestPublicPostSorting(t *testing.T) {
	s := newServer(t)
	s.create("posts", `{"slug":"early","title":"早的文章","publishStatus":"published","publishedOn":"2026-01-01","publishedAt":"2026-01-01T10:00:00Z"}`)
	s.create("posts", `{"slug":"late","title":"晚的文章","publishStatus":"published","publishedOn":"2026-02-01","publishedAt":"2026-02-01T10:00:00Z"}`)

	_, env := s.do(http.MethodGet, "/api/blog/posts?sort=desc", "")
	items := pageList(t, env)
	if len(items) != 2 || items[0]["title"] != "晚的文章" {
		t.Fatalf("desc 首项应为最新：%v", items)
	}
	_, env = s.do(http.MethodGet, "/api/blog/posts?sort=asc", "")
	items = pageList(t, env)
	if len(items) != 2 || items[0]["title"] != "早的文章" {
		t.Fatalf("asc 首项应为最早：%v", items)
	}

	// 分页：pageSize=1 时第二页应拿到另一篇
	_, env = s.do(http.MethodGet, "/api/blog/posts?pageSize=1&page=2&sort=desc", "")
	if items = pageList(t, env); len(items) != 1 || items[0]["title"] != "早的文章" {
		t.Fatalf("第二页内容不符合预期：%v", items)
	}
}

// TestPublicCatalogVisibility 对应 P-API-01 / P-API-05 / P-API-06 / P-API-07 / P-API-08
func TestPublicCatalogVisibility(t *testing.T) {
	s := newServer(t)
	s.createCategory("post", "tech", "技术")
	s.createCategory("post", "draft-only", "只有草稿用")
	s.createCategory("portfolio", "web", "网站")
	s.createCategory("bookmark", "tools", "工具站")

	tagID := s.createTag("Go")
	draftTagID := s.createTag("仅草稿标签")

	// 已发布文章：让 post 分类与标签「有用例支撑」
	s.create("posts", fmt.Sprintf(`{"slug":"cat-post","title":"带分类的文章","publishStatus":"published","publishedOn":%q,"publishedAt":%q,"categoryId":"tech","tagIds":[%q]}`, pastDate, pastTime, tagID))
	// 草稿文章：它的分类与标签不应出现在公开分类/标签里
	s.create("posts", fmt.Sprintf(`{"slug":"draft-post","title":"草稿文章","publishStatus":"draft","categoryId":"draft-only","tagIds":[%q]}`, draftTagID))
	// 公开 ID：日记/收藏集用 publicId，作品/利器用 slug
	const diary = "d1"
	const work = "w1"
	const ownTool = "t1"
	const linkTool = "t2"
	const bookmark = "b1"
	s.create("diaries", fmt.Sprintf(`{"publicId":%q,"title":"公开日记","publishStatus":"published","publishedOn":%q,"publishedAt":%q,"tagIds":[%q]}`, diary, pastDate, pastTime, tagID))
	s.create("diaries", `{"publicId":"d2","title":"草稿日记","publishStatus":"draft"}`)
	s.create("portfolio-items", fmt.Sprintf(`{"slug":%q,"title":"作品一","publishStatus":"published","publishedOn":%q,"publishedAt":%q,"categoryId":"web","techStack":["Go"],"gallery":["/files/blog/a.png"]}`, work, pastDate, pastTime))
	s.create("tools", fmt.Sprintf(`{"slug":%q,"kind":"own","name":"自研工具","developmentStatus":"可用","publishStatus":"published","publishedAt":%q}`, ownTool, pastTime))
	s.create("tools", fmt.Sprintf(`{"slug":%q,"kind":"link","name":"外部工具","url":"https://example.com","publishStatus":"published","publishedAt":%q}`, linkTool, pastTime))
	s.create("bookmarks", fmt.Sprintf(`{"publicId":%q,"title":"收藏一","url":"https://example.com","publishStatus":"published","publishedAt":%q,"tagIds":[%q],"categoryId":"tools"}`, bookmark, pastTime, tagID))

	// 日记：只出现公开的那条
	_, env := s.do(http.MethodGet, "/api/blog/diaries", "")
	items := pageList(t, env)
	if len(items) != 1 || items[0]["id"] != diary {
		t.Fatalf("公开日记列表 = %v，想要只含 d1", items)
	}
	if status, _ := s.do(http.MethodGet, "/api/blog/diaries/"+diary, ""); status != http.StatusOK {
		t.Fatalf("公开日记详情 = %d", status)
	}

	// 作品：列表 + 详情结构
	_, env = s.do(http.MethodGet, "/api/blog/portfolio/items", "")
	items = list(t, env)
	if len(items) != 1 || items[0]["id"] != work {
		t.Fatalf("公开作品列表 = %v", items)
	}
	_, env = s.do(http.MethodGet, "/api/blog/portfolio/items/"+work, "")
	detail := object(t, env)
	// 四个数组字段一视同仁：无数据时必须是 []，不能是 null（F4 回归点）
	for _, field := range []string{"gallery", "links", "metrics", "techStack"} {
		if _, ok := detail[field].([]any); !ok {
			t.Fatalf("作品详情 %s 应为数组，实际 %#v", field, detail[field])
		}
	}

	// 利器：自研走详情、外链带 url
	_, env = s.do(http.MethodGet, "/api/blog/tools", "")
	tools := list(t, env)
	if len(tools) != 2 {
		t.Fatalf("公开利器数 = %d，想要 2", len(tools))
	}
	status, env := s.do(http.MethodGet, "/api/blog/tools/"+ownTool, "")
	if status != http.StatusOK {
		t.Fatalf("自研利器详情 = %d", status)
	}
	if object(t, env)["type"] != "own" {
		t.Fatalf("自研利器 type 不是 own：%s", string(env.Data))
	}
	// 外链型没有详情页：后端只对 kind=own 提供详情（sitemap 也只收录自研），
	// 它的 url 在列表里返回，前端直接把整卡片渲染成外链。
	status, _ = s.do(http.MethodGet, "/api/blog/tools/"+linkTool, "")
	if status != http.StatusNotFound {
		t.Fatalf("外链利器详情 = %d，想要 404", status)
	}
	foundLink := false
	for _, tool := range tools {
		if tool["id"] == linkTool {
			foundLink = true
			if tool["url"] != "https://example.com" || tool["type"] != "link" {
				t.Fatalf("外链利器列表项不符合预期：%v", tool)
			}
		}
	}
	if !foundLink {
		t.Fatalf("列表里没有外链利器：%v", tools)
	}

	// 收藏集
	_, env = s.do(http.MethodGet, "/api/blog/bookmarks", "")
	bookmarks := list(t, env)
	if len(bookmarks) != 1 || bookmarks[0]["id"] != bookmark {
		t.Fatalf("公开收藏集 = %v", bookmarks)
	}

	// 分类/标签：只返回「已发布内容」用到的（草稿独占的分类/标签必须隐藏）
	_, env = s.do(http.MethodGet, "/api/blog/categories", "")
	if categories := list(t, env); len(categories) != 1 || categories[0]["id"] != "tech" {
		t.Fatalf("公开分类 = %v，想要只含 tech（draft-only 属于草稿，不应出现）", categories)
	}
	_, env = s.do(http.MethodGet, "/api/blog/tags", "")
	if tags := list(t, env); len(tags) != 1 || tags[0]["name"] != "Go" {
		t.Fatalf("公开标签 = %v，想要只含 Go（仅草稿标签不应出现）", tags)
	}
	_, env = s.do(http.MethodGet, "/api/blog/portfolio/categories", "")
	if categories := list(t, env); len(categories) != 1 || categories[0]["id"] != "web" {
		t.Fatalf("作品分类 = %v", categories)
	}
	_, env = s.do(http.MethodGet, "/api/blog/bookmark/categories", "")
	if categories := list(t, env); len(categories) != 1 || categories[0]["id"] != "tools" {
		t.Fatalf("收藏分类 = %v", categories)
	}
}

// TestPublicSiteProfileStats 对应 P-API-09 / P-API-10 / P-API-11
func TestPublicSiteProfileStats(t *testing.T) {
	s := newServer(t)

	// 个人资料：没有任何记录时也必须返回数组而不是 null（前端会直接 .map）
	_, env := s.do(http.MethodGet, "/api/blog/profile", "")
	profile := object(t, env)
	for _, field := range []string{"links", "skills"} {
		if _, ok := profile[field].([]any); !ok {
			t.Fatalf("空库 profile.%s 应为数组，实际 %#v", field, profile[field])
		}
	}

	// 站点设置：没有任何记录时返回完整默认模块（非严格归一化）
	_, env = s.do(http.MethodGet, "/api/blog/site", "")
	site := object(t, env)
	modules, ok := site["modules"].([]any)
	if !ok || len(modules) != 7 {
		t.Fatalf("空库 site.modules 应为 7 个固定模块，实际 %#v", site["modules"])
	}
	for _, field := range []string{"techStack", "milestones"} {
		if _, ok := site[field].([]any); !ok {
			t.Fatalf("空库 site.%s 应为数组，实际 %#v", field, site[field])
		}
	}

	// 直接写入含未知模块的记录：读取时应丢弃未知模块并补齐缺失模块
	unknownModules := `[{"id":"unknown","name":"X","path":"/X","marker":"X"}]`
	if err := s.db.Exec("INSERT INTO td_blog_site (id, site_key, name, slogan, intro, tech_stack, modules, milestones, home, footer, banner, maintenance_mode, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))",
		"site-1", "default", "站点", "", "", "[]", unknownModules, "[]", "{}", "{}", "{}", false).Error; err != nil {
		t.Fatalf("写入站点记录失败: %v", err)
	}
	_, env = s.do(http.MethodGet, "/api/blog/site", "")
	site = object(t, env)
	modules, _ = site["modules"].([]any)
	if len(modules) != 7 {
		t.Fatalf("含未知模块的站点记录归一化后有 %d 个模块，想要 7", len(modules))
	}
	for _, module := range modules {
		item, _ := module.(map[string]any)
		if item["id"] == "unknown" {
			t.Fatalf("未知模块未被丢弃：%s", string(env.Data))
		}
	}

	// 统计：只统计已发布内容
	s.create("posts", fmt.Sprintf(`{"slug":"p1","title":"公开文章","publishStatus":"published","publishedOn":%q,"publishedAt":%q}`, pastDate, pastTime))
	s.create("posts", `{"slug":"p2","title":"草稿文章","publishStatus":"draft"}`)
	s.create("diaries", `{"publicId":"d9","title":"草稿日记","publishStatus":"draft"}`)
	_, env = s.do(http.MethodGet, "/api/blog/stats", "")
	stats := object(t, env)
	if posts, _ := stats["posts"].(float64); int(posts) != 1 {
		t.Fatalf("公开统计 posts = %v，想要 1（草稿不应计入）", stats["posts"])
	}
	if diaries, _ := stats["diaries"].(float64); int(diaries) != 0 {
		t.Fatalf("公开统计 diaries = %v，想要 0", stats["diaries"])
	}
}

// TestPortfolioDetailArrayFields 对应 P-API-06 的数组契约（F4 回归用例）：
// 曾经 gallery/techStack 是 []，而 links/metrics 在无数据时是 null ——
// 根因是 decodeJSON 里 `json.Unmarshal("null", &result)` 会把预分配的切片重置为 nil，
// 现在统一归一成空切片。
func TestPortfolioDetailArrayFields(t *testing.T) {
	s := newServer(t)
	s.create("portfolio-items", fmt.Sprintf(`{"slug":"no-arrays","title":"无外链作品","publishStatus":"published","publishedOn":%q,"publishedAt":%q}`, pastDate, pastTime))

	_, env := s.do(http.MethodGet, "/api/blog/portfolio/items/no-arrays", "")
	detail := object(t, env)
	for _, field := range []string{"gallery", "links", "metrics", "techStack"} {
		if _, ok := detail[field].([]any); !ok {
			t.Fatalf("作品详情 %s 应为数组，实际 %#v", field, detail[field])
		}
	}
}

// TestSiteFiles 对应 X-01 / X-02 / X-03 / X-05
func TestSiteFiles(t *testing.T) {
	s := newServer(t)
	s.create("posts", fmt.Sprintf(`{"slug":"rss-post","title":"RSS 里的文章","summary":"摘要","publishStatus":"published","publishedOn":%q,"publishedAt":%q}`, pastDate, pastTime))
	s.create("posts", `{"slug":"rss-draft","title":"RSS 草稿","publishStatus":"draft"}`)

	// 说明：/ping 由 module/user/router 的内联 handler 注册，不可复用，因此不在本用例覆盖（清单 X-05）
	status, body := s.getRaw("/robots.txt")
	if status != http.StatusOK || !strings.Contains(body, "Sitemap:") {
		t.Fatalf("/robots.txt = %d %s", status, body)
	}

	status, body = s.getRaw("/sitemap.xml")
	if status != http.StatusOK {
		t.Fatalf("/sitemap.xml = %d", status)
	}
	for _, want := range []string{"/Blog/Portfolio", "/Blog/About", "rss-post"} {
		if !strings.Contains(body, want) {
			t.Fatalf("/sitemap.xml 缺少 %q：%s", want, body)
		}
	}
	if strings.Contains(body, "rss-draft") {
		t.Fatalf("/sitemap.xml 泄露了草稿：%s", body)
	}

	status, body = s.getRaw("/rss.xml")
	if status != http.StatusOK {
		t.Fatalf("/rss.xml = %d", status)
	}
	if !strings.Contains(body, "RSS 里的文章") {
		t.Fatalf("/rss.xml 缺少已发布文章：%s", body)
	}
	if strings.Contains(body, "RSS 草稿") {
		t.Fatalf("/rss.xml 泄露了草稿：%s", body)
	}
}
