package response

import (
	"dashlearn/internal/models"
	"time"
)

type CourseGeneralSettingsResponse struct {
	ID              uint                    `json:"id"`
	CourseID        uint                    `json:"course_id"`
	DifficultyLevel *models.DifficultyLevel `json:"difficulty_level"`
	MaximumStudent  *int32                  `json:"maximum_student"`
	Language        *string                 `json:"language"`
	CategoryID      uint                    `json:"category_id,omitempty"`
	Category        CategoryResponse        `json:"category"`
	Duration        *string                 `json:"duration"`
	CreatedAt       time.Time               `json:"created_at"`
	UpdatedAt       time.Time               `json:"updated_at"`
}
