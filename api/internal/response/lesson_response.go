package response

import (
	"dashlearn/internal/models"
	"dashlearn/internal/utils"
	"time"
)

type CourseLessonResponse struct {
	ID          uint                       `json:"id"`
	ChapterID   uint                       `json:"chapter_id"`
	Title       string                     `json:"title"`
	Description *string                    `json:"description"`
	LessonType  models.LessonType          `json:"lesson_type"`
	SourceType  models.LessonSourceType    `json:"source_type"`
	Source      utils.JSONB[models.Source] `json:"source"`
	IsPublished bool                       `json:"is_published,omitempty"`
	IsPublic    bool                       `json:"is_public"`
	Position    int                        `json:"position,omitempty"`
	Resources   []models.LessonResource    `json:"resources,omitempty"`
	CreatedAt   time.Time                  `json:"created_at"`
	UpdatedAt   time.Time                  `json:"updated_at"`
}
