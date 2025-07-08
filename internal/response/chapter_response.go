package response

import (
	"dashlearn/internal/models"
	"time"
)

type CourseChapterResponse struct {
	ID          uint                       `gorm:"primaryKey;autoIncrement" json:"id"`
	Position    int                        `gorm:"default:0" json:"position"`
	Title       string                     `json:"title"`
	Description *string                    `gorm:"type:text" json:"description"`
	Access      models.Access              `gorm:"enum('draft','published');default:'published'" json:"access"`
	CreatedAt   time.Time                  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time                  `gorm:"autoUpdateTime" json:"updated_at"`
	CourseID    uint                       `gorm:"column:course_id" json:"course_id"`
	Lessons     []CourseLessonResponse     `gorm:"foreignKey:ChapterID;references:ID" json:"course_lessons"`
	Assignments []CourseAssignmentResponse `gorm:"foreignKey:ChapterID;references:ID" json:"assignments"`
	Quizzes     []CourseQuizResponse       `gorm:"foreignKey:ChapterID;references:ID" json:"quizzes"`
}

func (CourseChapterResponse) TableName() string {
	return "course_chapters"
}
