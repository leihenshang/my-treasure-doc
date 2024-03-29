// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNameGoodsSku = "goods_sku"

// GoodsSku mapped from table <goods_sku>
type GoodsSku struct {
	ID           int32          `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Enabled      int32          `gorm:"column:enabled;not null;default:1" json:"enabled"`     // 1-可用，2-禁用
	GoodsID      int32          `gorm:"column:goods_id;not null" json:"goods_id"`             // 商品id
	GoodsSpecIds string         `gorm:"column:goods_spec_ids;not null" json:"goods_spec_ids"` // 规格id
	Price        float64        `gorm:"column:price;not null;default:0.0000" json:"price"`    // 价格
	Stock        int32          `gorm:"column:stock;not null" json:"stock"`                   // 库存
	CreatedAt    time.Time      `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName GoodsSku's table name
func (*GoodsSku) TableName() string {
	return TableNameGoodsSku
}
