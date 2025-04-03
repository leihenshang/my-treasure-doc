package model

import (
	"time"
)

type TeamUser struct {
	ID        string    `gorm:"column:id;primary_key"`
	UserID    string    `gorm:"column:user_id;NOT NULL"`
	TeamID    string    `gorm:"column:team_id;NOT NULL"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	DeletedAt time.Time `gorm:"column:deleted_at"`
}

func (m *TeamUser) TableName() string {
	return "td_team_user"
}
