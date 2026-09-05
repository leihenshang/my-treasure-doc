package response

import (
	"encoding/json"
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

func TestNormalizeSiteModulesKeepsOrderAndDropsUnknown(t *testing.T) {
	modules, err := NormalizeSiteModules([]SiteModule{
		{ID: "about", Name: "关于", Path: "/Blog/About", Visible: false},
		{ID: "legacy", Name: "旧模块", Path: "/Blog/Legacy"},
	}, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(modules) != 6 || modules[0].ID != "blog" || modules[5].ID != "about" {
		t.Fatalf("unexpected module set: %#v", modules)
	}
	if modules[5].Visible {
		t.Fatal("about should stay hidden")
	}
	if !modules[0].Visible {
		t.Fatal("missing module should default to visible=true")
	}
}
