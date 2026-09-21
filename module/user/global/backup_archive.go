package global

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fastduck/treasure-doc/module/user/config"
)

// BackupSchemaVersion 当前备份包 schema 版本。导入恢复时要求与备份一致，
// 避免用更旧/更新的结构覆盖运行中的库。随着表结构演进需要递增并在导入端兼容。
const BackupSchemaVersion = 1

// ManiftFile 描述归档中的一个文件条目，用于导入时的完整性校验（MD5/size）。
type ManifestFile struct {
	Path   string `json:"path"`
	Size   int64  `json:"size"`
	SHA256 string `json:"sha256"`
}

// Manifest 完整备份包的元数据描述。
type Manifest struct {
	AppVersion    string           `json:"appVersion"`
	SchemaVersion int              `json:"schemaVersion"`
	Driver        string           `json:"driver"`
	CreatedAt     string           `json:"createdAt"`
	IncludeAuth   bool             `json:"includeAuth"`
	Stats         map[string]int64 `json:"stats,omitempty"`
	Files         []ManifestFile   `json:"files"`
}

const (
	archiveDBName       = "treasure_doc.db"
	archiveManifestName = "manifest.json"
	archiveFileRoot     = "files"
)

// BuildBackupArchive 生成一份完整备份包（treasure_doc.db + files/ + manifest.json 的 tar.gz）
// 到 dstDir，文件名形如 treasure-backup-<yyyyMMdd>-<HHmmss>.tar.gz，返回完整路径。
// 先对当前库做 VACUUM INTO 快照（保证一致性），再连同上传目录一起打包。
// 仅 sqlite 驱动支持。includeAuth 目前恒为 true：本系统是单机个人博客，
// 恢复也需管理员授权，整库快照最正确，不额外剥离鉴权表。
func BuildBackupArchive(dstDir string) (string, error) {
	cfg := GetConf()
	if cfg == nil || cfg.Database.Driver != config.DriverSQLite {
		return "", fmt.Errorf("仅 SQLite 驱动支持备份")
	}

	if dstDir == "" {
		dstDir = "backup"
	}
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		return "", fmt.Errorf("创建备份目录失败: %w", err)
	}

	tmp, err := os.MkdirTemp("", "treasure-archive-*")
	if err != nil {
		return "", fmt.Errorf("创建临时目录失败: %w", err)
	}
	defer os.RemoveAll(tmp)

	// 1. 数据库一致性快照
	snapDir := filepath.Join(tmp, "snap")
	if err := os.MkdirAll(snapDir, 0o755); err != nil {
		return "", err
	}
	snapPath, err := BackupSqlite(snapDir, false)
	if err != nil {
		return "", err
	}

	// 2. 收集上传目录文件清单（relative path 以 files/ 开头）
	files, listErr := collectManifestFiles()
	if listErr != nil {
		return "", listErr
	}

	manifest := Manifest{
		AppVersion:    cfg.App.Name,
		SchemaVersion: BackupSchemaVersion,
		Driver:        cfg.Database.Driver,
		CreatedAt:     time.Now().Format(time.RFC3339),
		IncludeAuth:   true,
		Files:         files,
	}
	if stats := tableRowCounts(snapPath); len(stats) > 0 {
		manifest.Stats = stats
	}

	// 3. 打包 tar.gz
	name := fmt.Sprintf("treasure-backup-%s.tar.gz", time.Now().Format("20060102-150405"))
	dst := filepath.Join(dstDir, name)
	if err := writeArchive(dst, tmp, snapPath, manifest); err != nil {
		os.Remove(dst)
		return "", err
	}
	return dst, nil
}

// collectManifestFiles 递归遍历上传目录 files/，返回条目清单（相对路径以 files/ 开头）。
func collectManifestFiles() ([]ManifestFile, error) {
	var files []ManifestFile
	root := config.FilesPath
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		sum, size, err := sha256File(path)
		if err != nil {
			return err
		}
		files = append(files, ManifestFile{
			Path:   archiveFileRoot + "/" + filepath.ToSlash(rel),
			Size:   size,
			SHA256: sum,
		})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("扫描上传目录失败: %w", err)
	}
	return files, nil
}

