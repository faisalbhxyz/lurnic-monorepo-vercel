package student

import (
	"dashlearn/internal/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lucsky/cuid"
	"gorm.io/gorm"
)

type studentSessionInput struct {
	DeviceID   string
	DeviceName *string
	UserAgent  string
	IPAddress  string
}

func replaceStudentSession(db *gorm.DB, tenantID, studentID uint, input studentSessionInput) (string, error) {
	sessionID := cuid.New()
	now := time.Now()
	deviceName := resolveDeviceName(input.DeviceName, input.UserAgent)

	session := models.StudentSession{
		SessionID:  sessionID,
		StudentID:  studentID,
		TenantID:   tenantID,
		DeviceID:   input.DeviceID,
		DeviceName: &deviceName,
		LastSeenAt: now,
	}
	if input.UserAgent != "" {
		ua := truncateString(input.UserAgent, 512)
		session.UserAgent = &ua
	}
	if input.IPAddress != "" {
		ip := truncateString(input.IPAddress, 45)
		session.IPAddress = &ip
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("student_id = ?", studentID).Delete(&models.StudentSession{}).Error; err != nil {
			return err
		}
		return tx.Create(&session).Error
	})
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func invalidateStudentSession(db *gorm.DB, studentID uint) error {
	return db.Where("student_id = ?", studentID).Delete(&models.StudentSession{}).Error
}

func validateStudentSession(db *gorm.DB, studentID uint, sessionID string) (*models.StudentSession, error) {
	var session models.StudentSession
	err := db.Where("student_id = ? AND session_id = ?", studentID, sessionID).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func touchStudentSession(db *gorm.DB, sessionID string) {
	_ = db.Model(&models.StudentSession{}).
		Where("session_id = ?", sessionID).
		Update("last_seen_at", time.Now()).Error
}

func getStudentSessionForAdmin(db *gorm.DB, studentID uint) *models.StudentSession {
	var session models.StudentSession
	if err := db.Where("student_id = ?", studentID).First(&session).Error; err != nil {
		return nil
	}
	return &session
}

func sessionInputFromRequest(c *gin.Context, deviceID string, deviceName *string) studentSessionInput {
	ua := strings.TrimSpace(c.GetHeader("User-Agent"))
	return studentSessionInput{
		DeviceID:   deviceID,
		DeviceName: deviceName,
		UserAgent:  ua,
		IPAddress:  strings.TrimSpace(c.ClientIP()),
	}
}

func resolveDeviceName(deviceName *string, userAgent string) string {
	if deviceName != nil {
		name := strings.TrimSpace(*deviceName)
		if name != "" {
			return truncateString(name, 255)
		}
	}
	return truncateString(describeUserAgent(userAgent), 255)
}

func describeUserAgent(userAgent string) string {
	ua := strings.ToLower(userAgent)
	if ua == "" {
		return "Unknown device"
	}

	platform := "Unknown"
	switch {
	case strings.Contains(ua, "iphone"):
		platform = "iPhone"
	case strings.Contains(ua, "ipad"):
		platform = "iPad"
	case strings.Contains(ua, "android"):
		platform = "Android"
	case strings.Contains(ua, "windows"):
		platform = "Windows"
	case strings.Contains(ua, "mac os") || strings.Contains(ua, "macintosh"):
		platform = "Mac"
	case strings.Contains(ua, "linux"):
		platform = "Linux"
	}

	browser := "Browser"
	switch {
	case strings.Contains(ua, "edg/"):
		browser = "Edge"
	case strings.Contains(ua, "chrome/") && !strings.Contains(ua, "edg/"):
		browser = "Chrome"
	case strings.Contains(ua, "safari/") && !strings.Contains(ua, "chrome/"):
		browser = "Safari"
	case strings.Contains(ua, "firefox/"):
		browser = "Firefox"
	}

	return browser + " on " + platform
}

func truncateString(value string, max int) string {
	if len(value) <= max {
		return value
	}
	return value[:max]
}
