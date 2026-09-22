package router_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	blogmgrrouter "fastduck/treasure-doc/module/blog_mgr/router"
	"fastduck/treasure-doc/module/user/router/middleware"
)

// pubReq 以发布令牌调任意路径（method 可指定），返回 (status, envelope)。
func pubReq(s *testServer, method, path, body, token string) (int, envelope) {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set(middleware.PublishTokenHeader, token)
	}
	rec := httptest.NewRecorder()
	s.engine.ServeHTTP(rec, req)
	var env envelope
	if rec.Body.Len() > 0 {
		if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
			s.t.Fatalf("响应不是 JSON：%s", rec.Body.String())
		}
	}
	return rec.Code, env
}

func TestPublishMetaAndOverride(t *testing.T) {
	s := newTestServer(t)
	if status, env := s.do(http.MethodPut, "/api/blog-mgr/publish/token", `{"token":"psec"}`); status != http.StatusOK || env.Code != 0 {
		t.Fatalf("设发布令牌 = %d/%d", status, env.Code)
	}
	pub := s.engine.Group("/api/publish")
	pub.Use(middleware.RequirePublishToken())
	blogmgrrouter.RegisterPublish(pub)

	// 准备分类与标签
	s.create("categories", `{"scope":"post","slug":"tech","name":"技术"}`)
	s.create("tags", `{"name":"标签A"}`)

	// 发布端分类/标签列表
	_, env := pubReq(s, http.MethodGet, "/api/publish/categories?scope=post", "", "psec")
	var cats []map[string]any
	if err := json.Unmarshal(env.Data, &cats); err != nil {
		t.Fatalf("分类列表解析失败：%s", env.Data)
	}
	found := false
	for _, c := range cats {
		if c["slug"] == "tech" {
			found = true
		}
	}
	if !found {
		t.Fatalf("发布端分类列表缺少 tech：%s", env.Data)
	}
	_, env = pubReq(s, http.MethodGet, "/api/publish/tags", "", "psec")
	if !strings.Contains(string(env.Data), "标签A") {
		t.Fatalf("发布端标签列表缺少标签A：%s", env.Data)
	}

	// by-slug 查询：不存在 → exists=false
	_, env = pubReq(s, http.MethodGet, "/api/publish/posts/by-slug/nope", "", "psec")
	if got, _ := dataObject(t, env)["exists"].(bool); got {
		t.Fatalf("不存在文章 exists 应为 false：%s", env.Data)
	}

	// 创建发布 → 覆盖 → 版本自增 + 历史
	status, env := pubReq(s, http.MethodPost, "/api/publish/posts", `{"slug":"pub-over","title":"最初","content":"v1"}`, "psec")
	if status != http.StatusCreated {
		t.Fatalf("创建发布 = %d：%s", status, env.Msg)
	}
	postID, _ := dataObject(t, env)["id"].(string)

	// 覆盖
	if status, env = pubReq(s, http.MethodPut, "/api/publish/posts/pub-over", `{"slug":"pub-over","title":"覆盖后","content":"v2"}`, "psec"); status != http.StatusOK || env.Code != 0 {
		t.Fatalf("发布覆盖 = %d/%d %s", status, env.Code, env.Msg)
	}
	_, env = s.do(http.MethodGet, "/api/blog-mgr/posts/"+postID, "")
	detail := dataObject(t, env)
	if detail["title"] != "覆盖后" {
		t.Fatalf("覆盖后标题 = %v", detail["title"])
	}

	// 历史列表应至少含「最初」与「覆盖后」
	_, env = s.do(http.MethodGet, "/api/blog-mgr/history/posts/"+postID, "")
	var list []map[string]any
	if err := json.Unmarshal(env.Data, &list); err != nil {
		t.Fatalf("历史列表解析失败：%s", env.Data)
	}
	if len(list) < 2 {
		t.Fatalf("历史条数 = %d，想要 >=2（创建+覆盖）", len(list))
	}
	summaries := map[string]bool{}
	for _, item := range list {
		summaries[item["summary"].(string)] = true
	}
	if !summaries["最初"] || !summaries["覆盖后"] {
		t.Fatalf("历史应含最初与覆盖后：%v", list)
	}

	// 恢复「最初」版本（seq 由列表反查）
	restoreSeq := -1
	for _, item := range list {
		if item["summary"] == "最初" {
			restoreSeq = int(item["seq"].(float64))
		}
	}
	status, env = s.do(http.MethodPost, "/api/blog-mgr/history/posts/"+postID+"/"+strconv.Itoa(restoreSeq)+"/restore", "")
	if status != http.StatusOK || env.Code != 0 {
		t.Fatalf("恢复历史 = %d/%d %s", status, env.Code, env.Msg)
	}
	_, env = s.do(http.MethodGet, "/api/blog-mgr/posts/"+postID, "")
	if got, _ := dataObject(t, env)["title"].(string); got != "最初" {
		t.Fatalf("恢复后标题 = %q，想要 最初", got)
	}
	// 恢复会再记一条新历史
	_, env = s.do(http.MethodGet, "/api/blog-mgr/history/posts/"+postID, "")
	var list2 []map[string]any
	if err := json.Unmarshal(env.Data, &list2); err != nil {
		t.Fatalf("恢复后历史解析失败：%s", env.Data)
	}
	if len(list2) != len(list)+1 {
		t.Fatalf("恢复后历史条数 = %d，想要 %d", len(list2), len(list)+1)
	}
}
