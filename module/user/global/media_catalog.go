package global

import (
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"

	"fastduck/treasure-doc/module/blog/data/model"
	"fastduck/treasure-doc/module/user/config"

	"gorm.io/gorm"
)

// UpsertMediaMeta 写入/更新某上传文件的元数据记录（不存在则新增，存在则刷新）。
// 文件本体是否已落盘由调用方保证；这里只维护 files/blog 的目录登记。
func UpsertMediaMeta(name, sha256, ext, mime, originName string, size int64) error {
	if Db == nil || name == "" {
		return nil
	}
	path := filepath.Join(config.FilesPath, "blog", name)
	width, height := 0, 0
	if w, h, err := probeImageSize(path); err == nil {
		width, height = w, h
	}

	var existing model.Media
	err := Db.Where("name = ?", name).First(&existing).Error
	record := model.Media{
		Name:       name,
		Size:       size,
		Ext:        ext,
		Mime:       mime,
		Sha256:     sha256,
		OriginName: originName,
		Width:      width,
		Height:     height,
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return Db.Create(&record).Error
	}
	if err != nil {
		return err
	}
	// 保留首次上传时间（created_at）不变，只刷新可变的字段
	record.CreatedAt = existing.CreatedAt
	return Db.Model(&existing).Select("size", "ext", "mime", "sha256", "origin_name", "width", "height").Updates(record).Error
}

// RemoveMediaMeta 删除某文件的元数据记录（配合物理删除调用）。
func RemoveMediaMeta(name string) error {
	if Db == nil || name == "" {
		return nil
	}
	return Db.Where("name = ?", name).Delete(&model.Media{}).Error
}

// SyncMediaCatalog 将 files/blog 目录与元数据表做一次对齐：
// 磁盘上存在而库中没有 → 补齐登记；库中存在而磁盘已删 → 清除记录。
// 用于流程外的文件变化（手放文件、恢复备份后）也能被媒体库识别。
func SyncMediaCatalog() error {
	if Db == nil {
		return nil
	}
	dir := filepath.Join(config.FilesPath, "blog")

	// 磁盘上的文件集合
	diskNames := map[string]struct{}{}
	if entries, err := os.ReadDir(dir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			diskNames[entry.Name()] = struct{}{}
		}
	}

	// 磁盘有而库没有 → 新增
	for name := range diskNames {
		info, err := os.Stat(filepath.Join(dir, name))
		if err != nil || info.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(name))
		sha := strings.TrimSuffix(name, ext)
		if err := UpsertMediaMeta(name, sha, ext, mimeByExt(ext), "", info.Size()); err != nil {
			return err
		}
	}

	// 库有而磁盘没有 → 删除记录
	var rows []model.Media
	if err := Db.Find(&rows).Error; err != nil {
		return err
	}
	for _, row := range rows {
		if _, ok := diskNames[row.Name]; ok {
			continue
		}
		if err := RemoveMediaMeta(row.Name); err != nil {
			return err
		}
	}
	return nil
}

// probeImageSize 尽力而为地读取图片宽高；非图片或不支持格式返回错误。
func probeImageSize(path string) (int, int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()
	cfg, _, err := image.DecodeConfig(f)
	if err != nil {
		return 0, 0, err
	}
	return cfg.Width, cfg.Height, nil
}

// mimeByExt 按后缀兜底推断 MIME（与上传侧白名单一致）。
func mimeByExt(ext string) string {
	return map[string]string{
		".jpg": "image/jpeg", ".jpeg": "image/jpeg", ".png": "image/png",
		".gif": "image/gif", ".webp": "image/webp", ".bmp": "image/bmp",
		".mp4": "video/mp4", ".m4v": "video/mp4", ".webm": "video/webm",
		".mov": "video/quicktime", ".avi": "video/avi", ".mkv": "video/x-matroska",
		".mpeg": "video/mpeg", ".mpg": "video/mpeg",
	}[ext]
}
