package config

import (
	"fmt"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

const DefaultConfig = "config.toml"

type Config struct {
	App      App
	Database Database
	Redis    Redis
	Log      Log
	Debug    Debug
	BlogSeed BlogSeed
	Captcha  Captcha
	Backup   Backup
}

// applyConfigDefaults 设置配置文件可省略项的默认值。
// 验证码默认开启，避免新增配置项后因为漏配而失去登录保护。
func applyConfigDefaults(v *viper.Viper) {
	v.SetDefault("captcha.enabled", true)
	// SQLite 定时备份默认值：每日一次、压缩、保留 7 天；默认关闭，需显式 enable。
	v.SetDefault("backup.enable", false)
	v.SetDefault("backup.interval", 86400)
	v.SetDefault("backup.dir", "backup")
	v.SetDefault("backup.compress", true)
	v.SetDefault("backup.keepDays", 7)
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

func PublishConfig(cfg *Config) {
	if cfg == nil {
		return
	}
	snapshot := *cfg
	cfgMu.Lock()
	globalConfig = &snapshot
	cfgMu.Unlock()
}

func InitConf(path string) (err error) {
	fmt.Println("load config file:", path)

	v := viper.New()
	v.SetConfigFile(path)
	v.AddConfigPath(".")
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to load config: %w \n", err)
	}

	applyConfigDefaults(v)
	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w \n", err)
	}

	PublishConfig(cfg)
	cfgMu.Lock()
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

	applyConfigDefaults(v)
	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal reloaded config: %w", err)
	}

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
