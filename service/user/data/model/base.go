package model

import (
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"

	"fastduck/treasure-doc/service/user/gid"
)

// TimeFormat 时间格式
const TimeFormat = "2006-01-02 15:04:05"

type BaseModel struct {
	Id        int64          `json:"-"  gorm:"column:id;type:bigint(20) unsigned;primary_key;"`
	Gid       string         `json:"id"  gorm:"-"`
	UpdatedAt time.Time      `json:"createdAt" gorm:"column:updated_at;type:datetime"`
	CreatedAt time.Time      `json:"updatedAt" gorm:"column:created_at;type:datetime"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"column:deleted_at;type:datetime"`
}

func (m *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	if m.Id == 0 {
		m.Id = gid.GenId()
	}
	return
}

func (m *BaseModel) BeforeSave(tx *gorm.DB) (err error) {
	if m.Gid != "" {
		num, err := strconv.ParseInt(m.Gid, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to convert gid %s to int", m.Gid)
		}
		m.Id = num
	}
	return
}

func (m *BaseModel) AfterFind(tx *gorm.DB) (err error) {
	if m.Id > 0 {
		m.Gid = strconv.FormatInt(m.Id, 10)
	}
	return
}

func (m *BaseModel) AfterCreate(tx *gorm.DB) (err error) {
	if m.Id > 0 {
		m.Gid = strconv.FormatInt(m.Id, 10)
	}
	return
}

func (m *BaseModel) AfterSave(tx *gorm.DB) (err error) {
	if m.Id > 0 {
		m.Gid = strconv.FormatInt(m.Id, 10)
	}
	return
}
