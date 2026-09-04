package request

import "testing"

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
}
