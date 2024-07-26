package config

import (
	"fmt"
	"github.com/spf13/viper"
	"sync"
)

const DefaultCfg = "config.toml"

type Config struct {
	App   App
	Mysql Mysql
	Redis Redis
	Log   Log
	Debug Debug
}

var globalConfig *Config
var cfgOnce sync.Once

func GetConfig() *Config {
	cfgOnce.Do(func() {
		if globalConfig == nil {
			globalConfig = &Config{}
		}
	})
	return globalConfig
}

func InitConf(path string) (err error) {
	fmt.Println("load config file:", path)

	VIPER := viper.New()
	VIPER.SetConfigFile(path)
	VIPER.AddConfigPath(".")
	if err := VIPER.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to load config: %w \n", err)
	}

	if err := VIPER.Unmarshal(&globalConfig); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w \n", err)
	}
	return nil
}
