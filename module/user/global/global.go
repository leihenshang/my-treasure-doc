package global

import (
	"fmt"
	"log"
	"sync"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"fastduck/treasure-doc/module/user/config"
)

var (
	Conf  *config.Config
	Db    *gorm.DB
	Zap   *zap.Logger
	Log   *zap.SugaredLogger
	Redis *redis.Client
	Trans ut.Translator

	confMu sync.RWMutex
	connMu sync.Mutex
)

func GetConf() *config.Config {
	confMu.RLock()
	defer confMu.RUnlock()
	return Conf
}

func setConf(cfg *config.Config) {
	confMu.Lock()
	Conf = cfg
	confMu.Unlock()
}

func InitModule(cfgPath string) (destructFunc func(), err error) {
	if err = config.InitConf(cfgPath); err != nil {
		return
	}
	fmt.Println("初始化配置完成")
	setConf(config.GetConfig())

	if err = initLog(); err != nil {
		return
	}
	fmt.Println("初始化日志完成")

	if GetConf().Redis.Enable {
		if err = initRedis(); err != nil {
			return
		}
		fmt.Println("初始化redis完成")
	}

	if err = initMysql(); err != nil {
		return
	}
	fmt.Println("初始化mysql完成")

	if err = InitTrans("zh"); err != nil {
		log.Fatalf("init trans failed, err:%v\n", err)
		return
	}
	fmt.Println("初始化validator完成")

	if err = initConfigHotReload(); err != nil {
		return
	}
	fmt.Println("初始化配置热更新完成")

	return destructModule(), migrateDbTable()
}

func initConfigHotReload() error {
	return config.WatchConf(func(cfg *config.Config) {
		oldCfg := GetConf()
		setConf(cfg)
		if err := applyConnectionHotReload(oldCfg, cfg); err != nil {
			if Log != nil {
				Log.Errorf("config hot reload connection refresh failed: %v", err)
			} else {
				fmt.Printf("config hot reload connection refresh failed: %v\n", err)
			}
		}

		if Log != nil {
			Log.Infof("config hot reloaded")
			return
		}
		fmt.Println("config hot reloaded")
	})
}

func applyConnectionHotReload(oldCfg, newCfg *config.Config) error {
	if oldCfg == nil || newCfg == nil {
		return nil
	}

	connMu.Lock()
	defer connMu.Unlock()

	if oldCfg.Mysql != newCfg.Mysql {
		if err := reloadMysql(newCfg); err != nil {
			return err
		}
		if Log != nil {
			Log.Infof("mysql connection refreshed by config hot reload")
		}
	}

	if oldCfg.Redis != newCfg.Redis {
		if err := reloadRedis(newCfg); err != nil {
			return err
		}
		if Log != nil {
			Log.Infof("redis connection refreshed by config hot reload")
		}
	}

	return nil
}

func reloadMysql(cfg *config.Config) error {
	newDb, err := openMysqlWithConfig(cfg)
	if err != nil {
		return err
	}

	oldDb := Db
	Db = newDb

	if oldDb != nil {
		if err := closeMysql(oldDb); err != nil {
			if Log != nil {
				Log.Warnf("close old mysql connection failed: %v", err)
			} else {
				log.Printf("close old mysql connection failed: %v\n", err)
			}
		}
	}

	return nil
}

func reloadRedis(cfg *config.Config) error {
	if !cfg.Redis.Enable {
		if Redis != nil {
			if err := Redis.Close(); err != nil && Log != nil {
				Log.Warnf("close old redis connection failed: %v", err)
			}
			Redis = nil
		}
		return nil
	}

	newRedis, err := initRedisWithConfig(cfg)
	if err != nil {
		return err
	}

	oldRedis := Redis
	Redis = newRedis
	if oldRedis != nil {
		if err := oldRedis.Close(); err != nil && Log != nil {
			Log.Warnf("close old redis connection failed: %v", err)
		}
	}

	return nil
}

func destructModule() func() {
	return func() {
		if Zap != nil {
			if err := Zap.Sync(); err != nil {
				log.Printf("failed to sync Zap log,error:%v \n", err)
			}
		}

		if Log != nil {
			if err := Log.Sync(); err != nil {
				log.Printf("failed to sync Log log,error:%v \n", err)
			}
		}

		if Db != nil {
			if err := closeMysql(Db); err != nil {
				log.Printf("failed to close Db,error:%v \n", err)
			}
		}

		if Redis != nil {
			if err := Redis.Close(); err != nil {
				log.Printf("failed to close Redis,error:%v \n", err)
			}
		}
	}
}

func InitRestPwd(cfgPath string) error {
	if err := config.InitConf(cfgPath); err != nil {
		return err
	}
	setConf(config.GetConfig())
	fmt.Println("初始化配置完成")
	err := initMysql()
	if err != nil {
		return err
	}
	fmt.Println("初始化mysql完成")
	fmt.Printf("\n")
	return nil
}

func initLog() error {
	return initZapLogger()
}

func initRedis() error {
	return reloadRedis(GetConf())
}

func initRedisWithConfig(cfg *config.Config) (*redis.Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password, // no password set
		DB:       cfg.Redis.DbId,     // use default Db)
	})

	if err := client.Ping().Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("failed to initilize redis,%w", err)
	}

	return client, nil
}
