package conf

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

const DefaultConfig = "config.toml"

type Conf struct {
	DeepSeekToken string
}

var globalConfig *Conf
var confOnce sync.Once

func GetConf() *Conf {
	return globalConfig
}

func InitConf(path string) (err error) {
	vp := viper.New()
	vp.SetConfigFile(path)
	vp.AddConfigPath(".")
	if err := vp.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to load config: %w \n", err)
	}

	if err := vp.Unmarshal(&globalConfig); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w \n", err)
	}
	return nil
}
