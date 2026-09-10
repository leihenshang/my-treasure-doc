package api

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"fastduck/treasure-doc/module/user/config"
	"fastduck/treasure-doc/module/user/data/response"
	"fastduck/treasure-doc/module/user/global"

	"github.com/gin-gonic/gin"
)

type backupFile struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	ModTime string `json:"modTime"`
}

// BackupApi 暴露 SQLite 备份的查看、触发与下载，挂在管理端路由下（需管理员）。
type BackupApi struct{}

func NewBackupApi() *BackupApi { return &BackupApi{} }

// ListBackups 列出备份目录中的文件，按修改时间倒序。
func (b *BackupApi) ListBackups(c *gin.Context) {
	dir, ok := sqliteBackupDir(c)
	if !ok {
		return
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			response.OkWithData(c, []backupFile{})
			return
		}
		response.FailWithMessage(c, "读取备份目录失败")
		return
	}

	files := make([]backupFile, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		files = append(files, backupFile{Name: entry.Name(), Size: info.Size(), ModTime: info.ModTime().Format(time.RFC3339)})
	}
	sort.Slice(files, func(i, j int) bool { return files[i].ModTime > files[j].ModTime })
	response.OkWithData(c, files)
}

// CreateBackup 立即执行一次备份并返回文件信息。
func (b *BackupApi) CreateBackup(c *gin.Context) {
	cfg := global.GetConf()
	if cfg == nil || cfg.Database.Driver != config.DriverSQLite {
		response.FailWithMessage(c, "仅 SQLite 驱动支持备份")
		return
	}
	path, err := global.BackupSqlite(backupDirectory(cfg.Backup.Dir), cfg.Backup.Compress)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}

	file := backupFile{Name: filepath.Base(path)}
	if info, err := os.Stat(path); err == nil {
		file.Size = info.Size()
		file.ModTime = info.ModTime().Format(time.RFC3339)
	}
	response.OkWithData(c, file)
}

// DownloadBackup 下载指定备份文件。
func (b *BackupApi) DownloadBackup(c *gin.Context) {
	dir, ok := sqliteBackupDir(c)
	if !ok {
		return
	}
	name := c.Param("name")
	// 只允许备份目录下的普通文件名，阻断 ../ 越权
	if name != filepath.Base(name) || name == "." || name == ".." {
		response.FailWithMessage(c, "备份文件名不合法")
		return
	}
	path := filepath.Join(dir, name)
	if info, err := os.Stat(path); err != nil || info.IsDir() {
		response.FailWithMessage(c, "备份文件不存在")
		return
	}

	c.Header("Content-Disposition", `attachment; filename="`+name+`"`)
	c.File(path)
}

func sqliteBackupDir(c *gin.Context) (string, bool) {
	cfg := global.GetConf()
	if cfg == nil || cfg.Database.Driver != config.DriverSQLite {
		response.FailWithMessage(c, "仅 SQLite 驱动支持备份")
		return "", false
	}
	return backupDirectory(cfg.Backup.Dir), true
}

func backupDirectory(dir string) string {
	if strings.TrimSpace(dir) == "" {
		return "backup"
	}
	return dir
}
