package service

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	blogmodel "fastduck/treasure-doc/module/blog/data/model"
	blogresponse "fastduck/treasure-doc/module/blog/data/response"
	"fastduck/treasure-doc/module/blog_mgr/data/request"
)

func TestPublishedTimesDefaults(t *testing.T) {
	before := time.Now()
	on, at, err := publishedTimes(blogmodel.StatusPublished, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	if at.Before(before) || on.Format("2006-01-02") != at.Format("2006-01-02") {
		t.Fatalf("unexpected times: %v %v", on, at)
	}
}

func TestBuildModelDeduplicatesTags(t *testing.T) {
	item, tags, relation, err := buildModel("posts", request.Post{Slug: "post", Title: "Post", PublishStatus: blogmodel.StatusDraft, TagIDs: []string{"a", "a", "b"}})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := item.(*blogmodel.Post); !ok {
		t.Fatalf("unexpected model %T", item)
	}
	if relation != "td_blog_post_tag" || len(tags) != 2 {
		t.Fatalf("unexpected relation/tags: %s %#v", relation, tags)
	}
}

func TestTaggedManagementResponsesUseTagIDs(t *testing.T) {
	values := []interface{}{
		PostWithTags{Post: blogmodel.Post{BaseModel: blogmodel.BaseModel{ID: "post-1"}}, TagIDs: []string{}},
		DiaryWithTags{Diary: blogmodel.Diary{BaseModel: blogmodel.BaseModel{ID: "diary-1"}}, TagIDs: []string{"tag-1"}},
		BookmarkWithTags{Bookmark: blogmodel.Bookmark{BaseModel: blogmodel.BaseModel{ID: "bookmark-1"}}, TagIDs: []string{}},
	}
	for _, value := range values {
		data, err := json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(data), `"tagIds"`) {
			t.Fatalf("response %T does not contain tagIds: %s", value, data)
		}
	}
}

func TestNormalizeSiteModules(t *testing.T) {
	modules, err := normalizeSiteModules([]blogresponse.SiteModule{
		{ID: "diary", Icon: "D", Name: "日记", Desc: "记录", Path: "/Blog/Diary", Visible: false},
		{ID: "blog", Icon: "B", Name: "文章", Desc: "长文", Path: "/Blog", Visible: true},
	}, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(modules) != 6 {
		t.Fatalf("expected six modules, got %d", len(modules))
	}
	if modules[0].ID != "blog" || modules[1].ID != "diary" {
		t.Fatalf("modules are not in fixed order: %#v", modules)
	}
	if modules[1].Visible {
		t.Fatal("explicit visible=false was not preserved")
	}
	if !modules[2].Visible {
		t.Fatal("missing module should default to visible=true")
	}
}

func TestNormalizeSiteModulesRejectsInvalidPutPayload(t *testing.T) {
	modules := defaultSiteModules()
	modules[0].Path = "/Blog/Changed"
	if _, err := normalizeSiteModules(modules, true); !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected invalid path error, got %v", err)
	}

	if _, err := normalizeSiteModules(modules[:5], true); !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected missing module error, got %v", err)
	}
}
