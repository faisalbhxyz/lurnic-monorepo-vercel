package models

import "time"

type StudentLessonCompletion struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID    uint      `gorm:"column:tenant_id" json:"-"`
	StudentID   uint      `gorm:"column:student_id" json:"student_id"`
	CourseID    uint      `gorm:"column:course_id" json:"course_id"`
	LessonID    uint      `gorm:"column:lesson_id" json:"lesson_id"`
	CompletedAt time.Time `gorm:"column:completed_at" json:"completed_at"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
