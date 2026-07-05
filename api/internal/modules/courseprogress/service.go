package courseprogress

import (
	"dashlearn/internal/models"
	"dashlearn/internal/modules/certificate"
	"dashlearn/internal/progress"
	"errors"
	"time"

	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) GetCourseProgress(tenantID, studentID uint, slug string) (*progress.Breakdown, error) {
	course, err := s.loadEnrolledCourse(tenantID, studentID, slug)
	if err != nil {
		return nil, err
	}

	opts := progress.LoadOptions(s.db, course.ID)
	breakdown := progress.CalcBreakdown(s.db, tenantID, studentID, course.ID, opts, true)
	return &breakdown, nil
}

func (s *Service) MarkLessonComplete(tenantID, studentID uint, slug string, lessonID uint) (*progress.Breakdown, error) {
	course, err := s.loadEnrolledCourse(tenantID, studentID, slug)
	if err != nil {
		return nil, err
	}

	var lesson models.CourseLesson
	if err := s.db.
		Joins("JOIN course_chapters ON course_chapters.id = course_lessons.chapter_id").
		Where("course_lessons.id = ? AND course_chapters.course_id = ? AND course_lessons.is_published = ?", lessonID, course.ID, true).
		First(&lesson).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("lesson not found")
		}
		return nil, err
	}

	var existing models.StudentLessonCompletion
	err = s.db.
		Where("tenant_id = ? AND student_id = ? AND lesson_id = ?", tenantID, studentID, lessonID).
		First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		completion := models.StudentLessonCompletion{
			TenantID:    tenantID,
			StudentID:   studentID,
			CourseID:    course.ID,
			LessonID:    lessonID,
			CompletedAt: time.Now(),
		}
		if err := s.db.Create(&completion).Error; err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	_, _ = certificate.NewService(s.db).TryIssueCertificate(tenantID, studentID, course.ID)

	opts := progress.LoadOptions(s.db, course.ID)
	breakdown := progress.CalcBreakdown(s.db, tenantID, studentID, course.ID, opts, true)
	return &breakdown, nil
}

func (s *Service) loadEnrolledCourse(tenantID, studentID uint, slug string) (*models.CourseDetails, error) {
	var course models.CourseDetails
	if err := s.db.Where("tenant_id = ? AND slug = ?", tenantID, slug).First(&course).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("course not found")
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
		return nil, errors.New("enrollment required")
	}

	return &course, nil
}
