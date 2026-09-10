package response

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestSiteModuleUnmarshalDefaultsMissingVisibleToTrue(t *testing.T) {
	var module SiteModule
	if err := json.Unmarshal([]byte(`{"id":"blog","path":"/Blog","name":"文章"}`), &module); err != nil {
		t.Fatal(err)
	}
	if !module.Visible {
		t.Fatal("legacy module without visible should default to true")
	}
}

func TestSiteModuleUnmarshalKeepsExplicitVisible(t *testing.T) {
	var module SiteModule
	if err := json.Unmarshal([]byte(`{"id":"diary","path":"/Blog/Diary","name":"日记","visible":false}`), &module); err != nil {
		t.Fatal(err)
	}
	if module.Visible {
		t.Fatal("explicit visible=false should be preserved")
	}
}

func TestSiteModuleUnmarshalRejectsNonBooleanVisible(t *testing.T) {
	var module SiteModule
	if err := json.Unmarshal([]byte(`{"id":"blog","visible":"yes"}`), &module); err == nil {
		t.Fatal("expected error for non-boolean visible")
	}
}

func TestNormalizeSiteModulesRejectsOverlongNameForStrictPayload(t *testing.T) {
	modules := DefaultSiteModules()
	modules[0].Name = strings.Repeat("菜", MaxSiteModuleNameLength+1)
	if _, err := NormalizeSiteModules(modules, true); !errors.Is(err, ErrInvalidSiteModule) {
		t.Fatalf("expected overlong name rejection, got %v", err)
	}
}

func TestNormalizeSiteModulesRejectsOverlongSubtitleForStrictPayload(t *testing.T) {
	modules := DefaultSiteModules()
	modules[0].Desc = strings.Repeat("副", MaxSiteModuleSubtitleLength+1)
	if _, err := NormalizeSiteModules(modules, true); !errors.Is(err, ErrInvalidSiteModule) {
		t.Fatalf("expected overlong subtitle rejection, got %v", err)
	}
}

func TestNormalizeSiteModulesDefaultsLegacyMarker(t *testing.T) {
	modules, err := NormalizeSiteModules([]SiteModule{{ID: "blog", Name: "文章", Path: "/Blog"}}, false)
	if err != nil {
		t.Fatal(err)
	}
	if modules[1].Marker != "BLOG" {
		t.Fatalf("expected default marker BLOG, got %q", modules[1].Marker)
	}
}

func TestNormalizeSiteModulesDefaultsLegacyTitle(t *testing.T) {
	modules, err := NormalizeSiteModules([]SiteModule{{ID: "diary", Name: "日记", Path: "/Blog/Diary"}}, false)
	if err != nil {
		t.Fatal(err)
	}
	if modules[2].Title != "日记" {
		t.Fatalf("expected default title 日记, got %q", modules[2].Title)
	}
}

func TestNormalizeSiteModulesAllowsEmptyTitleForStrictPayload(t *testing.T) {
	modules := DefaultSiteModules()
	modules[0].Title = ""
	normalized, err := NormalizeSiteModules(modules, true)
	if err != nil {
		t.Fatalf("expected empty title to be allowed, got %v", err)
	}
	if normalized[0].Title != "" {
		t.Fatalf("expected empty title to remain empty, got %q", normalized[0].Title)
	}
}

func TestNormalizeSiteModulesRejectsEmptyMarkerForStrictPayload(t *testing.T) {
	modules := DefaultSiteModules()
	modules[0].Marker = ""
	if _, err := NormalizeSiteModules(modules, true); !errors.Is(err, ErrInvalidSiteModule) {
		t.Fatalf("expected empty marker rejection, got %v", err)
	}
}

func TestNormalizeSiteModulesKeepsOrderAndDropsUnknown(t *testing.T) {
	modules, err := NormalizeSiteModules([]SiteModule{
		{ID: "about", Name: "关于", Path: "/Blog/About", Visible: false},
		{ID: "legacy", Name: "旧模块", Path: "/Blog/Legacy"},
	}, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(modules) != 7 || modules[0].ID != "home" || modules[6].ID != "about" {
		t.Fatalf("unexpected module set: %#v", modules)
	}
	if modules[6].Visible {
		t.Fatal("about should stay hidden")
	}
	if !modules[0].Visible {
		t.Fatal("missing module should default to visible=true")
	}
}
