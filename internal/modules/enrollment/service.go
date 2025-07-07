package enrollment

import (
	"dashlearn/models"
	"dashlearn/response"
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
	var enrollments []response.EnrolledCourseRes

	err := s.db.
		Where(&models.Enrollment{
			TenantID:  tenantID,
			StudentID: studentID,
		}).
		Preload("Course").
		Find(&enrollments).Error

	return enrollments, err
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
