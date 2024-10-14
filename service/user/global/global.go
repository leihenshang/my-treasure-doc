package global

import (
	"fastduck/treasure-doc/service/user/config"
	"fmt"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"log"
)

var (
	CONFIG   *config.Config
	DB       *gorm.DB
	ZAP      *zap.Logger
	ZAPSUGAR *zap.SugaredLogger
	REDIS    *redis.Client
	TRANS    ut.Translator
)

func InitModule(cfgPath string) (destructFunc func(), err error) {
	if err = config.InitConf(cfgPath); err != nil {
		return
	}
	fmt.Println("初始化配置完成")
	CONFIG = config.GetConfig()

	if err = initLog(); err != nil {
		return
	}
	fmt.Println("初始化日志完成")

	if CONFIG.Redis.Enable {
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

	return destructModule(), migrateDbTable()
}

func destructModule() func() {
	return func() {
		if ZAP != nil {
			if err := ZAP.Sync(); err != nil {
				log.Printf("failed to sync ZAP log,error:%v \n", err)
			}
		}

		if ZAPSUGAR != nil {
			if err := ZAPSUGAR.Sync(); err != nil {
				log.Printf("failed to sync ZAPSUGAR log,error:%v \n", err)
			}
		}

		if DB != nil {
			if db, err := DB.DB(); err != nil {
				log.Printf("failed to get DB,error:%v \n", err)
			} else if err = db.Close(); err != nil {
				log.Printf("failed to close DB,error:%v \n", err)
			}
		}
	}
}

func InitRestPwd(cfgPath string) error {
	if err := config.InitConf(cfgPath); err != nil {
		return err
	}
	CONFIG = config.GetConfig()
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
	REDIS = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", CONFIG.Redis.Host, CONFIG.Redis.Port),
		Password: CONFIG.Redis.Password, // no password set
		DB:       CONFIG.Redis.DbId,     // use default DB)
	})

	if err := REDIS.Ping().Err(); err != nil {
		return fmt.Errorf("failed to initilize redis,%w", err)
	}
	return nil
}
