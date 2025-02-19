// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNameMerchantsMetum = "merchants_meta"

// MerchantsMetum mapped from table <merchants_meta>
type MerchantsMetum struct {
	ID        int64          `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	MID       int64          `gorm:"column:m_id;not null" json:"m_id"`
	Ac        string         `gorm:"column:ac;comment:钱包地址" json:"ac"` // 钱包地址
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName MerchantsMetum's table name
func (*MerchantsMetum) TableName() string {
	return TableNameMerchantsMetum
}
