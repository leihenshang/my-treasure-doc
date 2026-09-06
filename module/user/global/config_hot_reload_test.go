package global

import (
	"reflect"
	"testing"

	"fastduck/treasure-doc/module/user/config"
)

func TestEffectiveHotReloadConfig(t *testing.T) {
	current := &config.Config{
		App:      config.App{Port: 2021, RunMode: "release", RegisterEnabled: false},
		Database: config.Database{Driver: "postgres", Dsn: "host=db-old dbname=treasure_doc"},
		Redis:    config.Redis{Enable: true, Host: "redis-old"},
		Log:      config.Log{Level: "info"},
		Debug:    config.Debug{EnableMockLogin: false},
	}
	candidate := *current
	candidate.App.RegisterEnabled = true
	candidate.App.Port = 3030
	candidate.Database.Dsn = "host=db-new dbname=treasure_doc"
	candidate.Redis.Host = "redis-new"
	candidate.Log.Level = "debug"
	candidate.Debug.EnableMockLogin = true
	candidate.BlogSeed.Enabled = true

	effective, restartRequired := effectiveHotReloadConfig(current, &candidate)
	if !effective.App.RegisterEnabled {
		t.Fatal("registerEnabled should be hot reloaded")
	}
	if effective.App.Port != current.App.Port || effective.Database != current.Database || effective.Redis != current.Redis || effective.Log != current.Log || effective.Debug != current.Debug {
		t.Fatalf("restart-required config changed: %#v", effective)
	}
	wantSections := []string{"app", "database", "redis", "log", "debug", "blogSeed"}
	if !reflect.DeepEqual(restartRequired, wantSections) {
		t.Fatalf("restart sections = %#v, want %#v", restartRequired, wantSections)
	}
}

func TestEffectiveHotReloadConfigBusinessOnly(t *testing.T) {
	current := &config.Config{App: config.App{RegisterEnabled: false}}
	candidate := *current
	candidate.App.RegisterEnabled = true
	effective, restartRequired := effectiveHotReloadConfig(current, &candidate)
	if !effective.App.RegisterEnabled || len(restartRequired) != 0 {
		t.Fatalf("effective = %#v, restart = %#v", effective, restartRequired)
	}
}

func TestEffectiveHotReloadConfigCaptcha(t *testing.T) {
	current := &config.Config{Captcha: config.Captcha{Enabled: true}}
	candidate := *current
	candidate.Captcha.Enabled = false
	effective, restartRequired := effectiveHotReloadConfig(current, &candidate)
	if effective.Captcha.Enabled {
		t.Fatal("captcha should be hot reloaded")
	}
	if len(restartRequired) != 0 {
		t.Fatalf("restart sections = %#v, want none", restartRequired)
	}
}

func TestValidateStartupConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  config.Config
		wantErr bool
	}{
		{name: "dev mock", config: config.Config{App: config.App{RunMode: config.GinModeDev}, Debug: config.Debug{EnableMockLogin: true}}},
		{name: "release mock", config: config.Config{App: config.App{RunMode: config.GinModeRelease}, Debug: config.Debug{EnableMockLogin: true}}, wantErr: true},
		{name: "dev seed", config: config.Config{App: config.App{RunMode: config.GinModeDev}, BlogSeed: config.BlogSeed{Enabled: true}}},
		{name: "release seed", config: config.Config{App: config.App{RunMode: config.GinModeRelease}, BlogSeed: config.BlogSeed{Enabled: true}}, wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateStartupConfig(&test.config)
			if (err != nil) != test.wantErr {
				t.Fatalf("validateStartupConfig() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}
