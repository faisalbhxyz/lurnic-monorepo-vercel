package response

import (
	"dashlearn/internal/models"
	"time"
)

type CourseChapterResponse struct {
	ID          uint                       `json:"id"`
	Position    int                        `json:"position"`
	Title       string                     `json:"title"`
	Description *string                    `json:"description"`
	Access      models.Access              `json:"access"`
	CreatedAt   time.Time                  `json:"created_at"`
	UpdatedAt   time.Time                  `json:"updated_at"`
	CourseID    uint                       `json:"course_id"`
	Lessons     []CourseLessonResponse     `json:"course_lessons,omitempty"`
	Assignments []CourseAssignmentResponse `json:"assignments,omitempty"`
	Quizzes     []CourseQuizResponse       `json:"quizzes,omitempty"`
}
