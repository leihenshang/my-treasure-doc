package request

import (
	"errors"
	"strconv"
	"testing"
)

func TestListNormalize(t *testing.T) {
	query := List{}
	if err := query.Normalize(); err != nil {
		t.Fatal(err)
	}
	if query.Page != 1 || query.PageSize != 20 || query.Sort != "desc" || query.Deleted != "exclude" {
		t.Fatalf("unexpected defaults: %#v", query)
	}
	query = List{PageSize: 101}
	if query.Normalize() == nil {
		t.Fatal("pageSize > 100 should fail")
	}
}

func TestValidation(t *testing.T) {
	if !ValidDate("2026-09-04") || ValidDate("2026/09/04") {
		t.Fatal("date validation mismatch")
	}
	if !ValidURL("https://example.com", false) || ValidURL("http://example.com", false) {
		t.Fatal("URL validation mismatch")
	}
	if ValidateTool(Tool{Slug: "mdn", Kind: "link", Name: "MDN", URL: "https://developer.mozilla.org", PublishStatus: "published"}) != nil {
		t.Fatal("valid link tool rejected")
	}
	if ValidateTool(Tool{Slug: "mdn", Kind: "link", Name: "MDN", URL: "javascript:alert(1)", PublishStatus: "published"}) == nil {
		t.Fatal("unsafe URL accepted")
	}
	// 缺失字段必须给出可定位的细分错误，而不是笼统的 ErrInvalid
	if err := ValidateTool(Tool{Slug: "own", Kind: "own", Name: "自研", PublishStatus: "draft"}); !errors.Is(err, ErrToolStatusRequired) {
		t.Fatalf("own tool without developmentStatus = %v, want ErrToolStatusRequired", err)
	}
	if err := ValidateTool(Tool{Slug: "ext", Kind: "link", Name: "外链", URL: "http://example.com", PublishStatus: "draft"}); !errors.Is(err, ErrToolURLRequired) {
		t.Fatalf("link tool without https url = %v, want ErrToolURLRequired", err)
	}
}

func TestBatchIDsNormalize(t *testing.T) {
	if (&BatchIDs{}).Normalize() == nil {
		t.Fatal("empty list should fail")
	}
	if (&BatchIDs{IDs: []string{"ok", "  "}}).Normalize() == nil {
		t.Fatal("blank id should fail")
	}
	payload := BatchIDs{IDs: []string{" b ", "a", "a"}}
	if err := payload.Normalize(); err != nil {
		t.Fatal(err)
	}
	if len(payload.IDs) != 2 || payload.IDs[0] != "b" || payload.IDs[1] != "a" {
		t.Fatalf("unexpected normalized ids: %#v", payload.IDs)
	}
	tooMany := make([]string, MaxBatchIDs+1)
	for i := range tooMany {
		tooMany[i] = strconv.Itoa(i)
	}
	if (&BatchIDs{IDs: tooMany}).Normalize() == nil {
		t.Fatal("over limit should fail")
	}
}
