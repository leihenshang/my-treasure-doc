package global

import (
	"errors"
	"strings"

	"fastduck/treasure-doc/module/user/data/model"

	"gorm.io/gorm"
)

// 系统级设置键名。
const SystemSettingBackupTokenKey = "backup.apiToken"

// GetSystemSetting 读取系统设置，存在返回 (value, true)。数据表未迁移/查询失败时返回错误。
func GetSystemSetting(key string) (string, bool, error) {
	if Db == nil {
		return "", false, nil
	}
	var setting model.SystemSetting
	err := Db.Where("setting_key = ?", key).First(&setting).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return setting.Value, true, nil
}

// SetSystemSetting 写入（或覆盖）系统设置。空 value 表示显式置空该键。
func SetSystemSetting(key, value string) error {
	if Db == nil {
		return nil
	}
	var setting model.SystemSetting
	err := Db.Where("setting_key = ?", key).First(&setting).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		setting = model.SystemSetting{Key: key, Value: value}
		return Db.Create(&setting).Error
	}
	if err != nil {
		return err
	}
	setting.Value = value
	return Db.Save(&setting).Error
}

// EffectiveBackupApiToken 返回当前生效的 NAS 备份令牌：
// 后台（DB）优先 —— 已显式存储则以 DB 值为准（空值 = 禁用 NAS）；
// 后台从未配过则回退到配置文件的 [backup].apiToken，保持零配置兼容。
func EffectiveBackupApiToken() string {
	value, exists, err := GetSystemSetting(SystemSettingBackupTokenKey)
	if err == nil && exists {
		return strings.TrimSpace(value)
	}
	cfg := GetConf()
	if cfg != nil {
		return strings.TrimSpace(cfg.Backup.ApiToken)
	}
	return ""
}
