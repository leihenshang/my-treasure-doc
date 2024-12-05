package global

import (
	"fmt"
	"log"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"fastduck/treasure-doc/service/user/config"
)

var (
	Conf  *config.Config
	Db    *gorm.DB
	Zap   *zap.Logger
	Log   *zap.SugaredLogger
	Redis *redis.Client
	Trans ut.Translator
)

func InitModule(cfgPath string) (destructFunc func(), err error) {
	if err = config.InitConf(cfgPath); err != nil {
		return
	}
	fmt.Println("初始化配置完成")
	Conf = config.GetConfig()

	if err = initLog(); err != nil {
		return
	}
	fmt.Println("初始化日志完成")

	if Conf.Redis.Enable {
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
			if db, err := Db.DB(); err != nil {
				log.Printf("failed to get Db,error:%v \n", err)
			} else if err = db.Close(); err != nil {
				log.Printf("failed to close Db,error:%v \n", err)
			}
		}
	}
}

func InitRestPwd(cfgPath string) error {
	if err := config.InitConf(cfgPath); err != nil {
		return err
	}
	Conf = config.GetConfig()
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
	Redis = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", Conf.Redis.Host, Conf.Redis.Port),
		Password: Conf.Redis.Password, // no password set
		DB:       Conf.Redis.DbId,     // use default Db)
	})

	if err := Redis.Ping().Err(); err != nil {
		return fmt.Errorf("failed to initilize redis,%w", err)
	}
	return nil
}
