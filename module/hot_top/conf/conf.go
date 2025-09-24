package conf

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/spf13/viper"
)

const DefaultConfig = "config.toml"
const DefaultPort = 2025
const DefaultHost = ""
const DefaultHotExpiredInterval = time.Second * 10

type Conf struct {
	App App
	Ai  Ai
	Hot Hot
}

const GinModeRelease = "release"
const GinModeDev = "dev"

type App struct {
	Host    string
	Port    int
	Name    string
	RunMode string
}

type Ai struct {
	DeepSeekToken string
}

type Hot struct {
	ExpiredCheckInterval       string
	ExpiredCheckIntervalParsed time.Duration `toml:"-"`
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

			if duration, err := time.ParseDuration(globalConfig.Hot.ExpiredCheckInterval); err != nil || duration <= 0 {
				log.Println("parse hot expired check time failed,err: [%v], use default value: ", err, DefaultHotExpiredInterval)
				globalConfig.Hot.ExpiredCheckIntervalParsed = DefaultHotExpiredInterval
			} else {
				globalConfig.Hot.ExpiredCheckIntervalParsed = duration
			}
		}
	})

	return nil
}
