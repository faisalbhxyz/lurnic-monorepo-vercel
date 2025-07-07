package models

import (
	"time"
)

type Role struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	TenantID  *uint     `gorm:"column:tenant_id" json:"tenant_id"`
	Tenant    Tenant    `gorm:"foreignKey:TenantID;references:ID" json:"-"`
	Users     []User    `gorm:"foreignKey:RoleID;references:ID" json:"user"`
}
