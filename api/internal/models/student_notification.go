package models

import "time"

type StudentNotification struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID  uint       `gorm:"column:tenant_id" json:"-"`
	StudentID uint       `gorm:"column:student_id" json:"student_id"`
	Title     string     `json:"title"`
	Body      *string    `gorm:"type:text" json:"body"`
	Type      string     `gorm:"type:varchar(50);default:'general'" json:"type"`
	Data      []byte     `gorm:"type:json" json:"data,omitempty"`
	ReadAt    *time.Time `gorm:"column:read_at" json:"read_at"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

func (StudentNotification) TableName() string {
	return "student_notifications"
}
