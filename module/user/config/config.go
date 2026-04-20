package config

import (
	"fmt"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

const DefaultConfig = "config.toml"

type Config struct {
	App   App
	Mysql Mysql
	Redis Redis
	Log   Log
	Debug Debug
}

var globalConfig *Config
var cfgMu sync.RWMutex
var cfgViper *viper.Viper

func GetConfig() *Config {
	cfgMu.RLock()
	cfg := globalConfig
	cfgMu.RUnlock()
	if cfg != nil {
		return cfg
	}

	cfgMu.Lock()
	defer cfgMu.Unlock()
	if globalConfig == nil {
		globalConfig = &Config{}
	}
	return globalConfig
}

func InitConf(path string) (err error) {
	fmt.Println("load config file:", path)

	v := viper.New()
	v.SetConfigFile(path)
	v.AddConfigPath(".")
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to load config: %w \n", err)
	}

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w \n", err)
	}

	cfgMu.Lock()
	globalConfig = cfg
	cfgViper = v
	cfgMu.Unlock()

	return nil
}

func ReloadConf() (*Config, error) {
	cfgMu.RLock()
	v := cfgViper
	cfgMu.RUnlock()
	if v == nil {
		return nil, fmt.Errorf("config not initialized")
	}

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to reload config: %w", err)
	}

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal reloaded config: %w", err)
	}

	cfgMu.Lock()
	globalConfig = cfg
	cfgMu.Unlock()

	return cfg, nil
}

func WatchConf(onChange func(*Config)) error {
	cfgMu.RLock()
	v := cfgViper
	cfgMu.RUnlock()
	if v == nil {
		return fmt.Errorf("config not initialized")
	}

	v.OnConfigChange(func(event fsnotify.Event) {
		cfg, err := ReloadConf()
		if err != nil {
			fmt.Printf("config hot reload failed: %v\n", err)
			return
		}

		fmt.Printf("config hot reloaded: %s (%s)\n", event.Name, event.Op.String())
		if onChange != nil {
			onChange(cfg)
		}
	})
	v.WatchConfig()

	return nil
}
