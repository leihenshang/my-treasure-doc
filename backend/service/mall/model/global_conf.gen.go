// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNameGlobalConf = "global_conf"

// GlobalConf mapped from table <global_conf>
type GlobalConf struct {
	ID        int64          `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Key       string         `gorm:"column:key;not null" json:"key"`
	Value     string         `gorm:"column:value;not null" json:"value"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	Version   int32          `gorm:"column:version;not null" json:"version"`
	CreatedBy int64          `gorm:"column:created_by;not null" json:"created_by"`
}

// TableName GlobalConf's table name
func (*GlobalConf) TableName() string {
	return TableNameGlobalConf
}
