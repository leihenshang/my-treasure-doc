package global

import (
	"reflect"
	"testing"

	"fastduck/treasure-doc/module/user/config"
)

func TestEffectiveHotReloadConfig(t *testing.T) {
	current := &config.Config{
		App:      config.App{Port: 2021, RunMode: "release", RegisterEnabled: false},
		Database: config.Database{Host: "db-old", Port: 5432},
		Redis:    config.Redis{Enable: true, Host: "redis-old"},
		Log:      config.Log{Level: "info"},
		Debug:    config.Debug{EnableMockLogin: false},
	}
	candidate := *current
	candidate.App.RegisterEnabled = true
	candidate.App.Port = 3030
	candidate.Database.Host = "db-new"
	candidate.Redis.Host = "redis-new"
	candidate.Log.Level = "debug"
	candidate.Debug.EnableMockLogin = true

	effective, restartRequired := effectiveHotReloadConfig(current, &candidate)
	if !effective.App.RegisterEnabled {
		t.Fatal("registerEnabled should be hot reloaded")
	}
	if effective.App.Port != current.App.Port || effective.Database != current.Database || effective.Redis != current.Redis || effective.Log != current.Log || effective.Debug != current.Debug {
		t.Fatalf("restart-required config changed: %#v", effective)
	}
	wantSections := []string{"app", "database", "redis", "log", "debug"}
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
