package response

import (
	"dashlearn/internal/models"
	"dashlearn/internal/utils"
	"time"
)

// LessonResourceResponse is the storefront contract for downloadable lesson files (PDF, etc.).
type LessonResourceResponse struct {
	ID        uint      `json:"id"`
	CourseID  uint      `json:"course_id"`
	LessonID  uint      `json:"lesson_id"`
	MimeType  string    `json:"mime_type"`
	Title     string    `json:"title"`
	FilePath  string    `json:"file_path"`
	URL       string    `json:"url"`
	Position  int       `json:"position"`
	FileSize  int64     `json:"file_size"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func MapLessonResources(resources []models.LessonResource) []LessonResourceResponse {
	if len(resources) == 0 {
		return nil
	}

	mapped := make([]LessonResourceResponse, 0, len(resources))
	for _, res := range resources {
		mapped = append(mapped, LessonResourceResponse{
			ID:        res.ID,
			CourseID:  res.CourseID,
			LessonID:  res.LessonID,
			MimeType:  res.MimeType,
			Title:     res.Title,
			FilePath:  res.FilePath,
			URL:       res.FilePath,
			Position:  res.Position,
			FileSize:  res.FileSize,
			CreatedAt: res.CreatedAt,
			UpdatedAt: res.UpdatedAt,
		})
	}
	return mapped
}

type CourseLessonResponse struct {
	ID                  uint                       `json:"id"`
	ChapterID           uint                       `json:"chapter_id"`
	Title               string                     `json:"title"`
	Description         *string                    `json:"description"`
	LessonType          models.LessonType          `json:"lesson_type"`
	SourceType          models.LessonSourceType    `json:"source_type"`
	Source              utils.JSONB[models.Source] `json:"source"`
	IsPublished         bool                       `json:"is_published,omitempty"`
	IsPublic            bool                       `json:"is_public"`
	Position            int                        `json:"position,omitempty"`
	Resources           []LessonResourceResponse   `json:"resources,omitempty"`
	OfflineDownloadable bool                       `json:"offline_downloadable"`
	DownloadURL         *string                    `json:"download_url,omitempty"`
	CreatedAt           time.Time                  `json:"created_at"`
	UpdatedAt           time.Time                  `json:"updated_at"`
}
