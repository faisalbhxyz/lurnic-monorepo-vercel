package models

import "time"

// StudentFreeLesson is a lesson the student added to their Free Class library.
type StudentFreeLesson struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID  uint      `gorm:"column:tenant_id;index" json:"-"`
	StudentID uint      `gorm:"column:student_id;index" json:"-"`
	LessonID  uint      `gorm:"column:lesson_id;index" json:"lesson_id"`
	CourseID  uint      `gorm:"column:course_id" json:"course_id"`
	AddedAt   time.Time `gorm:"column:added_at" json:"added_at"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (StudentFreeLesson) TableName() string {
	return "student_free_lessons"
}
