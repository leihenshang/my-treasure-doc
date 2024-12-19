package model

import (
	"time"
)

type GlobalConf struct {
	ID        string    `gorm:"column:id;primary_key"`
	Key       string    `gorm:"column:key;NOT NULL"`
	Value     string    `gorm:"column:value;NOT NULL"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	DeletedAt time.Time `gorm:"column:deleted_at"`
	Version   int       `gorm:"column:version;NOT NULL"`
	CreatedBy int64     `gorm:"column:created_by;NOT NULL"`
}

func (m *GlobalConf) TableName() string {
	return "td_global_conf"
}
