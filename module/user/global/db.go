package global

import (
	"fmt"
	"log"
	"os"
	"time"

	blogmodel "fastduck/treasure-doc/module/blog/data/model"
	"fastduck/treasure-doc/module/user/config"
	"fastduck/treasure-doc/module/user/data/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var TableMigrate = append([]schema.Tabler{
	&model.User{},
	&model.UserToken{},
}, blogTables()...)

func blogTables() []schema.Tabler {
	tables := blogmodel.Tables()
	result := make([]schema.Tabler, 0, len(tables))
	for _, table := range tables {
		result = append(result, table.(schema.Tabler))
	}
	return result
}

func initDatabase() error {
	cfg := GetConf()
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}

	newDb, err := openDatabaseWithConfig(cfg)
	if err != nil {
		return err
	}

	Db = newDb
	return nil
}

func openDatabaseWithConfig(cfg *config.Config) (*gorm.DB, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DbName,
		cfg.Database.Port,
		cfg.Database.SSLMode,
		cfg.Database.TimeZone)

	// table prefix
	tablePrefix := cfg.Database.TablePrefix

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // Slow SQL threshold
			LogLevel:      logger.Silent, // Zap level
			Colorful:      true,          // Disable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
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
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return db, nil
}

func closeDatabase(db *gorm.DB) error {
	if db == nil {
		return nil
	}

	sqlDb, err := db.DB()
	if err != nil {
		return err
	}

	return sqlDb.Close()
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
