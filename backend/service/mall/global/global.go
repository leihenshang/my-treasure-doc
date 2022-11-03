package global

import (
	"fastduck/treasure-doc/service/mall/config"
	"fmt"
	"log"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"fastduck/treasure-doc/service/mall/data/query"
)

var (
	Config   *config.Config
	Viper    *viper.Viper
	DbIns    *gorm.DB
	Zap      *zap.Logger
	ZapSugar *zap.SugaredLogger
	Redis    *redis.Client
	Trans    ut.Translator
)

// 环境标签
type srvType string

func GlobalInit(srv srvType) {
	fmt.Println("start global init")
	//读取配置
	configFile := config.ConfigFile
	fmt.Println("load config file:", path.Join(".", configFile))
	initConf(configFile)
	fmt.Println("初始化日志")
	initLogger()

	if Config.Redis.Enable {
		fmt.Println("初始化redis")
		initRedis()
	}

	fmt.Println("初始化mysql")
	initMysql()

	if srv != "cli" {
		fmt.Printf("初始化query,mode: %+v", Config.App.RunMode)
		initQuery(Config.App.RunMode == config.MODE_RELEASE)
	}

	fmt.Println("初始化validator")
	//初始化验证器
	if err := InitTrans("zh"); err != nil {
		log.Fatalf("init trans failed, err:%v\n", err)
		return
	}

	fmt.Println("global init successfully!")
}

func initQuery(IsRelease bool) {
	if IsRelease {
		query.SetDefault(DbIns)
	} else {
		query.SetDefault(DbIns.Debug())
	}
}

func initConf(configFile string) {
	Viper = viper.New()
	Viper.SetConfigFile(configFile)
	Viper.AddConfigPath(".")
	err := Viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	err = Viper.Unmarshal(&Config)
	if err != nil {
		panic("解析app配置失败")
	}
}

func initMysql() {
	mysqlPort := strconv.Itoa(Config.Mysql.Port)
	var err error
	//初始化数据库
	dsn := Config.Mysql.User + ":" + Config.Mysql.Password + "@tcp(" + Config.Mysql.Host + ":" + mysqlPort + ")/" +
		Config.Mysql.DbName + "?charset=" + Config.Mysql.Charset + "&parseTime=True&loc=Local"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // Slow SQL threshold
			LogLevel:      logger.Silent, // Log level
			Colorful:      true,          // Disable color
		},
	)

	DbIns, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		fmt.Println(err.Error())
		panic("初始化mysql失败")
	}

}

func initRedis() {
	//初始化redis
	Redis = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", Config.Redis.Host, Config.Redis.Port),
		Password: Config.Redis.Password, // no password set
		DB:       Config.Redis.DbId,     // use default DB)
	})

	if Redis.Ping().Err() != nil {
		panic("链接redis失败")
	}
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
		Trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", locale)
		}

		// 添加额外翻译
		_ = v.RegisterTranslation("required_with", Trans, func(ut ut.Translator) error {
			return ut.Add("required_with", "{0} 为必填字段!", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("required_with", fe.Field())
			return t
		})
		_ = v.RegisterTranslation("required_without", Trans, func(ut ut.Translator) error {
			return ut.Add("required_without", "{0} 为必填字段!", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("required_without", fe.Field())
			return t
		})
		_ = v.RegisterTranslation("required_without_all", Trans, func(ut ut.Translator) error {
			return ut.Add("required_without_all", "{0} 为必填字段!", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("required_without_all", fe.Field())
			return t
		})

		// 注册翻译器
		switch locale {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, Trans)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(v, Trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, Trans)
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

//handler中调用的错误翻译方法
func ErrResp(err error) string {
	errs, ok := err.(validator.ValidationErrors)
	// fmt.Println(reflect.TypeOf(err))
	if !ok {
		return errors.Wrap(err, "parse user data err!").Error()
	}
	errStruct := removeTopStruct(errs.Translate(Trans))
	for _, v := range errStruct {
		if val, ok := v.(string); ok {
			return val
		}
	}

	return "参数错误"
}
