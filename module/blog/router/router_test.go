package router

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"fastduck/treasure-doc/module/blog/data/request"
	"fastduck/treasure-doc/module/blog/data/response"
	"fastduck/treasure-doc/module/blog/internal/service"

	"github.com/gin-gonic/gin"
)

type fakeService struct{}

func (fakeService) Categories(context.Context, string) ([]response.Category, error) {
	return []response.Category{}, nil
}
func (fakeService) PostTags(context.Context) ([]response.Tag, error) {
	return []response.Tag{{ID: "tag-1", Name: "Go"}}, nil
}
func (fakeService) ListPosts(_ context.Context, query request.PostQuery) (response.Page, error) {
	return response.Page{List: []response.PostSummary{}, Pagination: response.Pagination{Page: query.Page, PageSize: query.PageSize, OrderBy: "date_" + query.Sort}}, nil
}
func (fakeService) GetPost(context.Context, string) (response.Post, error) {
	return response.Post{}, service.ErrPostNotFound
}
func (fakeService) DiaryTags(context.Context) ([]response.Tag, error) {
	return []response.Tag{{ID: "tag-1", Name: "Go"}}, nil
}
func (fakeService) ListDiaries(_ context.Context, query request.DiaryQuery) (response.Page, error) {
	return response.Page{List: []response.DiarySummary{}, Pagination: response.Pagination{Page: query.Page, PageSize: query.PageSize, OrderBy: "date_" + query.Sort}}, nil
}
func (fakeService) GetDiary(context.Context, string) (response.Diary, error) {
	return response.Diary{}, service.ErrDiaryNotFound
}
func (fakeService) ListPortfolio(context.Context, request.PortfolioQuery) ([]response.PortfolioSummary, error) {
	return []response.PortfolioSummary{}, nil
}
func (fakeService) GetPortfolio(context.Context, string) (response.PortfolioItem, error) {
	return response.PortfolioItem{}, service.ErrPortfolioNotFound
}
func (fakeService) ListTools(context.Context) ([]response.Tool, error) { return []response.Tool{}, nil }
func (fakeService) GetTool(context.Context, string) (response.Tool, error) {
	return response.Tool{}, service.ErrToolNotFound
}
func (fakeService) ListBookmarks(context.Context, request.BookmarkQuery) ([]response.Bookmark, error) {
	return []response.Bookmark{}, nil
}
func (fakeService) Profile(context.Context) (response.Profile, error) {
	return response.Profile{Links: []response.ProfileLink{}, Skills: []response.ProfileSkill{}}, nil
}
func (fakeService) Site(context.Context) (response.Site, error) {
	return response.Site{TechStack: []string{}, Modules: []response.SiteModule{}, Milestones: []response.SiteMilestone{}}, nil
}
func (fakeService) Stats(context.Context) (response.Stats, error) { return response.Stats{}, nil }

func newEngine() *gin.Engine {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	RegisterService(engine.Group("/api"), fakeService{})
	return engine
}

func TestPublicRoutes(t *testing.T) {
	routes := []string{
		"/api/blog/categories", "/api/blog/tags", "/api/blog/posts",
		"/api/blog/diary/tags", "/api/blog/diaries", "/api/blog/portfolio/categories",
		"/api/blog/portfolio/items", "/api/blog/tools", "/api/blog/bookmark/categories",
		"/api/blog/bookmarks", "/api/blog/profile", "/api/blog/site", "/api/blog/stats",
	}
	engine := newEngine()
	for _, path := range routes {
		t.Run(path, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			engine.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))
			if recorder.Code != http.StatusOK {
				t.Fatalf("GET %s status = %d, body = %s", path, recorder.Code, recorder.Body.String())
			}
			if !strings.Contains(recorder.Body.String(), `"code":0,"msg":""`) {
				t.Fatalf("GET %s returned unexpected envelope: %s", path, recorder.Body.String())
			}
		})
	}
}

func TestTagRoutesIncludeIDAndName(t *testing.T) {
	for _, path := range []string{"/api/blog/tags", "/api/blog/diary/tags"} {
		recorder := httptest.NewRecorder()
		newEngine().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))
		body := recorder.Body.String()
		if recorder.Code != http.StatusOK || !strings.Contains(body, `"id":"tag-1"`) || !strings.Contains(body, `"name":"Go"`) {
			t.Fatalf("GET %s returned unexpected tag response: %d %s", path, recorder.Code, body)
		}
	}
}

func TestPostQueryErrors(t *testing.T) {
	tests := []struct {
		path string
		code int
	}{
		{path: "/api/blog/posts?page=bad", code: response.CodeInvalidQuery},
		{path: "/api/blog/posts?pageSize=101", code: response.CodeInvalidQuery},
		{path: "/api/blog/posts?sort=newest", code: response.CodeUnsupportedSort},
	}
	engine := newEngine()
	for _, test := range tests {
		recorder := httptest.NewRecorder()
		engine.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, test.path, nil))
		if recorder.Code != http.StatusBadRequest {
			t.Fatalf("GET %s status = %d", test.path, recorder.Code)
		}
		if !strings.Contains(recorder.Body.String(), `"code":`+stringCode(test.code)) || !strings.Contains(recorder.Body.String(), `"data":null`) {
			t.Fatalf("GET %s returned unexpected body: %s", test.path, recorder.Body.String())
		}
	}
}

func TestDetailNotFoundCodes(t *testing.T) {
	tests := []struct {
		path string
		code int
	}{
		{path: "/api/blog/posts/missing", code: response.CodePostNotFound},
		{path: "/api/blog/diaries/missing", code: response.CodeDiaryNotFound},
		{path: "/api/blog/portfolio/items/missing", code: response.CodePortfolioMissing},
		{path: "/api/blog/tools/missing", code: response.CodeToolNotFound},
	}
	engine := newEngine()
	for _, test := range tests {
		recorder := httptest.NewRecorder()
		engine.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, test.path, nil))
		if recorder.Code != http.StatusNotFound || !strings.Contains(recorder.Body.String(), `"code":`+stringCode(test.code)) {
			t.Fatalf("GET %s status/body = %d %s", test.path, recorder.Code, recorder.Body.String())
		}
	}
}

func TestNoPublicWriteRoutes(t *testing.T) {
	recorder := httptest.NewRecorder()
	newEngine().ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/api/blog/posts", nil))
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("POST public route status = %d, want 404", recorder.Code)
	}
}

func TestLegacyPublicPrefixIsNotRegistered(t *testing.T) {
	for _, path := range []string{"/api/public/blog/tags", "/blog/tags"} {
		recorder := httptest.NewRecorder()
		newEngine().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))
		if recorder.Code != http.StatusNotFound {
			t.Fatalf("legacy route %s status = %d, want 404", path, recorder.Code)
		}
	}
}

func stringCode(code int) string {
	switch code {
	case response.CodeInvalidQuery:
		return "40001"
	case response.CodeUnsupportedSort:
		return "40002"
	case response.CodePostNotFound:
		return "40401"
	case response.CodeDiaryNotFound:
		return "40402"
	case response.CodePortfolioMissing:
		return "40403"
	case response.CodeToolNotFound:
		return "40404"
	default:
		panic(errors.New("unexpected code"))
	}
}
