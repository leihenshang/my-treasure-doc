package seed

import (
	"testing"
	"time"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		options Options
		wantErr bool
	}{
		{"disabled", Options{}, false}, {"local", Options{Enabled: true, Host: "127.0.0.1"}, false},
		{"remote blocked", Options{Enabled: true, Host: "db.example.com"}, true},
		{"remote allowed", Options{Enabled: true, Host: "db.example.com", AllowRemote: true}, false},
		{"release blocked", Options{Enabled: true, Host: "localhost", Release: true}, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if (Validate(test.options) != nil) != test.wantErr {
				t.Fatal("Validate() mismatch")
			}
		})
	}
}

func TestContentSeedDistribution(t *testing.T) {
	for name, items := range map[string][]contentSeed{"posts": postSeeds(), "diaries": diarySeeds()} {
		if len(items) != 20 {
			t.Fatalf("%s count = %d", name, len(items))
		}
		visible := 0
		keys := map[string]struct{}{}
		for _, item := range items {
			if _, ok := keys[item.Key]; ok {
				t.Fatalf("duplicate %s key %s", name, item.Key)
			}
			keys[item.Key] = struct{}{}
			if item.Status == "published" && !item.PublishedAt.After(time.Now()) && !item.Deleted {
				visible++
			}
		}
		if visible != 12 {
			t.Fatalf("%s visible count = %d", name, visible)
		}
	}
}
