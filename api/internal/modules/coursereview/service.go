package coursereview

import (
	"encoding/json"
	"errors"
	"log"
	"math"
	"strings"

	"dashlearn/internal/models"
	"dashlearn/internal/progress"
	"dashlearn/internal/response"

	"gorm.io/gorm"
)

var allowedTags = map[string]struct{}{
	"excellent_content":     {},
	"excellent_teaching":    {},
	"sufficient_resources":  {},
	"others":                {},
}

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

type ListReviewsParams struct {
	TenantID  uint
	Slug      string
	StudentID uint
	Authed    bool
	Page      int
	PerPage   int
}

func (s *Service) ListReviews(params ListReviewsParams) (*response.CourseReviewsSummary, error) {
	course, err := s.loadCourse(params.TenantID, params.Slug)
	if err != nil {
		return nil, err
	}

	page := params.Page
	if page < 1 {
		page = 1
	}
	perPage := params.PerPage
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 50 {
		perPage = 50
	}
	offset := (page - 1) * perPage

	emptySummary := func() *response.CourseReviewsSummary {
		return &response.CourseReviewsSummary{
			AverageRating: 0,
			TotalReviews:  0,
			Reviews:       []response.CourseReviewResponse{},
		}
	}

	var reviews []models.CourseReview
	if err := s.db.
		Preload("Student", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "first_name", "last_name", "profile_image")
		}).
		Where("tenant_id = ? AND course_id = ? AND status = ?", params.TenantID, course.ID, models.CourseReviewStatusPublished).
		Order("created_at DESC").
		Limit(perPage).
		Offset(offset).
		Find(&reviews).Error; err != nil {
		if isMissingTable(err) {
			log.Printf("[coursereview] course_reviews table missing; returning empty summary (run migration 00063)")
			return emptySummary(), nil
		}
		return nil, err
	}

	var total int64
	if err := s.db.Model(&models.CourseReview{}).
		Where("tenant_id = ? AND course_id = ? AND status = ?", params.TenantID, course.ID, models.CourseReviewStatusPublished).
		Count(&total).Error; err != nil {
		if isMissingTable(err) {
			return emptySummary(), nil
		}
		return nil, err
	}

	var avgRating float64
	if err := s.db.Model(&models.CourseReview{}).
		Where("tenant_id = ? AND course_id = ? AND status = ?", params.TenantID, course.ID, models.CourseReviewStatusPublished).
		Select("COALESCE(AVG(rating), 0)").
		Scan(&avgRating).Error; err != nil {
		if isMissingTable(err) {
			return emptySummary(), nil
		}
		return nil, err
	}

	summary := &response.CourseReviewsSummary{
		AverageRating: math.Round(avgRating*10) / 10,
		TotalReviews:  int(total),
		Reviews:       make([]response.CourseReviewResponse, 0, len(reviews)),
	}

	for _, review := range reviews {
		summary.Reviews = append(summary.Reviews, toReviewResponse(review))
	}

	if params.Authed && params.StudentID > 0 {
		var studentReview models.CourseReview
		err := s.db.
			Preload("Student", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "first_name", "last_name", "profile_image")
			}).
			Where("tenant_id = ? AND course_id = ? AND student_id = ?", params.TenantID, course.ID, params.StudentID).
			First(&studentReview).Error
		if err == nil {
			resp := toReviewResponse(studentReview)
			summary.StudentReview = &resp
		}

		canReview, err := s.canReview(params.TenantID, params.StudentID, course.ID)
		if err == nil {
			summary.CanReview = &canReview
		}
	}

	return summary, nil
}

type CreateReviewInput struct {
	Rating  uint8    `json:"rating" binding:"required,min=1,max=5"`
	Comment *string  `json:"comment"`
	Tags    []string `json:"tags"`
}

