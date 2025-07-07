package response

import (
	"dashlearn/models"
	"time"

	"gorm.io/datatypes"
)

type CourseAssignmentResponse struct {
	ID               uint                             `gorm:"primaryKey;autoIncrement" json:"id"`
	CourseID         uint                             `json:"course_id" gorm:"column:course_id"`
	ChapterID        uint                             `json:"chapter_id" gorm:"column:chapter_id"`
	Title            string                           `json:"title" gorm:"type:varchar(255)"`
	Instructions     string                           `json:"instructions" gorm:"type:text"`
	Attachments      *datatypes.JSON                  `gorm:"type:json" json:"attachments"`
	IsPublished      bool                             `json:"is_published" gorm:"column:is_published;default:false"`
	TimeLimit        int                              `json:"time_limit" gorm:"column:time_limit;default:1"`
	TimeLimitOption  models.CourseQuizTimeLimitOption `json:"time_limit_option" gorm:"column:time_limit_option;type:enum('minutes','hours','days','weeks','months');default:'weeks'"`
	FileUploadLimit  int                              `json:"file_upload_limit" gorm:"column:file_upload_limit;default:1"`
	TotalMarks       float32                          `json:"total_marks" gorm:"column:total_marks;default:1"`
	MinimumPassMarks float32                          `json:"minimum_pass_marks" gorm:"column:minimum_pass_marks;default:0"`
	CreatedAt        time.Time                        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time                        `gorm:"autoUpdateTime" json:"updated_at"`
}

func (CourseAssignmentResponse) TableName() string {
	return "course_assignments"
}
