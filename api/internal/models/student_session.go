package models

import "time"

type StudentSession struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	SessionID  string    `gorm:"column:session_id;uniqueIndex" json:"session_id"`
	StudentID  uint      `gorm:"column:student_id;uniqueIndex" json:"student_id"`
	TenantID   uint      `gorm:"column:tenant_id" json:"tenant_id"`
	DeviceID   string    `gorm:"column:device_id" json:"device_id"`
	DeviceName *string   `gorm:"column:device_name" json:"device_name"`
	UserAgent  *string   `gorm:"column:user_agent" json:"user_agent"`
	IPAddress  *string   `gorm:"column:ip_address" json:"ip_address"`
	LastSeenAt time.Time `gorm:"column:last_seen_at" json:"last_seen_at"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (StudentSession) TableName() string {
	return "student_sessions"
}
