package global

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	blogmodel "fastduck/treasure-doc/module/blog/data/model"
	"fastduck/treasure-doc/module/user/config"
	"fastduck/treasure-doc/module/user/data/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/glebarez/sqlite"
)

var TableMigrate = append([]schema.Tabler{
	&model.User{},
	&model.UserToken{},
	&model.SystemSetting{},
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

	gormConfig := &gorm.Config{
		Logger:                                   newLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   tablePrefix,
			SingularTable: true,
			NameReplacer:  nil,
			NoLowerCase:   false,
		},
	}

	var dialector gorm.Dialector
	switch cfg.Database.Driver {
	case config.DriverSQLite:
		dialector = sqlite.Open(cfg.Database.Dsn)
	case config.DriverPostgres:
		dialector = postgres.Open(cfg.Database.Dsn)
	default:
		return nil, fmt.Errorf("unsupported database driver: %q", cfg.Database.Driver)
	}

	// SQLite 需要数据库文件所在目录存在（driver 不会递归建目录）；
	// 从 dsn 解析出文件路径并创建父目录，便于配合挂载卷把库放到 ./data 子目录。
	if cfg.Database.Driver == config.DriverSQLite {
		if dbFile, perr := parseSqliteFileFromDsn(cfg.Database.Dsn); perr == nil {
			if dir := filepath.Dir(dbFile); dir != "" && dir != "." {
				if mkErr := os.MkdirAll(dir, 0o755); mkErr != nil {
					return nil, fmt.Errorf("failed to create sqlite db dir %q: %w", dir, mkErr)
				}
			}
		}
	}

	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database (driver=%s): %w", cfg.Database.Driver, err)
	}

	// SQLite 为单写者模型，限制连接池并在 DSN 已设置 busy_timeout/WAL，
	// 避免高并发下出现 database is locked。
	if cfg.Database.Driver == config.DriverSQLite {
		if sqlDb, err := db.DB(); err == nil {
			sqlDb.SetMaxOpenConns(1)
		}
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

	// 兼容旧库：收藏集不再使用公开 ID（public_id），移除历史遗留列。
	// AutoMigrate 只建表/加列不会删列，旧库若不清理，NOT NULL 的 public_id 会让新建收藏集失败。
	if Db.Migrator().HasColumn(&blogmodel.Bookmark{}, "public_id") {
		if err := Db.Migrator().DropColumn(&blogmodel.Bookmark{}, "public_id"); err != nil {
			return fmt.Errorf("failed to drop deprecated column td_blog_bookmark.public_id: %v", err)
		}
	}

	// 兼容旧库：模块「工具」曾命名为「利器」。seed 只增不改、服务端对非空 title 原样返回，
	// 因此改 seed/默认值不会让既有站点自愈，这里对存量的 title 做一次性幂等回写。
	if err := migrateSiteModuleTitle(); err != nil {
		return fmt.Errorf("failed to migrate site module title: %v", err)
	}

	return nil
}

// migrateSiteModuleTitle 把站点模块集合里 tools.title 旧的「利器」重写为「工具」。
// 用 Go + encoding/json 而非原生 SQL：JSONB/文本改写在 SQLite 与 Postgres 语法不同，
// Go 方案跨驱动、幂等（无命中则跳过），并与 seed/默认配置保持一致。
func migrateSiteModuleTitle() error {
	var sites []blogmodel.Site
	if err := Db.Find(&sites).Error; err != nil {
		return err
	}
	for _, site := range sites {
		if len(site.Modules) == 0 {
			continue
		}
		type moduleLike struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		}
		var modules []moduleLike
		if err := json.Unmarshal(site.Modules, &modules); err != nil {
			// 非严格解析失败（结构非法）时跳过，避免回写破坏异常数据。
			continue
		}
		changed := false
		for i := range modules {
			if modules[i].ID == "tools" && modules[i].Title == "利器" {
				modules[i].Title = "工具"
				changed = true
			}
		}
		if !changed {
			continue
		}
		data, err := json.Marshal(modules)
		if err != nil {
			return err
		}
		if err := Db.Model(&blogmodel.Site{}).Where("id = ?", site.ID).Update("modules", data).Error; err != nil {
			return err
		}
	}
	return nil
}
