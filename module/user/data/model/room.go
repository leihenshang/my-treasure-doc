package model

type RoomStatus string

const (
	RoomStatusNormal RoomStatus = "normal" // 正常
	RoomStatusClosed RoomStatus = "closed" // 关闭
)

type Room struct {
	BaseModel
	Name      string     `json:"name" gorm:"column:name;NOT NULL"`
	UserId    string     `json:"userId" gorm:"column:user_id;NOT NULL"` // 房主
	Status    RoomStatus `json:"status" gorm:"column:status;NOT NULL;default:'normal'"`
	IsDefault int8       `json:"isDefault" gorm:"column:is_default;NOT NULL;default:0"`
}

func (m *Room) TableName() string {
	return "td_room"
}
