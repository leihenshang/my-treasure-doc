package model

import (
	"time"

	"gorm.io/gorm"

	"fastduck/treasure-doc/module/user/gid"
)

// TimeFormat 时间格式
const TimeFormat = "2006-01-02 15:04:05"

type BaseModel struct {
	Id        string         `json:"id"  gorm:"column:id;type:varchar(100);primary_key;"`
	UpdatedAt time.Time      `json:"createdAt" gorm:"column:updated_at;type:datetime"`
	CreatedAt time.Time      `json:"updatedAt" gorm:"column:created_at;type:datetime"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"column:deleted_at;type:datetime"`
}

func (m *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	if m.Id == "" {
		m.Id = gid.GenId()
	}
	return
}
