package router_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	blogmgrrouter "fastduck/treasure-doc/module/blog_mgr/router"
	"fastduck/treasure-doc/module/user/router/middleware"
)

// publish 以发布令牌调 /api/publish 树下的接口，返回 (status, envelope)。
func publish(s *testServer, token string, path, body string) (int, envelope) {
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
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

// TestPublishApi M-PUB-01：机器令牌发布接口 —— 令牌校验 + 强制发布 + 默认分类 + 校验失败。
func TestPublishApi(t *testing.T) {
	s := newTestServer(t)

	// 后台配置发布令牌（走管理端点落库）
	status, env := s.do(http.MethodPut, "/api/blog-mgr/publish/token", `{"token":"publish-secret"}`)
	if status != http.StatusOK || env.Code != 0 {
		t.Fatalf("设置发布令牌 = %d/%d %s", status, env.Code, env.Msg)
	}

	// 发布路由挂在 user 模块 API 树，测试引擎需手动补挂（含令牌中间件）
	pub := s.engine.Group("/api/publish")
	pub.Use(middleware.RequirePublishToken())
	blogmgrrouter.RegisterPublish(pub)

	// 无令牌 → 401
	if code, _ := publish(s, "", "/api/publish/posts", `{"slug":"pub-a","title":"A","publishStatus":"published"}`); code != http.StatusUnauthorized {
		t.Fatalf("无令牌 = %d，想要 401", code)
	}
	// 错误令牌 → 401
	if code, _ := publish(s, "wrong", "/api/publish/posts", `{"slug":"pub-a","title":"A","publishStatus":"published"}`); code != http.StatusUnauthorized {
		t.Fatalf("错误令牌 = %d，想要 401", code)
	}

	// 正确令牌发布文章：入参 draft 也被强制为 published，未给分类落到默认分类
	status, env = publish(s, "publish-secret", "/api/publish/posts", `{"slug":"pub-a","title":"A","publishStatus":"draft"}`)
	if status != http.StatusCreated || env.Code != 0 {
		t.Fatalf("发布文章 = %d/%d %s", status, env.Code, env.Msg)
	}
	created := dataObject(t, env)
	if created["publishStatus"] != "published" {
		t.Fatalf("发布后状态 = %v，想要 published", created["publishStatus"])
	}
	if created["categoryID"] != "uncategorized" {
		t.Fatalf("未指定分类应落到默认分类，got = %v", created["categoryID"])
	}

	// 重复 slug → 冲突 409/40900
	if status, env = publish(s, "publish-secret", "/api/publish/posts", `{"slug":"pub-a","title":"A2","publishStatus":"published"}`); status != http.StatusConflict || env.Code != 40900 {
		t.Fatalf("重复 slug = %d/%d，想要 409/40900", status, env.Code)
	}

	// 外链利器缺地址 → 40003
	if status, env = publish(s, "publish-secret", "/api/publish/tools", `{"slug":"tool-x","kind":"link","name":"外链","publishStatus":"published"}`); status != http.StatusBadRequest || env.Code != 40003 {
		t.Fatalf("外链缺地址 = %d/%d，想要 400/40003", status, env.Code)
	}

	// 分类/标签不允许通过发布接口创建（路由不存在 → 404）
	req := httptest.NewRequest(http.MethodPost, "/api/publish/categories", strings.NewReader(`{"scope":"post","slug":"x","name":"x"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(middleware.PublishTokenHeader, "publish-secret")
	rec := httptest.NewRecorder()
	s.engine.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("发布分类 = %d，想要 404", rec.Code)
	}
}
