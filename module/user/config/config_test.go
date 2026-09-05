package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestUnmarshalCaptchaDefaultsToEnabled(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(path, []byte("[app]\nport = 2021\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		t.Fatal(err)
	}
	applyConfigDefaults(v)
	got := &Config{}
	if err := v.Unmarshal(got); err != nil {
		t.Fatal(err)
	}
	if !got.Captcha.Enabled {
		t.Fatal("captcha should be enabled when the config section is missing")
	}

	v.Set("captcha.enabled", false)
	got = &Config{}
	if err := v.Unmarshal(got); err != nil {
		t.Fatal(err)
	}
	if got.Captcha.Enabled {
		t.Fatal("captcha should be disabled when explicitly configured")
	}
}

func TestPublishConfigCopiesSnapshot(t *testing.T) {
	original := &Config{App: App{RegisterEnabled: false, Port: 2021}}
	PublishConfig(original)
	original.App.RegisterEnabled = true
	original.App.Port = 3030

	got := GetConfig()
	if got.App.RegisterEnabled || got.App.Port != 2021 {
		t.Fatalf("published config changed through source pointer: %#v", got.App)
	}
}
