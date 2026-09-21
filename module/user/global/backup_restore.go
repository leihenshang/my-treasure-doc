package global

import (
	"archive/tar"
	"compress/gzip"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fastduck/treasure-doc/module/user/config"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// OpenedArchive 描述一次解包结果。
type OpenedArchive struct {
	Dir       string // 解包根目录
	Manifest  *Manifest
	DBPath    string   // 解包出的 treasure_doc.db
	FilesRoot string   // 解包出的 files/ 目录（可能不存在）
	Extracted []string // 已落盘条目路径，供调用方清理
}

// ExtractBackupArchive 安全解包 tar.gz 到 destDir，返回包内容描述。
// 拒绝路径穿越（绝对路径、../、软链逃逸）。
func ExtractBackupArchive(src, destDir string) (*OpenedArchive, error) {
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return nil, err
	}
	f, err := os.Open(src)
	if err != nil {
		return nil, fmt.Errorf("打开备份包失败: %w", err)
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return nil, fmt.Errorf("备份包不是有效的 gzip: %w", err)
	}
	defer gz.Close()

	destAbs, err := filepath.Abs(destDir)
	if err != nil {
		return nil, err
	}

	archived := &OpenedArchive{Dir: destDir, FilesRoot: filepath.Join(destDir, archiveFileRoot)}
	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("解析备份包失败: %w", err)
		}
		name := filepath.ToSlash(strings.TrimPrefix(hdr.Name, "./"))
		if !backupArchiveLinkSafe(name) {
			return nil, fmt.Errorf("备份包包含非法路径: %q", hdr.Name)
		}
		target, err := filepath.Abs(filepath.Join(destAbs, filepath.FromSlash(name)))
		if err != nil {
			return nil, err
		}
		if !strings.HasPrefix(target, destAbs+string(os.PathSeparator)) {
			return nil, fmt.Errorf("备份包路径越界: %q", hdr.Name)
		}
		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0o755); err != nil {
				return nil, err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return nil, err
			}
			out, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
			if err != nil {
				return nil, err
			}
			if _, err := io.Copy(out, tr); err != nil {
				out.Close()
				return nil, err
			}
			out.Close()
			archived.Extracted = append(archived.Extracted, target)
		default:
			// 拒绝硬链接/软链等，避免越权写入
			return nil, fmt.Errorf("备份包包含不支持的条目类型: %q", hdr.Name)
		}
	}

	dbPath := filepath.Join(destDir, archiveDBName)
	if _, err := os.Stat(dbPath); err != nil {
		return nil, fmt.Errorf("备份包缺少数据库文件 %s", archiveDBName)
	}
	data, err := os.ReadFile(filepath.Join(destDir, archiveManifestName))
	if err != nil {
		return nil, fmt.Errorf("备份包缺少清单文件 %s", archiveManifestName)
	}
	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("备份包清单解析失败: %w", err)
	}
	archived.DBPath = dbPath
	archived.Manifest = &manifest
	return archived, nil
}

// ValidateBackupSchema 校验备份清单版本与当前版本一致。
func ValidateBackupSchema(manifest *Manifest) error {
	if manifest == nil {
		return fmt.Errorf("备份包清单为空")
	}
	if manifest.SchemaVersion != BackupSchemaVersion {
		return fmt.Errorf("备份包版本不兼容：备份为 v%d，当前支持 v%d", manifest.SchemaVersion, BackupSchemaVersion)
	}
	return nil
}

// ValidateSqliteFile 对解包的数据库做完整性校验（integrity_check 与 foreign_key_check）。
func ValidateSqliteFile(dbPath string) error {
	db, err := gorm.Open(sqlite.Open("file:"+filepath.ToSlash(dbPath)+"?_mode=ro&_busy_timeout=3000"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return fmt.Errorf("无法打开备份数据库: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	var integrity string
	if err := sqlDB.QueryRow("PRAGMA integrity_check").Scan(&integrity); err != nil {
		return fmt.Errorf("备份数据库完整性校验失败: %w", err)
	}
	if strings.TrimSpace(integrity) != "ok" {
		return fmt.Errorf("备份数据库已损坏（integrity_check: %s）", integrity)
	}

	violations, err := countForeignKeyViolations(sqlDB)
	if err != nil {
		return fmt.Errorf("外键校验失败: %w", err)
	}
	if violations > 0 {
		return fmt.Errorf("备份数据库存在 %d 处外键违例", violations)
	}
	return nil
}

func countForeignKeyViolations(sqlDB *sql.DB) (int, error) {
	rows, err := sqlDB.Query("PRAGMA foreign_key_check")
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	count := 0
	for rows.Next() {
		count++
	}
	return count, rows.Err()
}

// RestoreSqliteSnap 用解包出的干净快照替换当前 SQLite 库：
// 先对当前库做安全快照（回滚点），关闭连接，把快照原子 rename 覆盖，
// 再重新打开并迁移。仅 sqlite 驱动可用；恢复期间服务短暂无写能力，
// 生产环境建议先开启站点维护模式。
func RestoreSqliteSnap(snapPath string) error {
	cfg := GetConf()
	if cfg == nil || cfg.Database.Driver != config.DriverSQLite {
		return fmt.Errorf("仅 SQLite 驱动支持恢复")
	}
	target, err := parseSqliteFileFromDsn(cfg.Database.Dsn)
	if err != nil {
		return err
	}

	// 回滚点：先备一份当前库（失败不影响，仅作为保险）。
	backupDir := cfg.Backup.Dir
	if strings.TrimSpace(backupDir) == "" {
		backupDir = "backup"
	}
	if _, backupErr := BackupSqlite(backupDir, false); backupErr != nil {
		logBackupError("restore: pre-restore snapshot failed", backupErr)
	}

	// 关闭当前连接
	if Db != nil {
		if err := closeDatabase(Db); err != nil {
			return fmt.Errorf("关闭当前数据库失败: %w", err)
		}
	}

	dir := filepath.Dir(target)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	tmp := filepath.Join(dir, ".restore-"+fmt.Sprint(time.Now().UnixNano())+".db")
	if err := copyFile(snapPath, tmp); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("写入恢复文件失败: %w", err)
	}
	// 清除残留 WAL/SHM 边车，避免陈旧日志污染新库
	_ = os.Remove(target + "-wal")
	_ = os.Remove(target + "-shm")
	if err := os.Rename(tmp, target); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("替换数据库失败: %w", err)
	}

	// 重新打开并迁移
	reopened, err := openDatabaseWithConfig(cfg)
	if err != nil {
		return fmt.Errorf("恢复后重新打开数据库失败: %w", err)
	}
	Db = reopened
	if err := migrateDbTable(); err != nil {
		return fmt.Errorf("恢复后迁移表结构失败: %w", err)
	}
	return nil
}

// RestoreFilesFromExtract 把解包出的 files/ 合并进上传目录（缺哪补哪，不删除现有文件）。
func RestoreFilesFromExtract(srcFilesDir string) error {
	if srcFilesDir == "" {
		return nil
	}
	return copyTree(srcFilesDir, config.FilesPath)
}

func copyTree(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, rerr := filepath.Rel(src, path)
		if rerr != nil {
			return rerr
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		return copyFile(path, target)
	})
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}

// openReadonlySqlite 以只读模式打开一个 sqlite 文件并返回 gorm 实例（调用方负责关闭）。
func openReadonlySqlite(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file:"+filepath.ToSlash(dbPath)+"?_mode=ro&_busy_timeout=3000"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}
