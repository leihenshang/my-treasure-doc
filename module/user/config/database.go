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

// Backup 控制 SQLite 定时备份（仅 sqlite 驱动生效，postgres 不启用）。
type Backup struct {
	// Enable 是否启用定时备份。
	Enable bool ``
	// Interval 备份周期，单位秒（例如 86400 = 每天）。
	Interval int ``
	// Dir 备份文件存放目录，缺省 backup。
	Dir string ``
	// Compress 是否 gzip 压缩（生成 .db.gz），节省磁盘占用。
	Compress bool ``
	// KeepDays 保留天数，超过该天数的旧备份会被自动清理；0 表示不清理。
	KeepDays int ``
}
