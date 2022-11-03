// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNameOrder = "order"

// Order mapped from table <order>
type Order struct {
	ID        int32          `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	OrderNo   string         `gorm:"column:order_no;not null" json:"order_no"`            // 订单号
	UserID    int32          `gorm:"column:user_id;not null" json:"user_id"`              // 用户id
	Amount    float64        `gorm:"column:amount;not null;default:0.0000" json:"amount"` // 金额
	Status    int32          `gorm:"column:status;not null" json:"status"`                // 状态,0-异常,1-待支付,2-已支付,3-支付失败,4-用户取消,5-系统取消,6-订单异常
	CreatedAt time.Time      `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName Order's table name
func (*Order) TableName() string {
	return TableNameOrder
}
