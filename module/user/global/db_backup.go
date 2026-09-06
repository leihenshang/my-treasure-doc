package global

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fastduck/treasure-doc/module/user/config"
)

// parseSqliteFileFromDsn 从 sqlite 的 DSN 中提取数据库文件路径。
// 形如 file:./treasure_doc.db?_busy_timeout=5000&_journal_mode=WAL
func parseSqliteFileFromDsn(dsn string) (string, error) {
	const prefix = "file:"
	if !strings.HasPrefix(dsn, prefix) {
		return "", fmt.Errorf("unsupported sqlite dsn: %q", dsn)
	}
	rest := dsn[len(prefix):]
	// 去掉查询参数（? 之后）。
	if idx := strings.IndexByte(rest, '?'); idx >= 0 {
		rest = rest[:idx]
	}
	rest = strings.TrimSpace(rest)
	if rest == "" {
		return "", fmt.Errorf("empty sqlite db file in dsn: %q", dsn)
	}
	return rest, nil
}

// escapeSqliteLiteral 转义 VACUUM INTO 语句里路径中的单引号，避免语法错误。
func escapeSqliteLiteral(p string) string {
	return strings.ReplaceAll(p, "'", "''")
}

// BackupSqlite 使用 VACUUM INTO 对当前 SQLite 数据库做一次在线热备（无需停服）。
// 当驱动不是 sqlite 时直接返回空字符串与 nil（调用方应跳过）。
// compress=true 时先 VACUUM 到临时 db，再 gzip 压缩为目标 .db.gz。
func BackupSqlite(backupDir string, compress bool) (string, error) {
	cfg := GetConf()
	if cfg == nil || cfg.Database.Driver != config.DriverSQLite {
		return "", nil
	}
	src, err := parseSqliteFileFromDsn(cfg.Database.Dsn)
	if err != nil {
		return "", err
	}

	if backupDir == "" {
		backupDir = "backup"
	}
	if err := os.MkdirAll(backupDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create backup dir: %w", err)
	}

	base := filepath.Base(src)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	ts := time.Now().Format("20060102_150405")
	tmp := filepath.Join(backupDir, fmt.Sprintf("%s_%s.db", name, ts))

	// VACUUM INTO 生成一个干净的独立数据库文件（已落盘全部 WAL 数据）。
	if err := Db.Exec(fmt.Sprintf("VACUUM INTO '%s'", escapeSqliteLiteral(tmp))).Error; err != nil {
		os.Remove(tmp)
		return "", fmt.Errorf("vacuum into failed: %w", err)
	}

	if compress {
		dst := tmp + ".gz"
		if err := gzipFile(tmp, dst); err != nil {
			os.Remove(tmp)
			return "", err
		}
		os.Remove(tmp)
		return dst, nil
	}
	return tmp, nil
}

func gzipFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	gw := gzip.NewWriter(out)
	defer gw.Close()

	if _, err := io.Copy(gw, in); err != nil {
		return err
	}
	return gw.Close()
}

// StartSqliteBackupScheduler 启动按 interval 周期执行的 SQLite 定时备份。
// keepDays>0 时清理超过该天数的旧备份；非 sqlite 驱动不启动，返回空停止函数。
func StartSqliteBackupScheduler(interval time.Duration, backupDir string, compress bool, keepDays int) func() {
	cfg := GetConf()
	if cfg == nil || cfg.Database.Driver != config.DriverSQLite {
		return func() {}
	}

	stop := make(chan struct{})
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-stop:
				ticker.Stop()
				return
			case <-ticker.C:
				path, err := BackupSqlite(backupDir, compress)
				if err != nil {
					logBackupError("sqlite backup failed", err)
					continue
				}
				if Log != nil {
					Log.Infof("sqlite backup done: %s", path)
				} else {
					fmt.Printf("sqlite backup done: %s\n", path)
				}
				if keepDays > 0 {
					if err := cleanOldBackups(backupDir, keepDays); err != nil {
						logBackupError("clean old backups failed", err)
					}
				}
			}
		}
	}()
	return func() { close(stop) }
}

func logBackupError(msg string, err error) {
	if Log != nil {
		Log.Errorf("%s: %v", msg, err)
	} else {
		fmt.Printf("%s: %v\n", msg, err)
	}
}

// cleanOldBackups 删除 backupDir 中修改时间早于 keepDays 天前的备份文件。
func cleanOldBackups(dir string, keepDays int) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	cutoff := time.Now().AddDate(0, 0, -keepDays)
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		if info.ModTime().Before(cutoff) {
			if err := os.Remove(filepath.Join(dir, e.Name())); err != nil {
				return err
			}
		}
	}
	return nil
}
