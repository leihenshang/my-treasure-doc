// Code generated by sql2gorm. DO NOT EDIT.
package model

import (
	"time"
)

// 团队成员表
type TeamUser struct {
	ID        uint64    `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	UserID    uint64    `gorm:"column:user_id;NOT NULL"`
	TeamID    uint64    `gorm:"column:team_id;NOT NULL"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	DeletedAt time.Time `gorm:"column:deleted_at"`
}

func (m *TeamUser) TableName() string {
	return "team_user"
}


