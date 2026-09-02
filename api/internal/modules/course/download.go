package course

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"dashlearn/internal/models"

	"gorm.io/gorm"
)

var (
	ErrDownloadCourseNotFound = errors.New("course not found")
	ErrDownloadNotEnrolled    = errors.New("not enrolled")
	ErrDownloadLessonNotFound = errors.New("lesson not found")
	ErrDownloadNotDownloadable = errors.New("not downloadable")
)

type LessonDownloadResult struct {
	RedirectURL string
	FileName    string
	ContentType string
}

type DownloadService struct {
	db *gorm.DB
}

func NewDownloadService(db *gorm.DB) *DownloadService {
	return &DownloadService{db: db}
}

func (s *DownloadService) ResolveLessonDownload(tenantID, studentID uint, slug string, lessonID uint) (*LessonDownloadResult, error) {
	course, err := s.loadEnrolledCourse(tenantID, studentID, slug)
	if err != nil {
		return nil, err
	}

	lesson, err := s.loadPublishedLesson(course.ID, lessonID)
	if err != nil {
		return nil, err
	}

	return buildLessonDownloadResult(lesson)
}

func buildLessonDownloadResult(lesson *models.CourseLesson) (*LessonDownloadResult, error) {
	if !OfflineDownloadable(lesson.LessonType, lesson.Source.Data) {
		return nil, ErrDownloadNotDownloadable
	}

	sourceURL := strings.TrimSpace(lesson.Source.Data.Data)
	if sourceURL == "" {
		return nil, ErrDownloadNotDownloadable
	}

	if fileID, ok := ExtractDriveFileID(sourceURL); ok {
		return &LessonDownloadResult{
			RedirectURL: fmt.Sprintf(
				"https://drive.usercontent.google.com/download?id=%s&export=download&confirm=t",
				fileID,
			),
			FileName:    lessonDownloadFileName(lesson.Title, ".mp4"),
			ContentType: "video/mp4",
		}, nil
	}

	if IsDirectVideoURL(sourceURL) {
		ext := directVideoExtension(sourceURL)
		return &LessonDownloadResult{
			RedirectURL: sourceURL,
			FileName:    lessonDownloadFileName(lesson.Title, ext),
			ContentType: mimeTypeForVideoExt(ext),
		}, nil
	}

	return nil, ErrDownloadNotDownloadable
}

func (s *DownloadService) loadEnrolledCourse(tenantID, studentID uint, slug string) (*models.CourseDetails, error) {
	var course models.CourseDetails
	if err := s.db.Where("tenant_id = ? AND slug = ?", tenantID, slug).First(&course).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDownloadCourseNotFound
		}
		return nil, err
	}

	var count int64
	if err := s.db.Model(&models.Enrollment{}).
		Where("tenant_id = ? AND student_id = ? AND course_id = ?", tenantID, studentID, course.ID).
		Count(&count).Error; err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, ErrDownloadNotEnrolled
	}

	return &course, nil
}

func (s *DownloadService) loadPublishedLesson(courseID, lessonID uint) (*models.CourseLesson, error) {
	var lesson models.CourseLesson
	if err := s.db.
		Joins("JOIN course_chapters ON course_chapters.id = course_lessons.chapter_id").
		Where(
			"course_lessons.id = ? AND course_chapters.course_id = ? AND course_lessons.is_published = ?",
			lessonID, courseID, true,
		).
		First(&lesson).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDownloadLessonNotFound
		}
		return nil, err
	}
	return &lesson, nil
}

var unsafeFileNamePattern = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)

func lessonDownloadFileName(title, ext string) string {
	base := strings.TrimSpace(title)
	if base == "" {
		base = "lesson"
	}
	base = strings.ToLower(unsafeFileNamePattern.ReplaceAllString(base, "-"))
	base = strings.Trim(base, "-")
	if base == "" {
		base = "lesson"
	}
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	return base + ext
}

func directVideoExtension(url string) string {
	lower := strings.ToLower(url)
	switch {
	case strings.Contains(lower, ".m4v"):
		return ".m4v"
	case strings.Contains(lower, ".webm"):
		return ".webm"
	case strings.Contains(lower, ".mov"):
		return ".mov"
	default:
		return ".mp4"
	}
}

func mimeTypeForVideoExt(ext string) string {
	switch ext {
	case ".m4v":
		return "video/x-m4v"
	case ".webm":
		return "video/webm"
	case ".mov":
		return "video/quicktime"
	default:
		return "video/mp4"
	}
}
