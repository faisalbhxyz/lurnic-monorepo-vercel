package enrollment

import (
	"dashlearn/internal/models"
	coursemodule "dashlearn/internal/modules/course"
	"dashlearn/internal/response"
	"errors"

	"gorm.io/gorm"
)

type EnrollmentService interface {
	GetAll(tenantID uint) ([]response.EnrollmentResponse, error)
	GetEnrolledCourses(tenantID uint, studentID uint) ([]response.EnrolledCourseRes, error)
	Create(input models.Enrollment, tenantID uint) error
	Delete(id uint, tenantID uint) error
}

type enrollmentService struct {
	db *gorm.DB
}

func NewEnrollmentService(db *gorm.DB) EnrollmentService {
	return &enrollmentService{
		db: db,
	}
}

func (s *enrollmentService) GetAll(tenantID uint) ([]response.EnrollmentResponse, error) {
	var enrollments []response.EnrollmentResponse

	err := s.db.
		Where("tenant_id = ?", tenantID).
		Preload("Student", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "first_name", "last_name", "email")
		}).
		Preload("Course", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "title")
		}).
		Find(&enrollments).Error

	return enrollments, err
}

func (s *enrollmentService) GetEnrolledCourses(tenantID uint, studentID uint) ([]response.EnrolledCourseRes, error) {
	var modelEnrollments []models.Enrollment
	err := s.db.
		Where(&models.Enrollment{
			TenantID:  tenantID,
			StudentID: studentID,
		}).
		Order("created_at DESC").
		Find(&modelEnrollments).Error
	if err != nil {
		return nil, err
	}

	if len(modelEnrollments) == 0 {
		return []response.EnrolledCourseRes{}, nil
	}

	courseIDs := make([]uint, len(modelEnrollments))
	for i, enrollment := range modelEnrollments {
		courseIDs[i] = enrollment.CourseID
	}

	coursesByID, err := coursemodule.LoadPublicCoursesByIDs(s.db, tenantID, courseIDs)
	if err != nil {
		return nil, err
	}

	enrollments := make([]response.EnrolledCourseRes, 0, len(modelEnrollments))
	for _, enrollment := range modelEnrollments {
		course, ok := coursesByID[enrollment.CourseID]
		if !ok || course == nil {
			continue
		}

		enrollments = append(enrollments, response.EnrolledCourseRes{
			ID:        enrollment.ID,
			CourseID:  enrollment.CourseID,
			Course:    *course,
			StudentID: enrollment.StudentID,
			CreatedAt: enrollment.CreatedAt,
			UpdatedAt: enrollment.UpdatedAt,
		})
	}

	return enrollments, nil
}

func (s *enrollmentService) Create(input models.Enrollment, tenantID uint) error {

	// check if student is already enrolled in this course
	if s.db.Where("student_id = ? AND course_id = ? AND tenant_id = ?", input.StudentID, input.CourseID, tenantID).First(&models.Enrollment{}).RowsAffected > 0 {
		return errors.New("student is already enrolled in this course")
	}

	newEnrollment := models.Enrollment{
		StudentID: input.StudentID,
		CourseID:  input.CourseID,
		TenantID:  tenantID,
	}

	return s.db.Create(&newEnrollment).Error
}

func (s *enrollmentService) Delete(id uint, tenantID uint) error {
	return s.db.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&models.Enrollment{}).Error
}
