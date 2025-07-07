package response

import (
	"dashlearn/models"
	"dashlearn/utils"
	"time"

	"gorm.io/datatypes"
)

type CourseLessonResponse struct {
	ID          uint                       `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string                     `json:"title"`
	Description *string                    `gorm:"type:text" json:"description"`
	LessonType  models.LessonType          `gorm:"enum('video','live_session','audio','text');default:'video'" json:"lesson_type"`
	SourceType  models.LessonSourceType    `gorm:"enum('youtube','vimeo', 'sound_cloud','spotify','custom_code','upload');default:'youtube'" json:"source_type"`
	Source      utils.JSONB[models.Source] `gorm:"type:json" json:"source"`
	IsPublished bool                       `gorm:"default:false" json:"is_published"`
	IsPublic    bool                       `gorm:"default:false" json:"is_public"`
	Resources   datatypes.JSON             `gorm:"type:json" json:"resources"` // filename, mimetype, url, size
	Position    int                        `gorm:"default:0" json:"position"`
	CreatedAt   time.Time                  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time                  `gorm:"autoUpdateTime" json:"updated_at"`
	ChapterID   uint                       `gorm:"column:chapter_id" json:"chapter_id"`
}

func (CourseLessonResponse) TableName() string {
	return "course_lessons"
}
