package models

import (
	"time"

	"gorm.io/datatypes"
)

type QuizAttemptSession struct {
	ID            uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID      uint           `gorm:"column:tenant_id" json:"-"`
	StudentID     uint           `gorm:"column:student_id" json:"student_id"`
	QuizID        uint           `gorm:"column:quiz_id" json:"quiz_id"`
	AttemptNumber int            `gorm:"column:attempt_number" json:"attempt_number"`
	QuestionOrder datatypes.JSON `gorm:"column:question_order;type:json" json:"question_order"`
	StartedAt     time.Time      `gorm:"column:started_at" json:"started_at"`
	ExpiresAt     *time.Time     `gorm:"column:expires_at" json:"expires_at,omitempty"`
	SubmittedAt   *time.Time     `gorm:"column:submitted_at" json:"submitted_at,omitempty"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}

func (QuizAttemptSession) TableName() string {
	return "quiz_attempt_sessions"
}

func (s *QuizAttemptSession) IsExpired() bool {
	if s.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*s.ExpiresAt)
}
