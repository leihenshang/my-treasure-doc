package model

type RoomStatus string

const (
	RoomStatusNormal RoomStatus = "normal" // 正常
	RoomStatusClosed RoomStatus = "closed" // 关闭
)

type Room struct {
	BaseModel
	Name   string     `gorm:"column:name;NOT NULL"`
	UserId string     `gorm:"column:user_id;NOT NULL"` // 房主
	Status RoomStatus `gorm:"column:status;NOT NULL;default:'normal'"`
}

func (m *Room) TableName() string {
	return "td_room"
}
