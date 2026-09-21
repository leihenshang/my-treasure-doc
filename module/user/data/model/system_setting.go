package model

// SystemSetting 系统级键值对设置（如 NAS 备份令牌）。
// 这类配置只供服务端 / 后台读取，不属于公开博客内容，因此独立成表、
// 绝不并入 blogresponse.Site，避免被公开站点接口泄露。
type SystemSetting struct {
	BaseModel
	// Key 设置键，如 backup.apiToken。
	Key string `gorm:"column:setting_key;size:128;not null;uniqueIndex"`
	// Value 设置值。空字符串表示该键已被显式置空（可用于"清除令牌"）。
	Value string `gorm:"column:setting_value;type:text"`
}

func (SystemSetting) TableName() string { return "td_system_setting" }
