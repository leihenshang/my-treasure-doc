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

// maxRestoreSize 导入备份包的大小上限（1GB）。
const maxRestoreSize = 1 << 30

// BackupTokenRequest 更新 NAS 备份令牌的请求体。
type BackupTokenRequest struct {
	// Token 新令牌；空字符串表示清除（禁用 NAS 导出）。
	Token string `json:"token"`
}

// GetBackupToken 返回当前生效的 NAS 备份令牌（后台 DB 优先，无配置时回退配置文件）。
func (b *BackupApi) GetBackupToken(c *gin.Context) {
	response.OkWithData(c, gin.H{"token": global.EffectiveBackupApiToken()})
}

// PutBackupToken 更新 NAS 备份令牌（仅管理员）。写入 DB 后即时生效，无需重启。
func (b *BackupApi) PutBackupToken(c *gin.Context) {
	var req BackupTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(c, "参数错误")
		return
	}
	if err := global.SetSystemSetting(global.SystemSettingBackupTokenKey, req.Token); err != nil {
		response.FailWithMessage(c, "保存失败")
		return
	}
	response.OkWithData(c, gin.H{"token": req.Token})
}

// ExportArchive 后台手动导出「完整备份包」（db + files/ + manifest），并存放进备份目录，
// 使其同时出现在备份列表中、可用已有下载接口取回。返回文件信息。
func (b *BackupApi) ExportArchive(c *gin.Context) {
	dir, ok := sqliteBackupDir(c)
	if !ok {
		return
	}
	path, err := global.BuildBackupArchive(dir)
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

// Restore 后台导入：上传 .tar.gz 备份包，先解包校验再恢复数据库与图片。
// 仅管理员可操作；校验任一失败即中止、不写任何数据。
func (b *BackupApi) Restore(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.FailWithMessage(c, "请选择备份文件(.tar.gz)")
		return
	}
	if file.Size <= 0 {
		response.FailWithMessage(c, "备份文件为空")
		return
	}
	if file.Size > maxRestoreSize {
		response.FailWithMessage(c, "备份文件过大")
		return
	}

	upload, err := os.CreateTemp("", "restore-upload-*.tar.gz")
	if err != nil {
		response.FailWithMessage(c, "创建临时文件失败")
		return
	}
	uploadPath := upload.Name()
	upload.Close()
	defer os.Remove(uploadPath)
	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		response.FailWithMessage(c, "保存上传文件失败")
		return
	}

	extractDir, err := os.MkdirTemp("", "restore-extract-*")
	if err != nil {
		response.FailWithMessage(c, "创建解包目录失败")
		return
	}
	defer os.RemoveAll(extractDir)

	archived, err := global.ExtractBackupArchive(uploadPath, extractDir)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	if err := global.ValidateBackupSchema(archived.Manifest); err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	if err := global.ValidateSqliteFile(archived.DBPath); err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}

	// 恢复数据库（含回滚保护），再合并图片
	if err := global.RestoreSqliteSnap(archived.DBPath); err != nil {
		response.FailWithMessage(c, "恢复数据库失败："+err.Error())
		return
	}
	if err := global.RestoreFilesFromExtract(archived.FilesRoot); err != nil {
		response.FailWithMessage(c, "恢复图片失败："+err.Error())
		return
	}

	response.OkWithData(c, gin.H{"restored": true, "createdAt": archived.Manifest.CreatedAt})
}

// NasExportArchive NAS 机器令牌下载完整备份：现场打包 tar.gz 并流式下载。
// 鉴权由中间件 BackupToken() 完成；该端点只读，不提供导入/删除能力。
func (b *BackupApi) NasExportArchive(c *gin.Context) {
	dir, ok := sqliteBackupDir(c)
	if !ok {
		return
	}
	tmpDir, err := os.MkdirTemp(dir, ".nas-*")
	if err != nil {
		response.FailWithMessage(c, "创建临时目录失败")
		return
	}
	defer os.RemoveAll(tmpDir)

	path, err := global.BuildBackupArchive(tmpDir)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	name := filepath.Base(path)
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
