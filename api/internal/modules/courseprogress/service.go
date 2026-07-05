package courseprogress

import (
	"dashlearn/internal/models"
	"dashlearn/internal/modules/certificate"
	"dashlearn/internal/progress"
	"errors"
	"math"
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

func (s *Service) GetLessonVideoProgress(tenantID, studentID uint, slug string, lessonID uint) (*LessonVideoProgressResponse, error) {
	course, err := s.loadEnrolledCourse(tenantID, studentID, slug)
	if err != nil {
		return nil, err
	}

	if _, err := s.loadPublishedLesson(course.ID, lessonID); err != nil {
		return nil, err
	}

	var row models.StudentLessonVideoProgress
	err = s.db.
		Where("tenant_id = ? AND student_id = ? AND lesson_id = ?", tenantID, studentID, lessonID).
		First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}
	if err != nil {
		return nil, err
	}

	return s.toLessonVideoProgressResponse(&row, s.isLessonCompleted(tenantID, studentID, lessonID)), nil
}

func (s *Service) UpdateLessonVideoProgress(tenantID, studentID uint, slug string, lessonID uint, req UpdateLessonVideoProgressRequest) (*LessonVideoProgressResponse, error) {
	course, err := s.loadEnrolledCourse(tenantID, studentID, slug)
	if err != nil {
		return nil, err
	}

	if _, err := s.loadPublishedLesson(course.ID, lessonID); err != nil {
		return nil, err
	}

	maxPosition := req.MaxPositionSeconds
	duration := req.DurationSeconds

	var existing models.StudentLessonVideoProgress
	err = s.db.
		Where("tenant_id = ? AND student_id = ? AND lesson_id = ?", tenantID, studentID, lessonID).
		First(&existing).Error
	if err == nil {
		if existing.MaxPositionSeconds > maxPosition {
			maxPosition = existing.MaxPositionSeconds
		}
		if existing.DurationSeconds > duration {
			duration = existing.DurationSeconds
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	progressPercent := calcVideoProgressPercent(maxPosition, duration)

	row := models.StudentLessonVideoProgress{
		TenantID:           tenantID,
		StudentID:          studentID,
		CourseID:           course.ID,
		LessonID:           lessonID,
		MaxPositionSeconds: maxPosition,
		DurationSeconds:    duration,
		ProgressPercent:    progressPercent,
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := s.db.Create(&row).Error; err != nil {
			return nil, err
		}
	} else {
		updates := map[string]interface{}{
			"max_position_seconds": maxPosition,
			"duration_seconds":     duration,
			"progress_percent":     progressPercent,
		}
		if err := s.db.Model(&existing).Updates(updates).Error; err != nil {
			return nil, err
		}
		row = existing
		row.MaxPositionSeconds = maxPosition
		row.DurationSeconds = duration
		row.ProgressPercent = progressPercent
		if err := s.db.First(&row, existing.ID).Error; err != nil {
			return nil, err
		}
	}

	return s.toLessonVideoProgressResponse(&row, s.isLessonCompleted(tenantID, studentID, lessonID)), nil
}

func (s *Service) MarkLessonComplete(tenantID, studentID uint, slug string, lessonID uint) (*progress.Breakdown, error) {
	course, err := s.loadEnrolledCourse(tenantID, studentID, slug)
	if err != nil {
		return nil, err
	}

	if _, err := s.loadPublishedLesson(course.ID, lessonID); err != nil {
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

func (s *Service) loadPublishedLesson(courseID, lessonID uint) (*models.CourseLesson, error) {
	var lesson models.CourseLesson
	if err := s.db.
		Joins("JOIN course_chapters ON course_chapters.id = course_lessons.chapter_id").
		Where("course_lessons.id = ? AND course_chapters.course_id = ? AND course_lessons.is_published = ?", lessonID, courseID, true).
		First(&lesson).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("lesson not found")
		}
		return nil, err
	}
	return &lesson, nil
}

func (s *Service) isLessonCompleted(tenantID, studentID, lessonID uint) bool {
	var count int64
	s.db.Model(&models.StudentLessonCompletion{}).
		Where("tenant_id = ? AND student_id = ? AND lesson_id = ?", tenantID, studentID, lessonID).
		Count(&count)
	return count > 0
}

func (s *Service) toLessonVideoProgressResponse(row *models.StudentLessonVideoProgress, completed bool) *LessonVideoProgressResponse {
	return &LessonVideoProgressResponse{
		LessonID:           row.LessonID,
		MaxPositionSeconds: row.MaxPositionSeconds,
		DurationSeconds:    row.DurationSeconds,
		ProgressPercent:    row.ProgressPercent,
		Completed:          completed,
		UpdatedAt:          row.UpdatedAt,
	}
}

func calcVideoProgressPercent(maxPosition, duration float64) float64 {
	if duration <= 0 {
		return 0
	}
	percent := maxPosition / duration * 100
	if percent > 100 {
		return 100
	}
	return math.Round(percent*10) / 10
}
