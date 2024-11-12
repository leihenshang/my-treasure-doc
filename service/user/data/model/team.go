package model

import (
	"time"
)

type Team struct {
	ID        int64     `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	Name      string    `gorm:"column:name;NOT NULL"`             // 名字
	Number    int       `gorm:"column:number;default:0;NOT NULL"` // 人数
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	DeletedAt time.Time `gorm:"column:deleted_at"`
}

func (m *Team) TableName() string {
	return "td_team"
}
