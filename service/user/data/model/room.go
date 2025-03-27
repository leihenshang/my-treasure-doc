package model

import "time"

type RoomStatus string

const (
	RoomStatusNormal RoomStatus = "normal" // 正常
	RoomStatusClosed RoomStatus = "closed" // 关闭
)

type Room struct {
	ID        string     `gorm:"column:id;primary_key;NOT NULL"`
	Name      string     `gorm:"column:name;NOT NULL"`
	UserId    string     `gorm:"column:user_id;NOT NULL"` // 房主
	Status    RoomStatus `gorm:"column:status;NOT NULL;default:normal"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt time.Time  `gorm:"column:deleted_at"`
}

func (m *Room) TableName() string {
	return "td_room"
}
