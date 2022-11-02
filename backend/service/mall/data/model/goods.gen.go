// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNameGood = "goods"

// Good mapped from table <goods>
type Good struct {
	ID        int32          `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Img       string         `gorm:"column:img" json:"img"`
	GoodsName string         `gorm:"column:goods_name;not null" json:"goods_name"`
	CreatedAt time.Time      `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName Good's table name
func (*Good) TableName() string {
	return TableNameGood
}
