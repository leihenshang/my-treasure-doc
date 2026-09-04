package router

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	blogresponse "fastduck/treasure-doc/module/blog/data/response"
	"fastduck/treasure-doc/module/blog_mgr/data/request"
	"fastduck/treasure-doc/module/blog_mgr/internal/service"

	"github.com/gin-gonic/gin"
)

type fakeManager struct{}

func (fakeManager) List(context.Context, string, request.List) (service.Page, error) {
	return service.Page{List: []interface{}{}}, nil
}
func (fakeManager) Get(context.Context, string, string) (interface{}, error) {
	return map[string]string{"id": "1"}, nil
}
func (fakeManager) Create(context.Context, string, interface{}) (interface{}, error) {
	return map[string]string{"id": "1"}, nil
}
func (fakeManager) Update(context.Context, string, string, interface{}) (interface{}, error) {
	return map[string]string{"id": "1"}, nil
}
func (fakeManager) Delete(context.Context, string, string) error  { return nil }
func (fakeManager) Restore(context.Context, string, string) error { return nil }
func (fakeManager) GetSetting(context.Context, string) (interface{}, error) {
	return map[string]string{}, nil
}
func (fakeManager) PutSetting(context.Context, string, interface{}) (interface{}, error) {
	return map[string]string{}, nil
}

func TestManagementRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	RegisterService(engine.Group("/api/blog-mgr"), fakeManager{})
	resources := []string{"categories", "tags", "posts", "diaries", "portfolio-items", "tools", "bookmarks"}
	bodies := map[string]string{"categories": `{"scope":"post","slug":"tech","name":"Tech"}`, "tags": `{"name":"Go"}`, "posts": `{"slug":"p","title":"P","publishStatus":"draft"}`, "diaries": `{"publicId":"d","title":"D","publishStatus":"draft"}`, "portfolio-items": `{"slug":"w","title":"W","publishStatus":"draft"}`, "tools": `{"slug":"t","kind":"own","name":"T","developmentStatus":"开发中","publishStatus":"draft"}`, "bookmarks": `{"publicId":"b","title":"B","url":"https://example.com","publishStatus":"draft"}`}
	for _, resource := range resources {
		tests := []struct {
			method, path, body string
			status             int
		}{{http.MethodGet, "/api/blog-mgr/" + resource, "", 200}, {http.MethodPost, "/api/blog-mgr/" + resource, bodies[resource], 201}, {http.MethodGet, "/api/blog-mgr/" + resource + "/1", "", 200}, {http.MethodPatch, "/api/blog-mgr/" + resource + "/1", bodies[resource], 200}, {http.MethodDelete, "/api/blog-mgr/" + resource + "/1", "", 200}, {http.MethodPost, "/api/blog-mgr/" + resource + "/1/restore", "", 200}}
		for _, test := range tests {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(test.method, test.path, strings.NewReader(test.body))
			req.Header.Set("Content-Type", "application/json")
			engine.ServeHTTP(rec, req)
			if rec.Code != test.status {
				t.Fatalf("%s %s = %d: %s", test.method, test.path, rec.Code, rec.Body.String())
			}
		}
	}
	for _, setting := range []string{"profile", "site"} {
		rec := httptest.NewRecorder()
		engine.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/blog-mgr/"+setting, nil))
		if rec.Code != 200 {
			t.Fatalf("GET %s = %d", setting, rec.Code)
		}
		var body string
		if setting == "profile" {
			body = `{"name":"N","links":[],"skills":[]}`
		} else {
			body = `{"name":"S","techStack":[],"modules":[],"milestones":[]}`
		}
		rec = httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/api/blog-mgr/"+setting, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		engine.ServeHTTP(rec, req)
		if rec.Code != 200 {
			t.Fatalf("PUT %s = %d: %s", setting, rec.Code, rec.Body.String())
		}
	}
	_ = blogresponse.Stats{}
}
