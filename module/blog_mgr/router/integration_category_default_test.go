package router_test

import (
	"fmt"
	"net/http"
	"testing"

	blogmodel "fastduck/treasure-doc/module/blog/data/model"
)

// getCategoryVersion 读取某资源详情的 version，供乐观并发更新使用。
func getCategoryVersion(t *testing.T, s *testServer, path string) float64 {
	t.Helper()
	_, env := s.do(http.MethodGet, path, "")
	return dataObject(t, env)["version"].(float64)
}

// TestDefaultCategoryAssignment M-CAT-01：未指定分类的内容自动落到默认分类，默认分类不可删除。
func TestDefaultCategoryAssignment(t *testing.T) {
	s := newTestServer(t)

	// 未指定分类 → 自动落到 post 作用域的默认分类并创建该分类
	postID := s.createPost("no-cat-post", "未分类文章", `,"publishStatus":"published"`)
	_, env := s.do(http.MethodGet, "/api/blog-mgr/posts/"+postID, "")
	if got, _ := dataObject(t, env)["categoryID"].(string); got != "uncategorized" {
		t.Fatalf("未指定分类的文章应为默认分类，got = %q", got)
	}
	var def blogmodel.Category
	if err := s.db.Where("scope = ? AND slug = ?", blogmodel.CategoryPost, "uncategorized").First(&def).Error; err != nil {
		t.Fatalf("默认分类未被自动创建：%v", err)
	}

	// 默认分类不可删除（无论是否被引用）
	status, delEnv := s.do(http.MethodDelete, "/api/blog-mgr/categories/"+def.ID, "")
	if status != http.StatusConflict || delEnv.Code != 40900 {
		t.Fatalf("删除默认分类 = %d/%d，想要 409/40900（%s）", status, delEnv.Code, delEnv.Msg)
	}

	// 显式指定分类的内容不受影响
	s.create("categories", `{"scope":"post","slug":"tech","name":"技术"}`)
	explicitID := s.createPost("cat-post", "有分类文章", `,"categoryId":"tech"`)
	_, env = s.do(http.MethodGet, "/api/blog-mgr/posts/"+explicitID, "")
	if got, _ := dataObject(t, env)["categoryID"].(string); got != "tech" {
		t.Fatalf("显式分类应保留 tech，got = %q", got)
	}

	// 更新时把分类清空 → 重新落到默认分类
	version := getCategoryVersion(t, s, "/api/blog-mgr/posts/"+explicitID)
	body := fmt.Sprintf(`{"slug":"cat-post","title":"有分类文章","publishStatus":"published","categoryId":"","version":%d}`, int(version))
	if status, upEnv := s.do(http.MethodPatch, "/api/blog-mgr/posts/"+explicitID, body); status != http.StatusOK || upEnv.Code != 0 {
		t.Fatalf("更新清空分类 = %d/%d %s", status, upEnv.Code, upEnv.Msg)
	}
	_, env = s.do(http.MethodGet, "/api/blog-mgr/posts/"+explicitID, "")
	if got, _ := dataObject(t, env)["categoryID"].(string); got != "uncategorized" {
		t.Fatalf("清空分类后应落到默认分类，got = %q", got)
	}
}
