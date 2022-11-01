package model

import (
	"time"
)

type Goods struct {
	Id        uint      `gorm:"column:id;type:int(10) unsigned;primary_key;AUTO_INCREMENT" json:"id"`
	GoodsName string    `gorm:"column:goods_name;type:varchar(100);NOT NULL" json:"goods_name"`
	Quantity  uint      `gorm:"column:quantity;type:int(10) unsigned;default:0;NOT NULL" json:"quantity"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"updated_at"`
	DeletedAt time.Time `gorm:"column:deleted_at;type:datetime" json:"deleted_at"`
}
