package global

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"gorm.io/gorm/schema"

	"fastduck/treasure-doc/service/user/config"
	"fastduck/treasure-doc/service/user/data/model"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	CONFIG   *config.Config
	DB       *gorm.DB
	ZAP      *zap.Logger
	ZAPSUGAR *zap.SugaredLogger
	REDIS    *redis.Client
	TRANS    ut.Translator
)

var TableMigrate = []any{
	&model.Doc{},
	&model.DocGroup{},
	&model.GlobalConf{},
	&model.Team{},
	&model.TeamUser{},
	&model.User{},
	&model.VerifyCode{},
}

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
		err = initRedis()
		if err != nil {
			return
		}
		fmt.Println("初始化redis完成")
	}

	err = initMysql()
	if err != nil {
		return
	}
	fmt.Println("初始化mysql完成")

	//初始化验证器
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

func initMysql() error {
	mysqlPort := strconv.Itoa(CONFIG.Mysql.Port)
	var err error
	//初始化数据库
	dsn := CONFIG.Mysql.User + ":" + CONFIG.Mysql.Password + "@tcp(" + CONFIG.Mysql.Host + ":" + mysqlPort + ")/" +
		CONFIG.Mysql.DbName + "?charset=" + CONFIG.Mysql.Charset + "&parseTime=True&loc=Local"

	// table prefix
	tablePrefix := CONFIG.Mysql.TablePrefix

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // Slow SQL threshold
			LogLevel:      logger.Silent, // Log level
			Colorful:      true,          // Disable color
		},
	)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                                   newLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   tablePrefix,
			SingularTable: true,
			NameReplacer:  nil,
			NoLowerCase:   false,
		},
	})
	if err != nil {
		return fmt.Errorf("faile to initialize mysql,%w", err)
	}

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

// InitTrans 初始化翻译器
func InitTrans(locale string) (err error) {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return
	}

	// 注册一个获取json tag的自定义方法
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	uni := ut.New(en.New(), zh.New())
	TRANS, ok = uni.GetTranslator(locale)
	if !ok {
		return fmt.Errorf("uni.GetTranslator(%s) failed", locale)
	}

	switch locale {
	case "en":
		err = enTranslations.RegisterDefaultTranslations(v, TRANS)
	case "zh":
		err = zhTranslations.RegisterDefaultTranslations(v, TRANS)
	default:
		ZAPSUGAR.Error("failed to get translation from locale [%s],use default translation [en]", locale)
		err = enTranslations.RegisterDefaultTranslations(v, TRANS)
	}
	return

}

func addValueToMap(fields map[string]string) map[string]interface{} {
	res := make(map[string]interface{})
	for field, err := range fields {
		fieldArr := strings.SplitN(field, ".", 2)
		if len(fieldArr) > 1 {
			NewFields := map[string]string{fieldArr[1]: err}
			returnMap := addValueToMap(NewFields)
			if res[fieldArr[0]] != nil {
				for k, v := range returnMap {
					res[fieldArr[0]].(map[string]interface{})[k] = v
				}
			} else {
				res[fieldArr[0]] = returnMap
			}
			continue
		} else {
			res[field] = err
			continue
		}
	}
	return res
}

// removeTopStruct 去掉结构体名称前缀
func removeTopStruct(fields map[string]string) map[string]interface{} {
	lowerMap := map[string]string{}
	for field, err := range fields {
		fieldArr := strings.SplitN(field, ".", 2)
		lowerMap[fieldArr[1]] = err
	}
	res := addValueToMap(lowerMap)
	return res
}

// ErrResp 响应中调用的错误翻译方法
func ErrResp(err error) string {
	var errs validator.ValidationErrors
	if !errors.As(err, &errs) {
		return err.Error()
	}
	errStruct := removeTopStruct(errs.Translate(TRANS))
	for _, v := range errStruct {
		if val, ok := v.(string); ok {
			return val
		}
	}

	return "参数错误"
}

func migrateDbTable() error {
	fmt.Println("start migrate tables")
	defer fmt.Println("end of migration tables")
	if DB == nil {
		return fmt.Errorf("the DB is not initialize")
	}

	if err := DB.AutoMigrate(TableMigrate...); err != nil {
		return fmt.Errorf("failed to migrate tables,error:%v,table[%+v]", err, TableMigrate)
	}

	return nil
}
