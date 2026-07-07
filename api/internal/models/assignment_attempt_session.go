package models

import "time"

type AssignmentAttemptSession struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID     uint       `gorm:"column:tenant_id" json:"-"`
	StudentID    uint       `gorm:"column:student_id" json:"student_id"`
	AssignmentID uint       `gorm:"column:assignment_id" json:"assignment_id"`
	StartedAt    time.Time  `gorm:"column:started_at" json:"started_at"`
	ExpiresAt    *time.Time `gorm:"column:expires_at" json:"expires_at,omitempty"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

func (AssignmentAttemptSession) TableName() string {
	return "assignment_attempt_sessions"
}

func (s *AssignmentAttemptSession) IsExpired() bool {
	if s.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*s.ExpiresAt)
}
