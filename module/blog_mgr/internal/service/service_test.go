package service

import (
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
