package models

import (
	"time"
)

type Instructor struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    string    `gorm:"column:user_id;uniqueIndex" json:"user_id"`
	FirstName string    `json:"first_name"`
	LastName  *string   `json:"last_name"`
	Phone     *string   `json:"phone"`
	Email     string    `gorm:"uniqueIndex" json:"email"`
	Password  string    `json:"-"`
	Status    bool      `json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	TenantID  uint      `gorm:"column:tenant_id" json:"-"`
	Tenant    Tenant    `gorm:"foreignKey:TenantID;references:ID" json:"Tenant"`
}
