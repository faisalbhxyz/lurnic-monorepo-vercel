package models

import "time"

type PaymentMethod struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string    `json:"title" gorm:"type:varchar(100)"`
	Image       *string   `json:"image" gorm:"type:text;null"`
	Instruction string    `json:"instruction" gorm:"type:text"`
	Status      bool      `json:"status"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	TenantID    uint      `gorm:"column:tenant_id" json:"-"`
	Tenant      Tenant    `gorm:"foreignKey:TenantID;references:ID" json:"-"`
}
