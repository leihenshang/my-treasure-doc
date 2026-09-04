package config

import "testing"

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
