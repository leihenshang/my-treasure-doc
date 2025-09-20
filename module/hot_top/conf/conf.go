package conf

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

const DefaultConfig = "config.toml"
const DefaultPort = 2025
const DefaultHost = ""

type Conf struct {
	App App
	Ai  Ai
}

const GinModeRelease = "release"
const GinModeDev = "dev"

type App struct {
	Host            string
	Port            int
	Name            string
	RunMode         string
	RegisterEnabled bool
}

type Ai struct {
	DeepSeekToken string
}

func (app *App) IsRelease() bool {
	return app.RunMode == GinModeRelease
}

func (app *App) IsDev() bool {
	return app.RunMode == GinModeDev
}

func (app *App) GetAddr() string {
	return fmt.Sprintf("%s:%d", app.Host, app.Port)
}

var globalConfig *Conf
var confOnce sync.Once

func GetConf() *Conf {
	return globalConfig
}

func InitConf(path string) (err error) {
	confOnce.Do(func() {
		if globalConfig == nil {
			vp := viper.New()
			vp.SetConfigFile(path)
			vp.AddConfigPath(".")
			if err = vp.ReadInConfig(); err != nil {
				err = fmt.Errorf("failed to load config: %w", err)
				return
			}

			if err = vp.Unmarshal(&globalConfig); err != nil {
				err = fmt.Errorf("failed to unmarshal config: %w", err)
				return
			}
			if globalConfig.App.Port == 0 {
				globalConfig.App.Port = DefaultPort
			}
			if globalConfig.App.Host == "" {
				globalConfig.App.Host = DefaultHost
			}
		}
	})

	return nil
}
