package models

import (
	"time"
)

type Instructor struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      string    `gorm:"column:user_id;uniqueIndex" json:"user_id"`
	FirstName   string    `gorm:"column:first_name" json:"first_name"`
	LastName    *string   `gorm:"column:last_name" json:"last_name"`
	Phone       *string   `json:"phone"`
	Role        *string   `json:"role"`
	Designation *string   `json:"designation"`
	Image       *string   `json:"image"`
	Email       string    `gorm:"uniqueIndex" json:"email"`
	Password    string    `json:"-"`
	Status      bool      `json:"status"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	TenantID    uint      `gorm:"column:tenant_id" json:"-"`
	Tenant      Tenant    `gorm:"foreignKey:TenantID;references:ID" json:"-"`
}

type InstructorResponseLite struct {
	ID        uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName string  `gorm:"column:first_name" json:"first_name"`
	LastName  *string `gorm:"column:last_name" json:"last_name"`
	Email     string  `gorm:"uniqueIndex" json:"email"`
}

func (InstructorResponseLite) TableName() string {
	return "instructors"
}
