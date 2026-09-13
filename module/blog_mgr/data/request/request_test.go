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
	if err := ValidateTool(Tool{Slug: "ext", Kind: "link", Name: "外链", URL: "javascript:alert(1)", PublishStatus: "draft"}); !errors.Is(err, ErrToolURLRequired) {
		t.Fatalf("link tool with unsafe scheme = %v, want ErrToolURLRequired", err)
	}
	// 链接类地址不限协议：http 与「没写协议」都应放行（后者由 NormalizeLinkURL 补 https）
	if ValidateTool(Tool{Slug: "ext", Kind: "link", Name: "外链", URL: "http://example.com", PublishStatus: "draft"}) != nil {
		t.Fatal("http link tool rejected")
	}
	if ValidateTool(Tool{Slug: "ext", Kind: "link", Name: "外链", URL: "example.com", PublishStatus: "draft"}) != nil {
		t.Fatal("scheme-less link tool rejected")
	}
}

func TestNormalizeLinkURL(t *testing.T) {
	tests := []struct {
		input string
		want  string
		ok    bool
	}{
		{input: "https://example.com", want: "https://example.com", ok: true},
		{input: "http://example.com/a?b=1", want: "http://example.com/a?b=1", ok: true},
		{input: "ftp://example.com/file", want: "ftp://example.com/file", ok: true},
		{input: "mailto:me@example.com", want: "mailto:me@example.com", ok: true},
		// 没写协议：按 https 补全
		{input: "example.com", want: "https://example.com", ok: true},
		{input: "  example.com/path  ", want: "https://example.com/path", ok: true},
		// 站内路径原样保留
		{input: "/files/blog/a.png", want: "/files/blog/a.png", ok: true},
		// 可执行脚本的协议与空值必须拒绝
		{input: "javascript:alert(1)", ok: false},
		{input: "data:text/html;base64,PHNjcmlwdD4=", ok: false},
		{input: "vbscript:msgbox(1)", ok: false},
		{input: "", ok: false},
		{input: "   ", ok: false},
	}
	for _, test := range tests {
		got, ok := NormalizeLinkURL(test.input)
		if ok != test.ok || (ok && got != test.want) {
			t.Fatalf("NormalizeLinkURL(%q) = (%q, %v), want (%q, %v)", test.input, got, ok, test.want, test.ok)
		}
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
