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

const defaultCertificateTemplate = "/templates/minar-academy"
const TemplateMinarAcademy = "/templates/minar-academy"
const defaultWatermarkOpacity uint8 = 30

func normalizeWatermarkOpacity(value *uint8) uint8 {
	if value == nil {
		return defaultWatermarkOpacity
	}
	if *value > 100 {
		return 100
	}
	return *value
}

func normalizeCertificateSettings(input CreateCertificateSettings) models.CourseCertificateSettings {
	templatePath := strings.TrimSpace(input.TemplatePath)
	if templatePath == "" || strings.HasPrefix(templatePath, "/images/Certificat-") {
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
		BrandLogo:           input.BrandLogo,
		WatermarkImage:      input.WatermarkImage,
		WatermarkOpacity:    normalizeWatermarkOpacity(input.WatermarkOpacity),
		OrganizationName:    input.OrganizationName,
		SignerName:          input.SignerName,
		SignerRole:          input.SignerRole,
		SignerOrg:           input.SignerOrg,
		DualSignersEnabled:  input.DualSignersEnabled,
		Signer2Name:         input.Signer2Name,
		Signer2Role:         input.Signer2Role,
		Signer2Org:          input.Signer2Org,
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
	if payload.BrandLogo == nil {
		payload.BrandLogo = existing.BrandLogo
	}
	if payload.WatermarkImage == nil {
		payload.WatermarkImage = existing.WatermarkImage
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
		"brand_logo",
		"watermark_image",
		"watermark_opacity",
		"organization_name",
		"signer_name",
		"signer_role",
		"signer_org",
		"dual_signers_enabled",
		"signer2_name",
		"signer2_role",
		"signer2_org",
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

	if brandLogoHeader, err := c.FormFile("certificate_brand_logo"); err == nil {
		file, err := brandLogoHeader.Open()
		if err != nil {
			return err
		}
		defer file.Close()

		url, err := utils.UploadToBunny(file, brandLogoHeader)
		if err != nil {
			return err
		}
		input.CertificateSettings.BrandLogo = &url
	}

	if watermarkHeader, err := c.FormFile("certificate_watermark_image"); err == nil {
		file, err := watermarkHeader.Open()
		if err != nil {
			return err
		}
		defer file.Close()

		url, err := utils.UploadToBunny(file, watermarkHeader)
		if err != nil {
			return err
		}
		input.CertificateSettings.WatermarkImage = &url
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
