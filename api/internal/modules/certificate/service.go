package certificate

import (
	"dashlearn/internal/models"
	"dashlearn/internal/progress"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lucsky/cuid"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

type CertificateResponse struct {
	ID                  uint      `json:"id"`
	CourseID            uint      `json:"course_id"`
	CourseTitle         string    `json:"course_title"`
	CertificateNumber   string    `json:"certificate_number"`
	StudentName         string    `json:"student_name"`
	ProgressPercent     float32   `json:"progress_percent"`
	TemplatePath        string    `json:"template_path"`
	Title               *string   `json:"title"`
	SubtitleOne         *string   `json:"subtitle_one"`
	SubtitleTwo         *string   `json:"subtitle_two"`
	OwnerSignature      *string   `json:"owner_signature"`
	InstructorSignature *string   `json:"instructor_signature"`
	IssuedAt            time.Time `json:"issued_at"`
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
	var cert models.StudentCertificate
	if err := s.db.
		Where("tenant_id = ? AND student_id = ? AND id = ?", tenantID, studentID, certificateID).
		First(&cert).Error; err != nil {
		return nil, err
	}
	res := toResponse(cert)
	return &res, nil
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
	cert := models.StudentCertificate{
		TenantID:            tenantID,
		StudentID:           studentID,
		CourseID:            courseID,
		CertificateNumber:   fmt.Sprintf("CERT-%s", strings.ToUpper(cuid.New())),
		StudentName:         studentName,
		CourseTitle:         course.Title,
		ProgressPercent:     currentProgress,
		TemplatePath:        settings.TemplatePath,
		Title:               title,
		SubtitleOne:         settings.SubtitleOne,
		SubtitleTwo:         settings.SubtitleTwo,
		OwnerSignature:      settings.OwnerSignature,
		InstructorSignature: settings.InstructorSignature,
		IssuedAt:            now,
	}

	if err := s.db.Create(&cert).Error; err != nil {
		return nil, err
	}

	return &cert, nil
}
