package request

import (
	"errors"
	"testing"
)

func TestPageQueryNormalize(t *testing.T) {
	tests := []struct {
		name    string
		query   PageQuery
		want    PageQuery
		wantErr error
	}{
		{name: "defaults", query: PageQuery{}, want: PageQuery{Page: 1, PageSize: 10, Sort: "desc"}},
		{name: "explicit", query: PageQuery{Page: 2, PageSize: 20, Sort: "asc"}, want: PageQuery{Page: 2, PageSize: 20, Sort: "asc"}},
		{name: "invalid page", query: PageQuery{Page: -1}, wantErr: ErrInvalidQuery},
		{name: "invalid page size", query: PageQuery{PageSize: 101}, wantErr: ErrInvalidQuery},
		{name: "invalid sort", query: PageQuery{Sort: "newest"}, wantErr: ErrInvalidSort},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.query.Normalize()
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("Normalize() error = %v, want %v", err, test.wantErr)
			}
			if test.wantErr == nil && test.query != test.want {
				t.Fatalf("Normalize() query = %#v, want %#v", test.query, test.want)
			}
		})
	}
}

func TestPostQueryNormalize(t *testing.T) {
	query := PostQuery{CategoryID: " tech ", Tag: " Vue ", Keyword: " performance "}
	if err := query.Normalize(); err != nil {
		t.Fatalf("Normalize() error = %v", err)
	}
	if query.CategoryID != "tech" || query.Tag != "Vue" || query.Keyword != "performance" {
		t.Fatalf("Normalize() did not trim filters: %#v", query)
	}
}

func TestEscapeLike(t *testing.T) {
	if got, want := EscapeLike(`50%_off\today`), `50\%\_off\\today`; got != want {
		t.Fatalf("EscapeLike() = %q, want %q", got, want)
	}
}

func TestValidPublicID(t *testing.T) {
	if ValidPublicID("") {
		t.Fatal("empty ID should be invalid")
	}
	if !ValidPublicID("hello-world") {
		t.Fatal("normal ID should be valid")
	}
	longID := make([]byte, 129)
	if ValidPublicID(string(longID)) {
		t.Fatal("ID longer than 128 bytes should be invalid")
	}
}
