package global

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"fastduck/treasure-doc/module/user/data/model"
)

var TableMigrate = []schema.Tabler{
	&model.Doc{},
	&model.DocGroup{},
	&model.GlobalConf{},
	&model.Team{},
	&model.TeamUser{},
	&model.User{},
	&model.VerifyCode{},
	&model.Note{},
	&model.DocHistory{},
	&model.UserConf{},
	&model.UserToken{},
	&model.Room{},
}

func initMysql() error {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		Conf.Mysql.User,
		Conf.Mysql.Password,
		Conf.Mysql.Host,
		Conf.Mysql.Port,
		Conf.Mysql.DbName,
		Conf.Mysql.Charset)

	// table prefix
	tablePrefix := Conf.Mysql.TablePrefix

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // Slow SQL threshold
			LogLevel:      logger.Silent, // Zap level
			Colorful:      true,          // Disable color
		},
	)

	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
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
	if Db == nil {
		return fmt.Errorf("the Db is not initialize")
	}

	for _, t := range TableMigrate {
		if err := Db.AutoMigrate(t); err != nil {
			return fmt.Errorf("failed to migrate tables,error:%v,table[%#v]", err, t.TableName())
		}
	}

	return nil
}