func (s *Service) CreateReview(tenantID, studentID uint, slug string, input CreateReviewInput) (*response.CourseReviewResponse, error) {
	course, err := s.loadCourse(tenantID, slug)
	if err != nil {
		return nil, err
	}

	if !s.isEnrolled(tenantID, studentID, course.ID) {
		return nil, errNotEnrolled
	}

	opts := progress.LoadOptions(s.db, course.ID)
	breakdown := progress.CalcBreakdown(s.db, tenantID, studentID, course.ID, opts, false)
	if breakdown.Percent < 100 {
		return nil, errCourseNotCompleted
	}

	var existing models.CourseReview
	err = s.db.Where("tenant_id = ? AND course_id = ? AND student_id = ?", tenantID, course.ID, studentID).First(&existing).Error
	if err == nil {
		return nil, errReviewExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if input.Comment != nil && len(strings.TrimSpace(*input.Comment)) > 2000 {
		return nil, errValidation
	}

	tagsJSON, err := encodeTags(input.Tags)
	if err != nil {
		return nil, err
	}

	review := models.CourseReview{
		TenantID:  tenantID,
		CourseID:  course.ID,
		StudentID: studentID,
		Rating:    input.Rating,
		Comment:   input.Comment,
		Tags:      tagsJSON,
		Status:    models.CourseReviewStatusPublished,
	}

	if err := s.db.Create(&review).Error; err != nil {
		return nil, err
	}

	if err := s.db.
		Preload("Student", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "first_name", "last_name", "profile_image")
		}).
		First(&review, review.ID).Error; err != nil {
		return nil, err
	}

	resp := toReviewResponse(review)
	return &resp, nil
}

func (s *Service) loadCourse(tenantID uint, slug string) (*models.CourseDetails, error) {
	var course models.CourseDetails
	if err := s.db.Where("tenant_id = ? AND slug = ?", tenantID, slug).First(&course).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errCourseNotFound
		}
		return nil, err
	}
	return &course, nil
}

func (s *Service) isEnrolled(tenantID, studentID, courseID uint) bool {
	var count int64
	s.db.Model(&models.Enrollment{}).
		Where("tenant_id = ? AND student_id = ? AND course_id = ?", tenantID, studentID, courseID).
		Count(&count)
	return count > 0
}

func (s *Service) canReview(tenantID, studentID, courseID uint) (bool, error) {
	if !s.isEnrolled(tenantID, studentID, courseID) {
		return false, nil
	}

	var existing int64
	s.db.Model(&models.CourseReview{}).
		Where("tenant_id = ? AND course_id = ? AND student_id = ?", tenantID, courseID, studentID).
		Count(&existing)
	if existing > 0 {
		return false, nil
	}

	opts := progress.LoadOptions(s.db, courseID)
	breakdown := progress.CalcBreakdown(s.db, tenantID, studentID, courseID, opts, false)
	return breakdown.Percent >= 100, nil
}

func encodeTags(tags []string) ([]byte, error) {
	if len(tags) == 0 {
		return nil, nil
	}
	filtered := make([]string, 0, len(tags))
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}
		if _, ok := allowedTags[tag]; !ok {
			return nil, errValidation
		}
		filtered = append(filtered, tag)
	}
	if len(filtered) == 0 {
		return nil, nil
	}
	return json.Marshal(filtered)
}

func toReviewResponse(review models.CourseReview) response.CourseReviewResponse {
	var tags []string
	if len(review.Tags) > 0 {
		_ = json.Unmarshal(review.Tags, &tags)
	}

	resp := response.CourseReviewResponse{
		ID:        review.ID,
		CourseID:  review.CourseID,
		StudentID: review.StudentID,
		Rating:    review.Rating,
		Comment:   review.Comment,
		Tags:      tags,
		CreatedAt: review.CreatedAt,
		UpdatedAt: review.UpdatedAt,
	}

	if review.Student.ID > 0 {
		resp.Student = &response.CourseReviewStudent{
			ID:           review.Student.ID,
			FirstName:    review.Student.FirstName,
			LastName:     review.Student.LastName,
			ProfileImage: review.Student.ProfileImage,
		}
	}

	return resp
}

var (
	errCourseNotFound     = errors.New("course not found")
	errNotEnrolled        = errors.New("not enrolled")
	errCourseNotCompleted = errors.New("course not completed")
	errReviewExists       = errors.New("review already exists")
	errValidation         = errors.New("validation error")
)
