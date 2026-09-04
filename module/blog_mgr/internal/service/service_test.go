package service

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	blogmodel "fastduck/treasure-doc/module/blog/data/model"
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
