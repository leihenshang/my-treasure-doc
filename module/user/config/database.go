package config

// Driver 支持的数据库方言。默认 sqlite（零依赖、单文件），可选 postgres。
const (
	DriverSQLite   = "sqlite"
	DriverPostgres = "postgres"
)

type Database struct {
	// Driver 选择数据库方言："sqlite"（默认）或 "postgres"。
	Driver string ``
	// Dsn 连接串：
	//   - sqlite：file:./treasure_doc.db?_busy_timeout=5000&_journal_mode=WAL&_foreign_keys=on
	//   - postgres：host=127.0.0.1 user=postgres password=postgres dbname=treasure_doc port=5432 sslmode=disable TimeZone=Asia/Shanghai
	Dsn string ``
	// TablePrefix 表名前缀，沿用现有 td_ 约定。
	TablePrefix string ``
}
