package global

import (
	"fastduck/treasure-doc/service/mall/config"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"time"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-redis/redis"
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
	initMysql(Config.Mysql.Debug)

	if srv != "cli" {
		fmt.Printf("初始化query,mode: %+v", Config.App.RunMode)
		initQuery()
	}

	fmt.Println("初始化validator")
	//初始化验证器
	if err := InitTrans("zh"); err != nil {
		log.Fatalf("init trans failed, err:%v\n", err)
		return
	}

	fmt.Println("global init successfully!")
}

func initQuery() {
	query.SetDefault(DbIns)
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

func initMysql(debug bool) {
	fmt.Printf("hhhh:%t \n", debug)
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

	// 启动debug模式
	if debug {
		DbIns = DbIns.Debug()
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
