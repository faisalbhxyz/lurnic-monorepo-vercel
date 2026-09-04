package freelesson

import (
	"dashlearn/internal/models"
	"dashlearn/internal/utils"
	"time"
)

// CatalogItem is one free/public lesson in the storefront catalog.
type CatalogItem struct {
	LessonID        uint                       `json:"lesson_id"`
	LessonTitle     string                     `json:"lesson_title"`
	ChapterID       uint                       `json:"chapter_id"`
	ChapterTitle    string                     `json:"chapter_title"`
	CourseID        uint                       `json:"course_id"`
	CourseSlug      string                     `json:"course_slug"`
	CourseTitle     string                     `json:"course_title"`
	FeaturedImage   *string                    `json:"featured_image"`
	LessonType      models.LessonType          `json:"lesson_type"`
	SourceType      models.LessonSourceType    `json:"source_type"`
	Source          utils.JSONB[models.Source] `json:"source"`
	IsPublic        bool                       `json:"is_public"`
	ClassSlugs      []string                   `json:"class_slugs"`
	DurationSeconds *float64                   `json:"duration_seconds,omitempty"`
}

// LibraryItem is a lesson in the student's free-lesson library (with watch progress).
type LibraryItem struct {
	LessonID        uint                       `json:"lesson_id"`
	LessonTitle     string                     `json:"lesson_title"`
	ChapterTitle    string                     `json:"chapter_title"`
	CourseID        uint                       `json:"course_id"`
	CourseSlug      string                     `json:"course_slug"`
	CourseTitle     string                     `json:"course_title"`
	FeaturedImage   *string                    `json:"featured_image"`
	SourceType      models.LessonSourceType    `json:"source_type"`
	Source          utils.JSONB[models.Source] `json:"source"`
	AddedAt         time.Time                  `json:"added_at"`
	WatchPercent    float64                    `json:"watch_percent"`
	WatchSeconds    float64                    `json:"watch_seconds"`
	DurationSeconds float64                    `json:"duration_seconds"`
	Completed       bool                       `json:"completed"`
}

type AddFreeLessonsRequest struct {
	LessonIDs []uint `json:"lesson_ids" binding:"required"`
}

type CatalogMeta struct {
	Total  int64 `json:"total"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
}
