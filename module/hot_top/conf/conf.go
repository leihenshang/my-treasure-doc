package conf

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

const (
	// gin configuration
	GinModeRelease = "release"
	GinModeDev     = "dev"

	// default configurations for hot
	DefaultConfig             = "config.toml"
	DefaultPort               = 2025
	DefaultHost               = ""
	DefaultHotExpiredInterval = time.Second * 10
	DefaultHotPullInterval    = time.Hour
	DefaultHotCacheFilePath   = "hot_cache"
)

var globalConfig *Conf = &Conf{}
var confOnce sync.Once
var vip *viper.Viper = viper.New()
var fsEventTicker *time.Timer
var lock *sync.RWMutex = &sync.RWMutex{}

type Conf struct {
	App App
	Ai  Ai
	Hot Hot
}

func GetConf() *Conf {
	lock.RLock()
	defer lock.RUnlock()
	return globalConfig
}

func (c *Conf) Validate(vip *viper.Viper) (err error) {
	lock.Lock()
	defer lock.Unlock()

	if err = vip.Unmarshal(&c); err != nil {
		err = fmt.Errorf("failed to unmarshal config: %w", err)
		return
	}

	if err = c.Ai.Validate(); err != nil {
		err = fmt.Errorf("failed to validate ai config: %w", err)
		return
	}

	if err = c.App.Validate(); err != nil {
		err = fmt.Errorf("failed to validate app config: %w", err)
		return
	}

	if err = c.Hot.Validate(); err != nil {
		err = fmt.Errorf("failed to validate hot config: %w", err)
		return
	}
	return nil
}

func (c *Conf) PrintJson() string {
	b, err := json.Marshal(c)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func InitConf(path string) (err error) {
	confOnce.Do(func() {
		vip.SetConfigFile(path)
		vip.OnConfigChange(func(in fsnotify.Event) {
			if fsEventTicker != nil {
				fsEventTicker.Stop()
			}

			fsEventTicker = time.AfterFunc(time.Millisecond*100, func() {
				if err = globalConfig.Validate(vip); err != nil {
					log.Printf("failed to validate config [%s] when reload a file: %v", in.Name, err)
				} else {
					log.Printf("config [%s] is updated: %s \n", in.Name, globalConfig.PrintJson())
				}
			})
		})
		vip.WatchConfig()
		if err = vip.ReadInConfig(); err != nil {
			err = fmt.Errorf("failed to load config from [%s], err: %w", path, err)
			return
		}

		if err = globalConfig.Validate(vip); err != nil {
			err = fmt.Errorf("failed to unmarshal config: %w", err)
			return
		}
	})
	return err
}

type App struct {
	Host    string
	Port    int
	Name    string
	RunMode string
}

func (app *App) GetAddr() string {
	return fmt.Sprintf("%s:%d", app.Host, app.Port)
}

func (a *App) Validate() error {
	if a.Port == 0 {
		a.Port = DefaultPort
	}
	if a.Host == "" {
		a.Host = DefaultHost
	}
	return nil
}

type Ai struct {
	DeepSeekToken string
}

func (a *Ai) Validate() error {
	if a.DeepSeekToken == "" {
		log.Println("deepseek token is empty, please check your config")
	}
	return nil
}

type Hot struct {
	ExpiredCheckInterval       string
	ExpiredCheckIntervalParsed time.Duration `toml:"-"`
	HotPullInterval            string
	HotPullIntervalParsed      time.Duration `toml:"-"`
	HotFileCachePath           string
	HotConf                    map[Source]*HotConf
}

func (h *Hot) Validate() error {
	if h.ExpiredCheckInterval == "" {
		h.ExpiredCheckIntervalParsed = DefaultHotExpiredInterval
	} else if duration, err := time.ParseDuration(h.ExpiredCheckInterval); err != nil || duration <= 0 {
		log.Printf("parse hot expired check time failed,err: [%v], use default value: %v \n", err, DefaultHotExpiredInterval)
		h.ExpiredCheckIntervalParsed = DefaultHotExpiredInterval
	} else {
		h.ExpiredCheckIntervalParsed = duration
	}

	if h.HotPullInterval == "" {
		h.HotPullIntervalParsed = DefaultHotPullInterval
	} else if duration, err := time.ParseDuration(h.HotPullInterval); err != nil || duration <= 0 {
		log.Printf("parse hot expired check time failed,err: [%v], use default value: %v \n", err, DefaultHotPullInterval)
		h.HotPullIntervalParsed = DefaultHotPullInterval
	} else {
		h.HotPullIntervalParsed = duration
	}

	if h.HotFileCachePath == "" {
		h.HotFileCachePath = DefaultHotCacheFilePath
	}

	h.HotConf = make(map[Source]*HotConf)
	for _, v := range HotConfList {
		h.HotConf[v.Source] = v
	}

	return nil
}
