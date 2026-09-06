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
		{"disabled", Options{}, false},
		{"sqlite local", Options{Enabled: true, Driver: "sqlite"}, false},
		{"pg local", Options{Enabled: true, Driver: "postgres", Dsn: "host=127.0.0.1 user=postgres dbname=treasure_doc"}, false},
		{"pg remote blocked", Options{Enabled: true, Driver: "postgres", Dsn: "host=db.example.com dbname=treasure_doc"}, true},
		{"pg remote allowed", Options{Enabled: true, Driver: "postgres", Dsn: "host=db.example.com dbname=treasure_doc", AllowRemote: true}, false},
		{"release blocked", Options{Enabled: true, Driver: "sqlite", Release: true}, true},
		{"unknown driver blocked", Options{Enabled: true, Driver: "mysql", Dsn: "host=db.example.com"}, true},
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
