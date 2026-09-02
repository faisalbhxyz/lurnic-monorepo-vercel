package course

import (
	"errors"
	"regexp"
	"strings"

	"dashlearn/internal/models"
	"dashlearn/internal/response"
	"dashlearn/internal/utils"
)

var driveFileIDPatterns = []*regexp.Regexp{
	regexp.MustCompile(`drive\.google\.com/file/d/([a-zA-Z0-9_-]+)`),
	regexp.MustCompile(`drive\.google\.com/open\?(?:[^#&]*&)*id=([a-zA-Z0-9_-]+)`),
	regexp.MustCompile(`drive\.google\.com/uc\?(?:[^#&]*&)*id=([a-zA-Z0-9_-]+)`),
	regexp.MustCompile(`docs\.google\.com/uc\?(?:[^#&]*&)*id=([a-zA-Z0-9_-]+)`),
}

var directVideoURLPattern = regexp.MustCompile(`(?i)\.(mp4|m4v|webm|mov)(\?|$)`)

func ExtractDriveFileID(raw string) (string, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", false
	}
	for _, pattern := range driveFileIDPatterns {
		m := pattern.FindStringSubmatch(raw)
		if len(m) >= 2 {
			return m[1], true
		}
	}
	return "", false
}

func IsDirectVideoURL(url string) bool {
	return directVideoURLPattern.MatchString(strings.TrimSpace(url))
}

func OfflineDownloadable(lessonType models.LessonType, source models.Source) bool {
	if lessonType != models.Video {
		return false
	}

	url := strings.TrimSpace(source.Data)
	if _, ok := ExtractDriveFileID(url); ok {
		return true
	}
	if IsDirectVideoURL(url) {
		return true
	}
	return false
}

func DownloadURLForLesson(source models.Source) string {
	url := strings.TrimSpace(source.Data)
	if url == "" {
		return ""
	}
	if _, ok := ExtractDriveFileID(url); ok {
		return url
	}
	if IsDirectVideoURL(url) {
		return url
	}
	return ""
}

func NormalizeLessonSource(
	lessonType models.LessonType,
	sourceType models.LessonSourceType,
	source models.Source,
) (models.LessonSourceType, models.Source) {
	source.Data = strings.TrimSpace(source.Data)
	if lessonType != models.Video {
		return sourceType, source
	}

	fileID, ok := ExtractDriveFileID(source.Data)
	if !ok {
		return sourceType, source
	}

	source.IsFile = false
	source.DriveFileID = &fileID
	return models.GoogleDrive, source
}

func ValidateLessonSource(sourceType models.LessonSourceType, source models.Source) error {
	if sourceType != models.GoogleDrive {
		return nil
	}
	if _, ok := ExtractDriveFileID(source.Data); !ok {
		return errors.New("google_drive source requires a valid Google Drive URL")
	}
	return nil
}

func prepareLessonSource(
	lessonType models.LessonType,
	sourceType models.LessonSourceType,
	source models.Source,
) (models.LessonSourceType, utils.JSONB[models.Source], error) {
	normalizedType, normalizedSource := NormalizeLessonSource(lessonType, sourceType, source)
	if err := ValidateLessonSource(normalizedType, normalizedSource); err != nil {
		return "", utils.JSONB[models.Source]{}, err
	}
	return normalizedType, utils.JSONB[models.Source]{Data: normalizedSource}, nil
}

func mapPublicLesson(lesson models.CourseLesson) response.CourseLessonResponse {
	source := lesson.Source.Data
	downloadURL := DownloadURLForLesson(source)

	res := response.CourseLessonResponse{
		ID:                  lesson.ID,
		Title:               lesson.Title,
		Description:         lesson.Description,
		Position:            lesson.Position,
		CreatedAt:           lesson.CreatedAt,
		UpdatedAt:           lesson.UpdatedAt,
		ChapterID:           lesson.ChapterID,
		LessonType:          lesson.LessonType,
		SourceType:          lesson.SourceType,
		Source:              lesson.Source,
		IsPublic:            lesson.IsPublic,
		IsPublished:         lesson.IsPublished,
		Resources:           response.MapLessonResources(lesson.Resources),
		OfflineDownloadable: OfflineDownloadable(lesson.LessonType, source),
	}
	if downloadURL != "" {
		res.DownloadURL = &downloadURL
	}
	return res
}
