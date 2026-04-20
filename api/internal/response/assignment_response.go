package response

import (
	"dashlearn/internal/models"
	"time"

	"gorm.io/datatypes"
)

type CourseAssignmentResponse struct {
	ID               uint                             `json:"id"`
	CourseID         uint                             `json:"course_id"`
	ChapterID        uint                             `json:"chapter_id"`
	Title            string                           `json:"title"`
	Instructions     string                           `json:"instructions"`
	Attachments      *datatypes.JSON                  `json:"attachments"`
	IsPublished      bool                             `json:"is_published"`
	TimeLimit        int                              `json:"time_limit"`
	TimeLimitOption  models.CourseQuizTimeLimitOption `json:"time_limit_option"`
	FileUploadLimit  int                              `json:"file_upload_limit"`
	TotalMarks       float32                          `json:"total_marks"`
	MinimumPassMarks float32                          `json:"minimum_pass_marks"`
	CreatedAt        time.Time                        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time                        `gorm:"autoUpdateTime" json:"updated_at"`
}
