package certificate

import (
	"crypto/rand"
	"dashlearn/internal/models"
	"dashlearn/internal/progress"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lucsky/cuid"
	"gorm.io/gorm"
)

const issuedAtFormat = "2006-01-02 15:04:05"

func newCertificateNumber() string {
	buf := make([]byte, 7)
	if _, err := rand.Read(buf); err != nil {
		return strings.ToLower(cuid.New())
	}
	return hex.EncodeToString(buf)
}

func (s *Service) generateUniqueCertificateNumber() (string, error) {
	for range 10 {
		number := newCertificateNumber()
		var count int64
		if err := s.db.Model(&models.StudentCertificate{}).
			Where("certificate_number = ?", number).
			Count(&count).Error; err != nil {
			return "", err
		}
		if count == 0 {
			return number, nil
		}
	}
	return "", fmt.Errorf("failed to generate unique certificate number")
}

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

type CertificateResponse struct {
	ID                  uint                      `json:"id"`
	CourseID            uint                      `json:"course_id"`
	CourseTitle         string                    `json:"course_title"`
	CertificateNumber   string                    `json:"certificate_number"`
	StudentName         string                    `json:"student_name"`
	ProgressPercent     float32                   `json:"progress_percent"`
	TemplatePath        string                    `json:"template_path"`
	Title               *string                   `json:"title"`
	SubtitleOne         *string                   `json:"subtitle_one"`
	SubtitleTwo         *string                   `json:"subtitle_two"`
	BrandLogo           *string                   `json:"brand_logo"`
	WatermarkImage      *string                   `json:"watermark_image"`
	WatermarkOpacity    uint8                     `json:"watermark_opacity"`
	OrganizationName    *string                   `json:"organization_name"`
	SignerName          *string                   `json:"signer_name"`
	SignerRole          *string                   `json:"signer_role"`
	SignerOrg           *string                   `json:"signer_org"`
	DualSignersEnabled  bool                      `json:"dual_signers_enabled"`
	Signer2Name         *string                   `json:"signer2_name"`
	Signer2Role         *string                   `json:"signer2_role"`
	Signer2Org          *string                   `json:"signer2_org"`
	PricingModel        models.CoursePricingModel `json:"pricing_model"`
	OwnerSignature      *string                   `json:"owner_signature"`
	InstructorSignature *string                   `json:"instructor_signature"`
	IssuedAt            time.Time                 `json:"issued_at"`
	DownloadURL         string                    `json:"download_url,omitempty"`
}

func toResponse(cert models.StudentCertificate) CertificateResponse {
	return CertificateResponse{
		ID:                  cert.ID,
		CourseID:            cert.CourseID,
		CourseTitle:         cert.CourseTitle,
		CertificateNumber:   cert.CertificateNumber,
		StudentName:         cert.StudentName,
		ProgressPercent:     cert.ProgressPercent,
		TemplatePath:        cert.TemplatePath,
		Title:               cert.Title,
		SubtitleOne:         cert.SubtitleOne,
		SubtitleTwo:         cert.SubtitleTwo,
		BrandLogo:           cert.BrandLogo,
		WatermarkImage:      cert.WatermarkImage,
		WatermarkOpacity:    cert.WatermarkOpacity,
		OrganizationName:    cert.OrganizationName,
		SignerName:          cert.SignerName,
		SignerRole:          cert.SignerRole,
		SignerOrg:           cert.SignerOrg,
		DualSignersEnabled:  cert.DualSignersEnabled,
		Signer2Name:         cert.Signer2Name,
		Signer2Role:         cert.Signer2Role,
		Signer2Org:          cert.Signer2Org,
		PricingModel:        cert.PricingModel,
		OwnerSignature:      cert.OwnerSignature,
		InstructorSignature: cert.InstructorSignature,
		IssuedAt:            cert.IssuedAt,
	}
}

func (s *Service) ListForStudent(tenantID, studentID uint) ([]CertificateResponse, error) {
	var rows []models.StudentCertificate
	if err := s.db.
		Where("tenant_id = ? AND student_id = ?", tenantID, studentID).
		Order("issued_at DESC").
		Find(&rows).Error; err != nil {
		return nil, err
	}

	out := make([]CertificateResponse, 0, len(rows))
	for _, row := range rows {
		out = append(out, toResponse(row))
	}
	return out, nil
}

