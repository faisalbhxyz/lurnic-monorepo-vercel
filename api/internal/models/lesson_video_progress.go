package models

import "time"

type StudentLessonVideoProgress struct {
	ID                 uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID           uint      `gorm:"column:tenant_id" json:"-"`
	StudentID          uint      `gorm:"column:student_id" json:"-"`
	CourseID           uint      `gorm:"column:course_id" json:"-"`
	LessonID           uint      `gorm:"column:lesson_id" json:"lesson_id"`
	MaxPositionSeconds float64   `gorm:"column:max_position_seconds" json:"max_position_seconds"`
	DurationSeconds    float64   `gorm:"column:duration_seconds" json:"duration_seconds"`
	ProgressPercent    float64   `gorm:"column:progress_percent" json:"progress_percent"`
	CreatedAt          time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (StudentLessonVideoProgress) TableName() string {
	return "student_lesson_video_progress"
}
