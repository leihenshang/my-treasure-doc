package model

import (
	"time"
)

type GlobalConf struct {
	ID        int64     `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	Key       string    `gorm:"column:key;NOT NULL"`
	Value     string    `gorm:"column:value;NOT NULL"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	DeletedAt time.Time `gorm:"column:deleted_at"`
	Version   int       `gorm:"column:version;default:0;NOT NULL"`
	CreatedBy int64     `gorm:"column:created_by;default:0;NOT NULL"`
}

func (m *GlobalConf) TableName() string {
	return "td_global_conf"
}