// tableRowCounts 在快照库上统计各业务表行数（尽力而为，失败不影响打包）。
func tableRowCounts(dbPath string) map[string]int64 {
	snap, err := openReadonlySqlite(dbPath)
	if err != nil {
		return nil
	}
	defer func() {
		if sqlDB, cerr := snap.DB(); cerr == nil {
			_ = sqlDB.Close()
		}
	}()

	var tables []string
	if err := snap.Raw("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'").Scan(&tables).Error; err != nil {
		return nil
	}

	stats := make(map[string]int64)
	for _, table := range tables {
		var count int64
		if err := snap.Raw("SELECT COUNT(*) FROM `" + table + "`").Scan(&count).Error; err != nil {
			continue
		}
		stats[table] = count
	}
	return stats
}

// writeArchive 把快照库（重命名为 treasure_doc.db）、上传目录与 manifest 写入 dst tar.gz。
func writeArchive(dst, tmpDir, snapPath string, manifest Manifest) error {
	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("创建归档文件失败: %w", err)
	}
	defer out.Close()

	gz := gzip.NewWriter(out)
	defer gz.Close()
	tw := tar.NewWriter(gz)
	defer tw.Close()

	// manifest.json
	if err := writeTarBytes(tw, archiveManifestName, manifestBytes(manifest)); err != nil {
		return err
	}
	// treasure_doc.db
	if err := writeTarFile(tw, archiveDBName, snapPath); err != nil {
		return err
	}
	// files/**
	return writeTarTree(tw, config.FilesPath, archiveFileRoot)
}

func manifestBytes(manifest Manifest) []byte {
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return []byte("{}")
	}
	return data
}

// writeTarFile 把本地文件作为单个 archive/name 条目写入。
func writeTarFile(tw *tar.Writer, name, src string) error {
	fi, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("读取待归档文件失败: %w", err)
	}
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()
	hdr := &tar.Header{
		Name:    name,
		Mode:    0o644,
		Size:    fi.Size(),
		ModTime: fi.ModTime(),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return err
	}
	if _, err := io.Copy(tw, f); err != nil {
		return err
	}
	return nil
}

// writeTarBytes 写入一个内存字节条目。
func writeTarBytes(tw *tar.Writer, name string, data []byte) error {
	hdr := &tar.Header{Name: name, Mode: 0o640, Size: int64(len(data)), ModTime: time.Now()}
	if err := tw.WriteHeader(hdr); err != nil {
		return err
	}
	_, err := tw.Write(data)
	return err
}

// writeTarTree 递归把本地目录 root 写入归档，条目前缀为 prefix（如 files/）。
func writeTarTree(tw *tar.Writer, root, prefix string) error {
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		rel, relErr := filepath.Rel(root, path)
		if relErr != nil {
			return relErr
		}
		name := prefix + "/" + filepath.ToSlash(rel)
		if d.IsDir() {
			// 目录也写条目（TarEntry 为空目录时保留结构）
			if rel == "." {
				return nil
			}
			info, ierr := d.Info()
			if ierr != nil {
				return ierr
			}
			return tw.WriteHeader(&tar.Header{Name: name + "/", Mode: 0o755, Typeflag: tar.TypeDir, ModTime: info.ModTime()})
		}
		return writeTarFile(tw, name, path)
	})
	if err != nil {
		return fmt.Errorf("打包 files/ 目录失败: %w", err)
	}
	return nil
}

// sha256File 计算文件 sha256 与大小。
func sha256File(path string) (string, int64, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", 0, err
	}
	defer f.Close()
	h := sha256.New()
	size, err := io.Copy(h, f)
	if err != nil {
		return "", 0, err
	}
	return hex.EncodeToString(h.Sum(nil)), size, nil
}

// backupArchiveLinkName 校验归档条目名是否安全（拒绝绝对路径与 ../ 逃逸），供导入解包使用。
func backupArchiveLinkSafe(name string) bool {
	clean := filepath.ToSlash(strings.TrimPrefix(name, "./"))
	if clean == "" || clean == "." || strings.HasPrefix(clean, "/") {
		return false
	}
	if strings.Contains(clean, "..") {
		for _, part := range strings.Split(clean, "/") {
			if part == ".." {
				return false
			}
		}
	}
	return true
}
