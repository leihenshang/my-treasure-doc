package global

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"sync"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"gorm.io/gorm"

	blogseed "fastduck/treasure-doc/module/blog/seed"
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
	if err = validateStartupConfig(GetConf()); err != nil {
		return
	}

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

	if err = initDatabase(); err != nil {
		return
	}
	fmt.Println("初始化PostgreSQL完成")

	if err = InitTrans("zh"); err != nil {
		log.Fatalf("init trans failed, err:%v\n", err)
		return
	}
	fmt.Println("初始化validator完成")

	if err = initConfigHotReload(); err != nil {
		return
	}
	fmt.Println("初始化配置热更新完成")

	destructFunc = destructModule()
	if err = migrateDbTable(); err != nil {
		destructFunc()
		destructFunc = nil
		return
	}
	if err = seedBlogData(); err != nil {
		destructFunc()
		destructFunc = nil
		return
	}
	return
}

func seedBlogData() error {
	cfg := GetConf()
	if cfg == nil || !cfg.BlogSeed.Enabled {
		return nil
	}
	result, err := blogseed.Seed(Db, blogseed.Options{
		Enabled: cfg.BlogSeed.Enabled, AllowRemote: cfg.BlogSeed.AllowRemote,
		RestoreDeleted: cfg.BlogSeed.RestoreDeleted, Release: cfg.App.IsRelease(),
		Driver: cfg.Database.Driver, Dsn: cfg.Database.Dsn,
	})
	if err != nil {
		return fmt.Errorf("failed to seed blog data: %w", err)
	}
	message := fmt.Sprintf("blog seed completed for %s: created=%d existing=%d restored=%d skipped=%d", cfg.Database.Driver, result.Created, result.Existing, result.Restored, result.Skipped)
	if Log != nil {
		Log.Info(message)
	} else {
		fmt.Println(message)
	}
	return nil
}

func initConfigHotReload() error {
	return config.WatchConf(func(candidate *config.Config) {
		effective, restartRequired := effectiveHotReloadConfig(GetConf(), candidate)
		setConf(effective)
		config.PublishConfig(effective)

		if len(restartRequired) > 0 {
			message := fmt.Sprintf("config changes require restart and were ignored: %s", strings.Join(restartRequired, ", "))
			if Log != nil {
				Log.Warn(message)
			} else {
				fmt.Println(message)
			}
		}
		if Log != nil {
			Log.Infof("business config hot reloaded")
		} else {
			fmt.Println("business config hot reloaded")
		}
	})
}

func effectiveHotReloadConfig(current, candidate *config.Config) (*config.Config, []string) {
	if current == nil {
		return candidate, nil
	}
	if candidate == nil {
		snapshot := *current
		return &snapshot, nil
	}

	effective := *current
	effective.App.RegisterEnabled = candidate.App.RegisterEnabled
	effective.Captcha = candidate.Captcha
	restartRequired := make([]string, 0, 6)
	currentApp := current.App
	candidateApp := candidate.App
	currentApp.RegisterEnabled = false
	candidateApp.RegisterEnabled = false
	if !reflect.DeepEqual(currentApp, candidateApp) {
		restartRequired = append(restartRequired, "app")
	}
	if !reflect.DeepEqual(current.Database, candidate.Database) {
		restartRequired = append(restartRequired, "database")
	}
	if !reflect.DeepEqual(current.Redis, candidate.Redis) {
		restartRequired = append(restartRequired, "redis")
	}
	if !reflect.DeepEqual(current.Log, candidate.Log) {
		restartRequired = append(restartRequired, "log")
	}
	if !reflect.DeepEqual(current.Debug, candidate.Debug) {
		restartRequired = append(restartRequired, "debug")
	}
	if !reflect.DeepEqual(current.BlogSeed, candidate.BlogSeed) {
		restartRequired = append(restartRequired, "blogSeed")
	}
	return &effective, restartRequired
}

func validateStartupConfig(cfg *config.Config) error {
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}
	if cfg.App.IsRelease() && cfg.Debug.EnableMockLogin {
		return fmt.Errorf("debug mock login is forbidden in release mode")
	}
	if cfg.App.IsRelease() && cfg.BlogSeed.Enabled {
		return fmt.Errorf("blog seed is forbidden in release mode")
	}
	// 数据库驱动校验：缺省视为 sqlite；非法驱动在启动期即失败，避免静默用空库。
	if cfg.Database.Driver == "" {
		cfg.Database.Driver = config.DriverSQLite
	}
	if cfg.Database.Driver != config.DriverSQLite && cfg.Database.Driver != config.DriverPostgres {
		return fmt.Errorf("unsupported database driver: %q (want %s or %s)", cfg.Database.Driver, config.DriverSQLite, config.DriverPostgres)
	}
	switch cfg.Database.Driver {
	case config.DriverSQLite:
		// 未配置 dsn 时给出默认单文件库，保证零配置即可启动。
		if cfg.Database.Dsn == "" {
			cfg.Database.Dsn = "file:./treasure_doc.db?_busy_timeout=5000&_journal_mode=WAL&_foreign_keys=on"
		}
	case config.DriverPostgres:
		if cfg.Database.Dsn == "" {
			return fmt.Errorf("database dsn is required for driver %q", cfg.Database.Driver)
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
			if err := closeDatabase(Db); err != nil {
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
	err := initDatabase()
	if err != nil {
		return err
	}
	fmt.Println("初始化PostgreSQL完成")
	fmt.Printf("\n")
	return nil
}

func initLog() error {
	return initZapLogger()
}

func initRedis() error {
	client, err := initRedisWithConfig(GetConf())
	if err != nil {
		return err
	}
	Redis = client
	return nil
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
