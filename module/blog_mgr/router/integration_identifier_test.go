package router_test

import (
	"fmt"
	"net/http"
	"testing"
)

// TestIdentifierAutoGenerate M-ID-01：posts 的 slug / diaries 的 publicId 省略时自动生成；
// portfolio-items / tools 的 slug 同样自动生成；更新省略时沿用原值，不改变公开 URL。
func TestIdentifierAutoGenerate(t *testing.T) {
	s := newTestServer(t)

	// 文章省略 slug → 按标题自动生成
	postID := s.createPost("", "Hello World 文章", `,"publishStatus":"published"`)
	_, env := s.do(http.MethodGet, "/api/blog-mgr/posts/"+postID, "")
	if got, _ := dataObject(t, env)["slug"].(string); got != "hello-world-文章" {
		t.Fatalf("文章自动生成 slug = %q", got)
	}

	// 标题相同 → 自动追加序号避免冲突
	postID2 := s.createPost("", "Hello World 文章", "")
	_, env = s.do(http.MethodGet, "/api/blog-mgr/posts/"+postID2, "")
	if got, _ := dataObject(t, env)["slug"].(string); got != "hello-world-文章-2" {
		t.Fatalf("冲突文章 slug = %q", got)
	}

	// 日记省略 publicId → 按标题自动生成
	diaryID := s.create("diaries", `{"title":"我的日记","publishStatus":"published"}`)
	_, env = s.do(http.MethodGet, "/api/blog-mgr/diaries/"+diaryID, "")
	if got, _ := dataObject(t, env)["publicID"].(string); got != "我的日记" {
		t.Fatalf("日记自动生成 publicId = %q", got)
	}

	// 作品省略 slug → 按标题自动生成
	workID := s.create("portfolio-items", `{"title":"我的作品","publishStatus":"published"}`)
	_, env = s.do(http.MethodGet, "/api/blog-mgr/portfolio-items/"+workID, "")
	if got, _ := dataObject(t, env)["slug"].(string); got != "我的作品" {
		t.Fatalf("作品自动生成 slug = %q", got)
	}

	// 工具省略 slug → 按名称自动生成
	toolID := s.create("tools", `{"kind":"own","name":"我的工具","developmentStatus":"开发中","publishStatus":"published"}`)
	_, env = s.do(http.MethodGet, "/api/blog-mgr/tools/"+toolID, "")
	if got, _ := dataObject(t, env)["slug"].(string); got != "我的工具" {
		t.Fatalf("工具自动生成 slug = %q", got)
	}

	// 更新时省略 slug → 沿用原值
	_, env = s.do(http.MethodGet, "/api/blog-mgr/posts/"+postID, "")
	detail := dataObject(t, env)
	version, _ := detail["version"].(float64)
	oldSlug, _ := detail["slug"].(string)
	body := fmt.Sprintf(`{"title":"重新命名","publishStatus":"published","version":%d}`, int(version))
	if status, up := s.do(http.MethodPatch, "/api/blog-mgr/posts/"+postID, body); status != http.StatusOK || up.Code != 0 {
		t.Fatalf("更新省略 slug = %d/%d %s", status, up.Code, up.Msg)
	}
	_, env = s.do(http.MethodGet, "/api/blog-mgr/posts/"+postID, "")
	if got, _ := dataObject(t, env)["slug"].(string); got != oldSlug {
		t.Fatalf("更新后 slug 变化 = %q，想要沿用 %q", got, oldSlug)
	}
}
