package course

import (
	"dashlearn/internal/models"
	"dashlearn/internal/utils"
	"encoding/json"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const defaultCertificateTemplate = "/images/Certificat-14.jpg"

func normalizeCertificateSettings(input CreateCertificateSettings) models.CourseCertificateSettings {
	templatePath := strings.TrimSpace(input.TemplatePath)
	if templatePath == "" {
		templatePath = defaultCertificateTemplate
	}

	completionPercent := input.CompletionPercent
	if completionPercent == 0 {
		completionPercent = 100
	}
	if completionPercent > 100 {
		completionPercent = 100
	}

	countLessons := input.CountLessons
	countQuizzes := input.CountQuizzes
	countAssignments := input.CountAssignments
	if !countLessons && !countQuizzes && !countAssignments {
		countLessons = true
		countQuizzes = true
		countAssignments = true
	}

	return models.CourseCertificateSettings{
		IsEnabled:           input.IsEnabled,
		CompletionPercent:   completionPercent,
		CountLessons:        countLessons,
		CountQuizzes:        countQuizzes,
		CountAssignments:    countAssignments,
		TemplatePath:        templatePath,
		Title:               input.Title,
		SubtitleOne:         input.SubtitleOne,
		SubtitleTwo:         input.SubtitleTwo,
		OwnerSignature:      input.OwnerSignature,
		InstructorSignature: input.InstructorSignature,
	}
}

func upsertCertificateSettings(db *gorm.DB, courseID uint, input CreateCertificateSettings) error {
	payload := normalizeCertificateSettings(input)
	payload.CourseID = courseID

	var existing models.CourseCertificateSettings
	err := db.Where("course_id = ?", courseID).First(&existing).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return db.Create(&payload).Error
		}
		return err
	}

	if payload.OwnerSignature == nil {
		payload.OwnerSignature = existing.OwnerSignature
	}
	if payload.InstructorSignature == nil {
		payload.InstructorSignature = existing.InstructorSignature
	}

	return db.Model(&existing).Select(
		"is_enabled",
		"completion_percent",
		"count_lessons",
		"count_quizzes",
		"count_assignments",
		"template_path",
		"title",
		"subtitle_one",
		"subtitle_two",
		"owner_signature",
		"instructor_signature",
	).Updates(payload).Error
}

func applyCertificateSettingsFromRequest(c *gin.Context, input *CourseDetailsInput) error {
	if settingsStr := c.PostForm("certificate_settings"); settingsStr != "" {
		var settings CreateCertificateSettings
		if err := json.Unmarshal([]byte(settingsStr), &settings); err != nil {
			return err
		}
		input.CertificateSettings = settings
	}

	if ownerHeader, err := c.FormFile("owner_signature"); err == nil {
		file, err := ownerHeader.Open()
		if err != nil {
			return err
		}
		defer file.Close()

		url, err := utils.UploadToBunny(file, ownerHeader)
		if err != nil {
			return err
		}
		input.CertificateSettings.OwnerSignature = &url
	}

	if instructorHeader, err := c.FormFile("instructor_signature"); err == nil {
		file, err := instructorHeader.Open()
		if err != nil {
			return err
		}
		defer file.Close()

		url, err := utils.UploadToBunny(file, instructorHeader)
		if err != nil {
			return err
		}
		input.CertificateSettings.InstructorSignature = &url
	}

	return nil
}
