package models

import "time"

type Category struct {
	ID            uint          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string        `json:"name"`
	Slug          string        `gorm:"column:slug" json:"slug"`
	Description   *string       `gorm:"type:text" json:"description"`
	Thumbnail     *string       `json:"thumbnail"`
	CreatedAt     time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	TenantID      uint          `gorm:"column:tenant_id" json:"-"`
	Tenant        Tenant        `gorm:"foreignKey:TenantID;references:ID" json:"-"`
	SubCategories []SubCategory `gorm:"foreignKey:CategoryID;references:ID" json:"sub_categories"`
}
