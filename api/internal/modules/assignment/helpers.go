package assignment

import (
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"dashlearn/internal/models"
)

const (
	MaxAssignmentFileSizeBytes    int64 = 2 * 1024 * 1024
	MaxAssignmentResponseTextLength int   = 50000
)

var AllowedAssignmentMimeTypes = []string{
	"application/pdf",
	"image/jpeg",
	"image/png",
	"image/gif",
	"image/webp",
	"application/zip",
	"application/x-zip-compressed",
	"application/msword",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	"text/plain",
}

func assignmentTimeLimitDuration(limit int, option models.CourseAssignmentTimeLimitOption) time.Duration {
	if limit <= 0 {
		return 0
	}
	switch option {
	case models.CourseAssignmentTimeLimitOptionMinute:
		return time.Duration(limit) * time.Minute
	case models.CourseAssignmentTimeLimitOptionHour:
		return time.Duration(limit) * time.Hour
	case models.CourseAssignmentTimeLimitOptionDay:
		return time.Duration(limit) * 24 * time.Hour
	case models.CourseAssignmentTimeLimitOptionWeek:
		return time.Duration(limit) * 7 * 24 * time.Hour
	case models.CourseAssignmentTimeLimitOptionMonth:
		return time.Duration(limit) * 30 * 24 * time.Hour
	default:
		return time.Duration(limit) * time.Minute
	}
}

func secondsRemaining(expiresAt *time.Time) *int {
	if expiresAt == nil {
		return nil
	}
	remaining := int(time.Until(*expiresAt).Seconds())
	if remaining < 0 {
		remaining = 0
	}
	return &remaining
}

func isAllowedMimeType(mimeType string) bool {
	mimeType = strings.TrimSpace(strings.ToLower(mimeType))
	if mimeType == "" {
		return false
	}
	for _, allowed := range AllowedAssignmentMimeTypes {
		allowed = strings.ToLower(allowed)
		if strings.HasSuffix(allowed, "/*") {
			prefix := strings.TrimSuffix(allowed, "/*")
			if strings.HasPrefix(mimeType, prefix+"/") {
				return true
			}
			continue
		}
		if mimeType == allowed {
			return true
		}
	}
	return false
}

func validateSubmissionFile(header *multipart.FileHeader) error {
	if header.Size > MaxAssignmentFileSizeBytes {
		return fmt.Errorf("file %q exceeds maximum size of 2 MB", header.Filename)
	}
	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	if !isAllowedMimeType(mimeType) {
		return fmt.Errorf("file type %q is not allowed", mimeType)
	}
	return nil
}

func submissionEditable(status models.AssignmentSubmissionStatus, sessionExpired bool) bool {
	return status == models.AssignmentSubmissionStatusPendingReview && !sessionExpired
}
