package global

import (
	"fastduck/treasure-doc/service/user/data/model"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"strconv"
	"time"
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
