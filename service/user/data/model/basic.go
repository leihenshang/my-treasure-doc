package model

import (
	"gorm.io/gorm"
	"time"
)

// TimeFormat 时间格式
const TimeFormat = "2006-01-02 15:04:05"

type BasicModel struct {
	Id        int64          `json:"id"  gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT"`
	UpdatedAt time.Time      `json:"createdAt" gorm:"column:updated_at;type:datetime"`
	CreatedAt time.Time      `json:"updatedAt" gorm:"column:created_at;type:datetime"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"column:deleted_at;type:datetime"`
}
