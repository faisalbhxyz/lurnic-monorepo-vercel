package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    string    `gorm:"column:user_id;uniqueIndex" json:"user_id"`
	Name      string    `json:"name"`
	Phone     *string   `json:"phone"`
	Email     string    `gorm:"uniqueIndex" json:"email"`
	Password  string    `json:"-"`
	Status    bool      `json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	TenantID  uint      `gorm:"column:tenant_id" json:"-"`
	Tenant    Tenant    `gorm:"foreignKey:TenantID;references:ID" json:"-"`
	RoleID    *uint     `gorm:"column:role_id" json:"role_id"`
	Role      Role      `gorm:"foreignKey:RoleID;references:ID" json:"role"`
}

// get users without role

type UserWithoutRole struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    string    `gorm:"column:user_id;uniqueIndex" json:"user_id"`
	Name      string    `json:"name"`
	Phone     *string   `json:"phone"`
	Email     string    `gorm:"uniqueIndex" json:"email"`
	Status    bool      `json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	// TenantID  uint      `gorm:"column:tenant_id" json:"-"`
	// Tenant    Tenant    `gorm:"foreignKey:TenantID;references:ID" json:"-"`
	// RoleID    *uint     `gorm:"column:role_id" json:"-"`
	// Role      Role      `gorm:"foreignKey:RoleID;references:ID" json:"-"`
}

func (UserWithoutRole) TableName() string {
	return "users"
}
