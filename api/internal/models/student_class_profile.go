package models

import "time"

// StudentClassProfile stores the student's class / grade preference for cross-device sync.
type StudentClassProfile struct {
	StudentID           uint      `gorm:"primaryKey;column:student_id" json:"-"`
	ClassLevel          string    `gorm:"column:class_level;size:32;not null" json:"class_level"`
	HscBatch            *string   `gorm:"column:hsc_batch;size:32" json:"hsc_batch"`
	Department          *string   `gorm:"column:department;size:32" json:"department"`
	PreferredClassSlug  *string   `gorm:"column:preferred_class_slug;size:64;index" json:"preferred_class_slug"`
	OnboardingCompleted bool      `gorm:"column:onboarding_completed;not null;default:false" json:"onboarding_completed"`
	CreatedAt           time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (StudentClassProfile) TableName() string {
	return "student_class_profiles"
}