func (s *Service) GetForStudent(tenantID, studentID, certificateID uint) (*CertificateResponse, error) {
	cert, err := s.GetCertificateModel(tenantID, studentID, certificateID)
	if err != nil {
		return nil, err
	}
	res := toResponse(*cert)
	return &res, nil
}

func (s *Service) GetCertificateModel(tenantID, studentID, certificateID uint) (*models.StudentCertificate, error) {
	var cert models.StudentCertificate
	if err := s.db.
		Where("tenant_id = ? AND student_id = ? AND id = ?", tenantID, studentID, certificateID).
		First(&cert).Error; err != nil {
		return nil, err
	}
	return &cert, nil
}

func (s *Service) GetForCourseSlug(tenantID, studentID uint, slug string) (*CertificateResponse, error) {
	var course models.CourseDetails
	if err := s.db.
		Where("tenant_id = ? AND slug = ?", tenantID, slug).
		First(&course).Error; err != nil {
		return nil, err
	}

	var cert models.StudentCertificate
	if err := s.db.
		Where("tenant_id = ? AND student_id = ? AND course_id = ?", tenantID, studentID, course.ID).
		First(&cert).Error; err != nil {
		return nil, err
	}

	res := toResponse(cert)
	return &res, nil
}

func (s *Service) TryIssueCertificate(tenantID, studentID, courseID uint) (*models.StudentCertificate, error) {
	var settings models.CourseCertificateSettings
	if err := s.db.Where("course_id = ?", courseID).First(&settings).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	if !settings.IsEnabled {
		return nil, nil
	}

	var existing models.StudentCertificate
	err := s.db.
		Where("tenant_id = ? AND student_id = ? AND course_id = ?", tenantID, studentID, courseID).
		First(&existing).Error
	if err == nil {
		return &existing, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	currentProgress := progress.CalcCourseProgress(
		s.db,
		tenantID,
		studentID,
		courseID,
		progress.LoadOptions(s.db, courseID),
	)
	if currentProgress < float32(settings.CompletionPercent) {
		return nil, nil
	}

	var student models.Student
	if err := s.db.Where("id = ? AND tenant_id = ?", studentID, tenantID).First(&student).Error; err != nil {
		return nil, err
	}

	var course models.CourseDetails
	if err := s.db.Where("id = ? AND tenant_id = ?", courseID, tenantID).First(&course).Error; err != nil {
		return nil, err
	}

	studentName := strings.TrimSpace(student.FirstName)
	if student.LastName != nil && strings.TrimSpace(*student.LastName) != "" {
		studentName = strings.TrimSpace(studentName + " " + strings.TrimSpace(*student.LastName))
	}

	title := settings.Title
	if title == nil || strings.TrimSpace(*title) == "" {
		defaultTitle := "Certificate of Completion"
		title = &defaultTitle
	}

	now := time.Now()
	certificateNumber, err := s.generateUniqueCertificateNumber()
	if err != nil {
		return nil, err
	}

	cert := models.StudentCertificate{
		TenantID:            tenantID,
		StudentID:           studentID,
		CourseID:            courseID,
		CertificateNumber:   certificateNumber,
		StudentName:         studentName,
		CourseTitle:         course.Title,
		ProgressPercent:     currentProgress,
		TemplatePath:        settings.TemplatePath,
		Title:               title,
		SubtitleOne:         settings.SubtitleOne,
		SubtitleTwo:         settings.SubtitleTwo,
		BrandLogo:           settings.BrandLogo,
		WatermarkImage:      settings.WatermarkImage,
		WatermarkOpacity:    settings.WatermarkOpacity,
		OrganizationName:    settings.OrganizationName,
		SignerName:          settings.SignerName,
		SignerRole:          settings.SignerRole,
		SignerOrg:           settings.SignerOrg,
		DualSignersEnabled:  settings.DualSignersEnabled,
		Signer2Name:         settings.Signer2Name,
		Signer2Role:         settings.Signer2Role,
		Signer2Org:          settings.Signer2Org,
		PricingModel:        course.PricingModel,
		OwnerSignature:      settings.OwnerSignature,
		InstructorSignature: settings.InstructorSignature,
		IssuedAt:            now,
	}

	if err := s.db.Create(&cert).Error; err != nil {
		return nil, err
	}

	return &cert, nil
}
