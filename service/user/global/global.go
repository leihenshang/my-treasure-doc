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
	//if err := InitTrans("zh"); err != nil {
	//	log.Fatalf("init trans failed, err:%v\n", err)
	//	return
	//}
	//fmt.Println("初始化validator完成")

	return destructModule(), migrateTable()
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
	//修改gin框架中的Validator属性，实现自定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册一个获取json tag的自定义方法
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		zhT := zh.New() //中文翻译器
		enT := en.New() //英文翻译器

		// 第一个参数是备用（fallback）的语言环境
		// 后面的参数是应该支持的语言环境（支持多个）
		// uni := ut.New(zhT, zhT) 也是可以的
		uni := ut.New(enT, zhT, enT)

		// locale 通常取决于 http 请求头的 'Accept-Language'
		var ok bool
		// 也可以使用 uni.FindTranslator(...) 传入多个locale进行查找
		TRANS, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", locale)
		}

		// 添加额外翻译
		_ = v.RegisterTranslation("required_with", TRANS, func(ut ut.Translator) error {
			return ut.Add("required_with", "{0} 为必填字段!", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("required_with", fe.Field())
			return t
		})
		_ = v.RegisterTranslation("required_without", TRANS, func(ut ut.Translator) error {
			return ut.Add("required_without", "{0} 为必填字段!", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("required_without", fe.Field())
			return t
		})
		_ = v.RegisterTranslation("required_without_all", TRANS, func(ut ut.Translator) error {
			return ut.Add("required_without_all", "{0} 为必填字段!", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("required_without_all", fe.Field())
			return t
		})

		// 注册翻译器
		switch locale {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, TRANS)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(v, TRANS)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, TRANS)
		}
		return
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

// 去掉结构体名称前缀
func removeTopStruct(fields map[string]string) map[string]interface{} {
	lowerMap := map[string]string{}
	for field, err := range fields {
		fieldArr := strings.SplitN(field, ".", 2)
		lowerMap[fieldArr[1]] = err
	}
	res := addValueToMap(lowerMap)
	return res
}

// handler中调用的错误翻译方法
func ErrResp(err error) string {
	var errs validator.ValidationErrors
	ok := errors.As(err, &errs)
	fmt.Println(reflect.TypeOf(err))
	if !ok {
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

func migrateTable() error {
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
