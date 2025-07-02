package models

import "time"

type Banner struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title     *string   `json:"title" form:"title" gorm:"type:varchar(100);null"`
	Url       *string   `json:"url" form:"url" gorm:"type:varchar(255);null"`
	Image     string    `json:"image" form:"image"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	TenantID  uint      `gorm:"column:tenant_id" json:"-"`
	Tenant    Tenant    `gorm:"foreignKey:TenantID;references:ID" json:"-"`
}
