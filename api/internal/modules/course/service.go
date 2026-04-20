package course

import (
	"dashlearn/internal/models"
	"dashlearn/internal/response"
	"dashlearn/internal/utils"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type CourseService interface {
	GetAll(tenantID uint) ([]models.CourseDetails, error)
	GetAllLite(tenantID uint) ([]struct {
		ID    uint   `json:"id"`
		Title string `json:"title"`
	}, error)
	SearchCourses(tenantID uint, query string) ([]struct {
		ID            uint   `json:"id"`
		Title         string `json:"title"`
		FeaturedImage string `json:"featured_image"`
		Slug          string `json:"slug"`
	}, error)
	GetAllPublic(tenantID uint, limitApplied bool, showItems int) ([]response.CourseDetailsPublicResponse, error)
	GetAllPublicByCategory(tenantID uint, categorySlug string) ([]response.CourseDetailsPublicResponse, error)
	GetAllPublicBySubCategory(tenantID uint, categorySlug string) ([]response.CourseDetailsPublicResponse, error)
	GetByID(tenantID uint, courseID uint) (models.CourseDetails, error)
	GetBySlugPublic(tenantID uint, slug string) (*response.CourseDetailsPublicResponse, error)
	Create(input CourseDetailsInput, tenantID uint, userID uint) error
	Update(courseID, tenantID, userID uint, input CourseDetailsInput) error
	Delete(id uint, tenantID uint) error
	DeleteLessonResource(courseID uint, resourceID uint) error
	ReorderCourses(tenantID uint, input ReorderRequest) error
}

type courseService struct {
	db *gorm.DB
}

func NewCourseService(db *gorm.DB) CourseService {
	return &courseService{
		db: db,
	}
}

func (s *courseService) GetAll(tenantID uint) ([]models.CourseDetails, error) {
	var courses []models.CourseDetails

	err := s.db.Where("tenant_id = ?", tenantID).Preload("Author").Preload("Chapters").Preload("Chapters.Lessons").Preload("GeneralSettings").Preload("GeneralSettings.Category").Preload("Instructors").Preload("Instructors.Instructor").Order("position ASC").Find(&courses).Error

	return courses, err
}

func (s *courseService) GetAllLite(tenantID uint) ([]struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}, error) {
	var courses []struct {
		ID    uint   `json:"id"`
		Title string `json:"title"`
	}

	err := s.db.Table("course_details").
		Where("tenant_id = ?", tenantID).
		Select("id", "title", "tenant_id").
		Order("position ASC").
		Find(&courses).Error

	return courses, err
}

func (s *courseService) SearchCourses(tenantID uint, query string) ([]struct {
	ID            uint   `json:"id"`
	Title         string `json:"title"`
	FeaturedImage string `json:"featured_image"`
	Slug          string `json:"slug"`
}, error) {
	var courses []struct {
		ID            uint   `json:"id"`
		Title         string `json:"title"`
		FeaturedImage string `json:"featured_image"`
		Slug          string `json:"slug"`
	}

	err := s.db.Table("course_details").
		Where("tenant_id = ?", tenantID).
		Where("title LIKE ? OR JSON_CONTAINS(tags, JSON_QUOTE(?)) = 1", "%"+query+"%", query).
		Select("id", "title", "featured_image", "slug").
		Order("position ASC").
		Find(&courses).Error

	return courses, err
}

func (s *courseService) GetAllPublic(tenantID uint, limitApplied bool, showItems int) ([]response.CourseDetailsPublicResponse, error) {
	var modelCourses []models.CourseDetails
	var publicResponses []response.CourseDetailsPublicResponse

	dbQuery := s.db.
		Where(models.CourseDetails{
			TenantID:   tenantID,
			Visibility: models.Public,
		}).
		Preload("GeneralSettings").
		Preload("GeneralSettings.Category").
		Order("position ASC")

	if limitApplied {
		dbQuery = dbQuery.Limit(showItems)
	}

	err := dbQuery.
		Find(&modelCourses).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []response.CourseDetailsPublicResponse{}, nil
		}
		return nil, fmt.Errorf("failed to retrieve courses: %w", err)
	}

	for _, course := range modelCourses {
		res := response.CourseDetailsPublicResponse{
			ID:          course.ID,
			Title:       course.Title,
			Slug:        course.Slug,
			Summary:     course.Summary,
			Visibility:  course.Visibility,
			IsScheduled: *course.IsScheduled,
			// ScheduleDate:    course.ScheduleDate,
			// ScheduleTime:    course.ScheduleTime,
			FeaturedImage: course.FeaturedImage,
			IntroVideo:    course.IntroVideo,
			PricingModel:  course.PricingModel,
			RegularPrice:  course.RegularPrice,
			SalePrice:     course.SalePrice,
			// ShowCommingSoom: course.ShowCommingSoom,
			Tags: course.Tags,
			GeneralSettings: &response.CourseGeneralSettingsResponse{
				ID:              course.GeneralSettings.ID,
				CourseID:        course.GeneralSettings.CourseID,
				DifficultyLevel: course.GeneralSettings.DifficultyLevel,
				Language:        course.GeneralSettings.Language,
				MaximumStudent:  course.GeneralSettings.MaximumStudent,
				Category: response.CategoryResponse{
					ID:          course.GeneralSettings.Category.ID,
					Name:        course.GeneralSettings.Category.Name,
					Slug:        course.GeneralSettings.Category.Slug,
					Description: utils.EmptyStringToNil(course.GeneralSettings.Category.Description),
					Thumbnail:   utils.EmptyStringToNil(course.GeneralSettings.Category.Thumbnail),
					CreatedAt:   course.GeneralSettings.Category.CreatedAt,
					UpdatedAt:   course.GeneralSettings.Category.UpdatedAt,
				},
				Duration:  course.GeneralSettings.Duration,
				CreatedAt: course.GeneralSettings.CreatedAt,
				UpdatedAt: course.GeneralSettings.UpdatedAt,
			},
		}
		publicResponses = append(publicResponses, res)
	}

	return publicResponses, err
}

func (s *courseService) GetAllPublicByCategory(tenantID uint, categorySlug string) ([]response.CourseDetailsPublicResponse, error) {
	var modelCourses []models.CourseDetails
	var publicResponses []response.CourseDetailsPublicResponse

	err := s.db.
		Joins("JOIN course_general_settings ON course_general_settings.course_id = course_details.id").
		Joins("JOIN categories ON categories.id = course_general_settings.category_id").
		Where("course_details.tenant_id = ? AND course_details.visibility = ?", tenantID, models.Public).
		Where("categories.slug = ?", categorySlug).
		Preload("GeneralSettings").
		Preload("GeneralSettings.Category").
		Order("position ASC").
		Find(&modelCourses).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []response.CourseDetailsPublicResponse{}, nil
		}
		return nil, fmt.Errorf("failed to retrieve courses: %w", err)
	}

	for _, course := range modelCourses {
		res := response.CourseDetailsPublicResponse{
			ID:          course.ID,
			Title:       course.Title,
			Slug:        course.Slug,
			Summary:     course.Summary,
			Visibility:  course.Visibility,
			IsScheduled: *course.IsScheduled,
			// ScheduleDate:    course.ScheduleDate,
			// ScheduleTime:    course.ScheduleTime,
			FeaturedImage: course.FeaturedImage,
			IntroVideo:    course.IntroVideo,
			PricingModel:  course.PricingModel,
			RegularPrice:  course.RegularPrice,
			SalePrice:     course.SalePrice,
			// ShowCommingSoom: course.ShowCommingSoom,
			Tags: course.Tags,
			GeneralSettings: &response.CourseGeneralSettingsResponse{
				ID:              course.GeneralSettings.ID,
				CourseID:        course.GeneralSettings.CourseID,
				DifficultyLevel: course.GeneralSettings.DifficultyLevel,
				Language:        course.GeneralSettings.Language,
				MaximumStudent:  course.GeneralSettings.MaximumStudent,
				Category: response.CategoryResponse{
					ID:          course.GeneralSettings.Category.ID,
					Name:        course.GeneralSettings.Category.Name,
					Slug:        course.GeneralSettings.Category.Slug,
					Description: utils.EmptyStringToNil(course.GeneralSettings.Category.Description),
					Thumbnail:   utils.EmptyStringToNil(course.GeneralSettings.Category.Thumbnail),
					CreatedAt:   course.GeneralSettings.Category.CreatedAt,
					UpdatedAt:   course.GeneralSettings.Category.UpdatedAt,
				},
				Duration:  course.GeneralSettings.Duration,
				CreatedAt: course.GeneralSettings.CreatedAt,
				UpdatedAt: course.GeneralSettings.UpdatedAt,
			},
		}
		publicResponses = append(publicResponses, res)
	}

	return publicResponses, err
}

func (s *courseService) GetAllPublicBySubCategory(tenantID uint, categorySlug string) ([]response.CourseDetailsPublicResponse, error) {
	var modelCourses []models.CourseDetails
	var publicResponses []response.CourseDetailsPublicResponse

	err := s.db.
		Joins("JOIN course_general_settings ON course_general_settings.course_id = course_details.id").
		Joins("JOIN categories ON categories.id = course_general_settings.category_id").
		Joins("JOIN sub_categories ON sub_categories.id = course_general_settings.sub_category_id").
		Where("course_details.tenant_id = ? AND course_details.visibility = ?", tenantID, models.Public).
		Where("sub_categories.slug = ?", categorySlug).
		Preload("GeneralSettings").
		Preload("GeneralSettings.Category").
		Preload("GeneralSettings.SubCategory").
		Order("position ASC").
		Find(&modelCourses).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []response.CourseDetailsPublicResponse{}, nil
		}
		return nil, fmt.Errorf("failed to retrieve courses: %w", err)
	}

	for _, course := range modelCourses {
		res := response.CourseDetailsPublicResponse{
			ID:          course.ID,
			Title:       course.Title,
			Slug:        course.Slug,
			Summary:     course.Summary,
			Visibility:  course.Visibility,
			IsScheduled: *course.IsScheduled,
			// ScheduleDate:    course.ScheduleDate,
			// ScheduleTime:    course.ScheduleTime,
			FeaturedImage: course.FeaturedImage,
			IntroVideo:    course.IntroVideo,
			PricingModel:  course.PricingModel,
			RegularPrice:  course.RegularPrice,
			SalePrice:     course.SalePrice,
			// ShowCommingSoom: course.ShowCommingSoom,
			Tags: course.Tags,
			GeneralSettings: &response.CourseGeneralSettingsResponse{
				ID:              course.GeneralSettings.ID,
				CourseID:        course.GeneralSettings.CourseID,
				DifficultyLevel: course.GeneralSettings.DifficultyLevel,
				Language:        course.GeneralSettings.Language,
				MaximumStudent:  course.GeneralSettings.MaximumStudent,
				Category: response.CategoryResponse{
					ID:          course.GeneralSettings.Category.ID,
					Name:        course.GeneralSettings.Category.Name,
					Slug:        course.GeneralSettings.Category.Slug,
					Description: utils.EmptyStringToNil(course.GeneralSettings.Category.Description),
					Thumbnail:   utils.EmptyStringToNil(course.GeneralSettings.Category.Thumbnail),
					CreatedAt:   course.GeneralSettings.Category.CreatedAt,
					UpdatedAt:   course.GeneralSettings.Category.UpdatedAt,
				},
				Duration:  course.GeneralSettings.Duration,
				CreatedAt: course.GeneralSettings.CreatedAt,
				UpdatedAt: course.GeneralSettings.UpdatedAt,
			},
		}
		publicResponses = append(publicResponses, res)
	}

	return publicResponses, err
}

func (s *courseService) GetByID(tenantID uint, courseID uint) (models.CourseDetails, error) {
	var course models.CourseDetails

	err := s.db.
		Where("tenant_id = ? AND id = ?", tenantID, courseID).
		Preload("Author").
		Preload("Chapters").
		Preload("Chapters.Lessons").
		Preload("Chapters.Lessons.Resources").
		Preload("Chapters.Assignments").
		Preload("Chapters.Quizzes").
		Preload("Chapters.Quizzes.Questions").
		Preload("GeneralSettings").
		Preload("GeneralSettings.Category").
		Preload("Instructors").
		Preload("Instructors.Instructor").
		Preload("Enrollments").
		First(&course).Error

	return course, err
}

// func (s *courseService) GetBySlugPublic(tenantID uint, slug string) (*response.CourseDetailsPublicResponse, error) {
// 	var modelCourse models.CourseDetails

// 	err := s.db.
// 		Where("tenant_id = ? AND slug = ?", tenantID, slug).
// 		Preload("Author").
// 		Preload("Chapters", "access = 'published'").
// 		Preload("Chapters.Lessons", "is_published = true").
// 		Preload("Chapters.Lessons.Resources").
// 		Preload("Chapters.Assignments", "is_published = true").
// 		Preload("Chapters.Quizzes", "is_published = true").
// 		Preload("Chapters.Quizzes.Questions").
// 		Preload("GeneralSettings").
// 		Preload("GeneralSettings.Category").
// 		Preload("Instructors").
// 		Preload("Instructors.Instructor").
// 		Preload("Enrollments").
// 		First(&modelCourse).Error

// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return &response.CourseDetailsPublicResponse{}, nil
// 		}
// 		return nil, fmt.Errorf("failed to retrieve course: %w", err)
// 	}

// 	chapters := make([]response.CourseChapterResponse, len(modelCourse.Chapters))
// 	for i, chapter := range modelCourse.Chapters {
// 		lessons := make([]response.CourseLessonResponse, len(chapter.Lessons))
// 		for j, lesson := range chapter.Lessons {
// 			lessons[j] = response.CourseLessonResponse{
// 				ID:          lesson.ID,
// 				Title:       lesson.Title,
// 				Description: lesson.Description,
// 				Position:    lesson.Position,
// 				CreatedAt:   lesson.CreatedAt,
// 				UpdatedAt:   lesson.UpdatedAt,
// 				ChapterID:   lesson.ChapterID,
// 				LessonType:  lesson.LessonType,
// 				SourceType:  lesson.SourceType,
// 				Source:      lesson.Source,
// 				IsPublic:    lesson.IsPublic,
// 				Resources:   lesson.Resources,
// 			}
// 		}

// 		assignments := make([]response.CourseAssignmentResponse, len(chapter.Assignments))
// 		for j, assignment := range chapter.Assignments {
// 			assignments[j] = response.CourseAssignmentResponse{
// 				ID:               assignment.ID,
// 				ChapterID:        assignment.ChapterID,
// 				CourseID:         assignment.CourseID,
// 				Title:            assignment.Title,
// 				Instructions:     assignment.Instructions,
// 				Attachments:      assignment.Attachments,
// 				IsPublished:      assignment.IsPublished,
// 				TimeLimit:        assignment.TimeLimit,
// 				TimeLimitOption:  assignment.TimeLimitOption,
// 				FileUploadLimit:  assignment.FileUploadLimit,
// 				TotalMarks:       assignment.TotalMarks,
// 				MinimumPassMarks: assignment.MinimumPassMarks,
// 				CreatedAt:        assignment.CreatedAt,
// 				UpdatedAt:        assignment.UpdatedAt,
// 			}
// 		}

// 		quizzes := make([]response.CourseQuizResponse, len(chapter.Quizzes))
// 		for j, quiz := range chapter.Quizzes {
// 			questions := make([]response.CourseQuizQuestionsResponse, len(quiz.Questions))
// 			for k, question := range quiz.Questions {
// 				questions[k] = response.CourseQuizQuestionsResponse{
// 					ID:                question.ID,
// 					QuizID:            question.QuizID,
// 					Title:             question.Title,
// 					Details:           question.Details,
// 					Media:             question.Media,
// 					Type:              question.Type,
// 					Marks:             question.Marks,
// 					AnswerRequired:    question.AnswerRequired,
// 					AnswerExplanation: question.AnswerExplanation,
// 					CreatedAt:         question.CreatedAt,
// 					UpdatedAt:         question.UpdatedAt,
// 				}
// 			}
// 			quizzes[j] = response.CourseQuizResponse{
// 				ID:                    quiz.ID,
// 				ChapterID:             quiz.ChapterID,
// 				CourseID:              quiz.CourseID,
// 				Title:                 quiz.Title,
// 				Instructions:          quiz.Instructions,
// 				IsPublished:           quiz.IsPublished,
// 				RandomizeQuestions:    quiz.RandomizeQuestions,
// 				SingleQuizView:        quiz.SingleQuizView,
// 				TimeLimit:             quiz.TimeLimit,
// 				TimeLimitOption:       quiz.TimeLimitOption,
// 				TotalVisibleQuestions: quiz.TotalVisibleQuestions,
// 				RevealAnswers:         quiz.RevealAnswers,
// 				EnableRetry:           quiz.EnableRetry,
// 				RetryAttempts:         quiz.RetryAttempts,
// 				MinimumPassPercentage: quiz.MinimumPassPercentage,
// 				Questions:             questions,
// 				CreatedAt:             quiz.CreatedAt,
// 				UpdatedAt:             quiz.UpdatedAt,
// 			}
// 		}

// 		chapters[i] = response.CourseChapterResponse{
// 			ID:          chapter.ID,
// 			Title:       chapter.Title,
// 			Description: chapter.Description,
// 			Position:    chapter.Position,
// 			Access:      chapter.Access,
// 			CreatedAt:   chapter.CreatedAt,
// 			UpdatedAt:   chapter.UpdatedAt,
// 			CourseID:    chapter.CourseID,
// 			Lessons:     lessons,
// 			Assignments: assignments,
// 			Quizzes:     quizzes,
// 		}
// 	}

// 	instructors := make([]response.CourseInstructorResponse, len(modelCourse.Instructors))
// 	for i, instructor := range modelCourse.Instructors {
// 		instructors[i] = response.CourseInstructorResponse{
// 			ID:           instructor.ID,
// 			CourseID:     instructor.CourseID,
// 			InstructorID: instructor.InstructorID,
// 			Instructor: response.InstructorResponse{
// 				ID:          instructor.Instructor.ID,
// 				FirstName:   instructor.Instructor.FirstName,
// 				LastName:    instructor.Instructor.LastName,
// 				Email:       instructor.Instructor.Email,
// 				Image:       utils.ZeroToNil(instructor.Instructor.Image),
// 				Phone:       instructor.Instructor.Phone,
// 				Role:        instructor.Instructor.Role,
// 				Designation: instructor.Instructor.Designation,
// 			},
// 		}
// 	}

// 	enrollments := make([]response.EnrolledCourseRes, len(modelCourse.Enrollments))
// 	for i, enrollment := range modelCourse.Enrollments {
// 		enrollments[i] = response.EnrolledCourseRes{
// 			ID:        enrollment.ID,
// 			CourseID:  enrollment.CourseID,
// 			StudentID: enrollment.StudentID,
// 		}
// 	}

// 	res := &response.CourseDetailsPublicResponse{
// 		ID:          modelCourse.ID,
// 		Title:       modelCourse.Title,
// 		Summary:     modelCourse.Summary,
// 		Description: modelCourse.Description,
// 		Visibility:  modelCourse.Visibility,
// 		IsScheduled: *modelCourse.IsScheduled,
// 		// ScheduleDate:    modelCourse.ScheduleDate,
// 		// ScheduleTime:    modelCourse.ScheduleTime,
// 		FeaturedImage: modelCourse.FeaturedImage,
// 		IntroVideo:    modelCourse.IntroVideo,
// 		PricingModel:  modelCourse.PricingModel,
// 		RegularPrice:  modelCourse.RegularPrice,
// 		SalePrice:     modelCourse.SalePrice,
// 		// ShowCommingSoom: modelCourse.ShowCommingSoom,
// 		Tags:     modelCourse.Tags,
// 		Overview: modelCourse.Overview,
// 		GeneralSettings: &response.CourseGeneralSettingsResponse{
// 			ID:              modelCourse.GeneralSettings.ID,
// 			CourseID:        modelCourse.GeneralSettings.CourseID,
// 			DifficultyLevel: modelCourse.GeneralSettings.DifficultyLevel,
// 			Language:        modelCourse.GeneralSettings.Language,
// 			MaximumStudent:  modelCourse.GeneralSettings.MaximumStudent,
// 			Category: response.CategoryResponse{
// 				ID:          modelCourse.GeneralSettings.Category.ID,
// 				Name:        modelCourse.GeneralSettings.Category.Name,
// 				Slug:        modelCourse.GeneralSettings.Category.Slug,
// 				Description: utils.EmptyStringToNil(modelCourse.GeneralSettings.Category.Description),
// 				Thumbnail:   utils.EmptyStringToNil(modelCourse.GeneralSettings.Category.Thumbnail),
// 				CreatedAt:   modelCourse.GeneralSettings.Category.CreatedAt,
// 				UpdatedAt:   modelCourse.GeneralSettings.Category.UpdatedAt,
// 			},
// 			Duration:  modelCourse.GeneralSettings.Duration,
// 			CreatedAt: modelCourse.GeneralSettings.CreatedAt,
// 			UpdatedAt: modelCourse.GeneralSettings.UpdatedAt,
// 		},
// 		Chapters:    chapters,
// 		Instructors: instructors,
// 		Enrollments: enrollments,
// 	}

// 	return res, err
// }

func (s *courseService) GetBySlugPublic(tenantID uint, slug string) (*response.CourseDetailsPublicResponse, error) {
	var modelCourse models.CourseDetails

	err := s.db.
		Where("tenant_id = ? AND slug = ?", tenantID, slug).
		Preload("Author").
		Preload("Chapters", "access = 'published'").
		Preload("Chapters.Lessons", "is_published = true").
		Preload("Chapters.Lessons.Resources").
		Preload("Chapters.Assignments", "is_published = true").
		Preload("Chapters.Quizzes", "is_published = true").
		Preload("Chapters.Quizzes.Questions").
		Preload("GeneralSettings").
		Preload("GeneralSettings.Category").
		Preload("Instructors").
		Preload("Instructors.Instructor").
		Preload("Enrollments").
		First(&modelCourse).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &response.CourseDetailsPublicResponse{}, nil
		}
		return nil, fmt.Errorf("failed to retrieve course: %w", err)
	}

	// Build chapters with filtering for empty ones
	chapters := make([]response.CourseChapterResponse, 0, len(modelCourse.Chapters))
	for _, chapter := range modelCourse.Chapters {
		lessons := make([]response.CourseLessonResponse, 0, len(chapter.Lessons))
		for _, lesson := range chapter.Lessons {
			lessons = append(lessons, response.CourseLessonResponse{
				ID:          lesson.ID,
				Title:       lesson.Title,
				Description: lesson.Description,
				Position:    lesson.Position,
				CreatedAt:   lesson.CreatedAt,
				UpdatedAt:   lesson.UpdatedAt,
				ChapterID:   lesson.ChapterID,
				LessonType:  lesson.LessonType,
				SourceType:  lesson.SourceType,
				Source:      lesson.Source,
				IsPublic:    lesson.IsPublic,
				Resources:   lesson.Resources,
			})
		}

		assignments := make([]response.CourseAssignmentResponse, 0, len(chapter.Assignments))
		for _, assignment := range chapter.Assignments {
			assignments = append(assignments, response.CourseAssignmentResponse{
				ID:               assignment.ID,
				ChapterID:        assignment.ChapterID,
				CourseID:         assignment.CourseID,
				Title:            assignment.Title,
				Instructions:     assignment.Instructions,
				Attachments:      assignment.Attachments,
				IsPublished:      assignment.IsPublished,
				TimeLimit:        assignment.TimeLimit,
				TimeLimitOption:  assignment.TimeLimitOption,
				FileUploadLimit:  assignment.FileUploadLimit,
				TotalMarks:       assignment.TotalMarks,
				MinimumPassMarks: assignment.MinimumPassMarks,
				CreatedAt:        assignment.CreatedAt,
				UpdatedAt:        assignment.UpdatedAt,
			})
		}

		quizzes := make([]response.CourseQuizResponse, 0, len(chapter.Quizzes))
		for _, quiz := range chapter.Quizzes {
			questions := make([]response.CourseQuizQuestionsResponse, 0, len(quiz.Questions))
			for _, question := range quiz.Questions {
				questions = append(questions, response.CourseQuizQuestionsResponse{
					ID:                question.ID,
					QuizID:            question.QuizID,
					Title:             question.Title,
					Details:           question.Details,
					Media:             question.Media,
					Type:              question.Type,
					Marks:             question.Marks,
					AnswerRequired:    question.AnswerRequired,
					AnswerExplanation: question.AnswerExplanation,
					CreatedAt:         question.CreatedAt,
					UpdatedAt:         question.UpdatedAt,
				})
			}

			// Only include quiz if it has questions or is published
			if len(questions) > 0 || quiz.IsPublished {
				quizzes = append(quizzes, response.CourseQuizResponse{
					ID:                    quiz.ID,
					ChapterID:             quiz.ChapterID,
					CourseID:              quiz.CourseID,
					Title:                 quiz.Title,
					Instructions:          quiz.Instructions,
					IsPublished:           quiz.IsPublished,
					RandomizeQuestions:    quiz.RandomizeQuestions,
					SingleQuizView:        quiz.SingleQuizView,
					TimeLimit:             quiz.TimeLimit,
					TimeLimitOption:       quiz.TimeLimitOption,
					TotalVisibleQuestions: quiz.TotalVisibleQuestions,
					RevealAnswers:         quiz.RevealAnswers,
					EnableRetry:           quiz.EnableRetry,
					RetryAttempts:         quiz.RetryAttempts,
					MinimumPassPercentage: quiz.MinimumPassPercentage,
					Questions:             questions,
					CreatedAt:             quiz.CreatedAt,
					UpdatedAt:             quiz.UpdatedAt,
				})
			}
		}

		// Only include chapter if it has at least one lesson, assignment, or quiz
		if len(lessons) > 0 || len(assignments) > 0 || len(quizzes) > 0 {
			chapters = append(chapters, response.CourseChapterResponse{
				ID:          chapter.ID,
				Title:       chapter.Title,
				Description: chapter.Description,
				Position:    chapter.Position,
				Access:      chapter.Access,
				CreatedAt:   chapter.CreatedAt,
				UpdatedAt:   chapter.UpdatedAt,
				CourseID:    chapter.CourseID,
				Lessons:     lessons,
				Assignments: assignments,
				Quizzes:     quizzes,
			})
		}
	}

	// Instructors
	instructors := make([]response.CourseInstructorResponse, len(modelCourse.Instructors))
	for i, instructor := range modelCourse.Instructors {
		instructors[i] = response.CourseInstructorResponse{
			ID:           instructor.ID,
			CourseID:     instructor.CourseID,
			InstructorID: instructor.InstructorID,
			Instructor: response.InstructorResponse{
				ID:          instructor.Instructor.ID,
				FirstName:   instructor.Instructor.FirstName,
				LastName:    instructor.Instructor.LastName,
				Email:       instructor.Instructor.Email,
				Image:       utils.ZeroToNil(instructor.Instructor.Image),
				Phone:       instructor.Instructor.Phone,
				Role:        instructor.Instructor.Role,
				Designation: instructor.Instructor.Designation,
			},
		}
	}

	// Enrollments
	enrollments := make([]response.EnrolledCourseRes, len(modelCourse.Enrollments))
	for i, enrollment := range modelCourse.Enrollments {
		enrollments[i] = response.EnrolledCourseRes{
			ID:        enrollment.ID,
			CourseID:  enrollment.CourseID,
			StudentID: enrollment.StudentID,
		}
	}

	// Build final response
	res := &response.CourseDetailsPublicResponse{
		ID:            modelCourse.ID,
		Title:         modelCourse.Title,
		Summary:       modelCourse.Summary,
		Description:   modelCourse.Description,
		Visibility:    modelCourse.Visibility,
		IsScheduled:   *modelCourse.IsScheduled,
		FeaturedImage: modelCourse.FeaturedImage,
		IntroVideo:    modelCourse.IntroVideo,
		PricingModel:  modelCourse.PricingModel,
		RegularPrice:  modelCourse.RegularPrice,
		SalePrice:     modelCourse.SalePrice,
		Tags:          modelCourse.Tags,
		Overview:      modelCourse.Overview,
		GeneralSettings: &response.CourseGeneralSettingsResponse{
			ID:              modelCourse.GeneralSettings.ID,
			CourseID:        modelCourse.GeneralSettings.CourseID,
			DifficultyLevel: modelCourse.GeneralSettings.DifficultyLevel,
			Language:        modelCourse.GeneralSettings.Language,
			MaximumStudent:  modelCourse.GeneralSettings.MaximumStudent,
			Category: response.CategoryResponse{
				ID:          modelCourse.GeneralSettings.Category.ID,
				Name:        modelCourse.GeneralSettings.Category.Name,
				Slug:        modelCourse.GeneralSettings.Category.Slug,
				Description: utils.EmptyStringToNil(modelCourse.GeneralSettings.Category.Description),
				Thumbnail:   utils.EmptyStringToNil(modelCourse.GeneralSettings.Category.Thumbnail),
				CreatedAt:   modelCourse.GeneralSettings.Category.CreatedAt,
				UpdatedAt:   modelCourse.GeneralSettings.Category.UpdatedAt,
			},
			Duration:  modelCourse.GeneralSettings.Duration,
			CreatedAt: modelCourse.GeneralSettings.CreatedAt,
			UpdatedAt: modelCourse.GeneralSettings.UpdatedAt,
		},
		Chapters:    chapters,
		Instructors: instructors,
		Enrollments: enrollments,
	}

	return res, err
}

func (s *courseService) Create(input CourseDetailsInput, tenantID uint, userID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var total int64
		if err := tx.Model(&models.CourseDetails{}).
			Where("tenant_id = ?", tenantID).
			Count(&total).Error; err != nil {
			return err
		}

		var videoPtr *models.IntroVideo

		var lastID uint
		err := tx.Model(&models.CourseDetails{}).
			Select("id").
			Order("id DESC").
			Limit(1).
			Pluck("id", &lastID).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		newID := lastID + 1

		if input.IntroVideo == nil ||
			(input.IntroVideo.Type == "" && input.IntroVideo.Source == "") {
			videoPtr = nil
		} else {
			videoPtr = &models.IntroVideo{
				Type:   input.IntroVideo.Type,
				Source: input.IntroVideo.Source,
			}
		}

		input.IntroVideo = videoPtr

		tagsJSON, err := utils.NormalizeTags(input.Tags)
		if err != nil {
			tagsJSON = nil
		}

		overviewJSON, err := utils.NormalizeTags(input.Overview)
		if err != nil {
			overviewJSON = nil
		}

		var introVideo *utils.JSONB[models.IntroVideo]
		if input.IntroVideo != nil {
			introVideo = &utils.JSONB[models.IntroVideo]{Data: *input.IntroVideo}
		}

		var scheduleDateForDB *string
		var scheduleTimeForDB *string

		if input.IsScheduled == "true" {
			cleanDate, err := NormalizeDate(*input.ScheduleDate)
			if err != nil {
				return errors.New("invalid schedule date format")
			}
			scheduleDateForDB = &cleanDate

			timeParsed, err := time.Parse("03:04 PM", *input.ScheduleTime)
			if err != nil {
				return errors.New("invalid schedule time format")
			}
			timeStr := timeParsed.Format("15:04:05")
			scheduleTimeForDB = &timeStr
		}

		newCourseDetails := models.CourseDetails{
			Title:       input.Title,
			Slug:        utils.Slugify(input.Title) + "-" + strconv.Itoa(int(newID)),
			Summary:     input.Summary,
			Description: utils.ZeroToNil(input.Description),
			Visibility:  input.Visibility,
			IsScheduled: func() *bool {
				if input.IsScheduled == "true" {
					return &[]bool{true}[0]
				}
				return &[]bool{false}[0]
			}(),
			ScheduleDate: scheduleDateForDB,
			ScheduleTime: scheduleTimeForDB,
			PricingModel: input.PricingModel,
			RegularPrice: input.RegularPrice,
			SalePrice:    input.SalePrice,
			ShowCommingSoon: func() *bool {
				if input.ShowCommingSoon != nil {
					b := false
					if *input.ShowCommingSoon == "true" {
						b = true
					}
					return &b
				}
				b := false
				return &b
			}(),
			Tags:          tagsJSON,
			Overview:      overviewJSON,
			FeaturedImage: input.FeaturedImage,
			IntroVideo:    introVideo,
			AuthorID:      userID,
			TenantID:      tenantID,
			Position:      total,
		}

		if err := tx.Create(&newCourseDetails).Error; err != nil {
			return err
		}

		// Create chapters and lessons
		for idx, chapter := range input.CourseChapters {
			newCourseChapter := models.CourseChapter{
				CourseID:    newCourseDetails.ID,
				Title:       chapter.Title,
				Description: utils.EmptyStringToNil(chapter.Description),
				Position:    idx,
				Access:      chapter.Access,
			}

			if err := tx.Create(&newCourseChapter).Error; err != nil {
				return err
			}

			// Create lessons
			for idx, lesson := range chapter.CourseLessons {
				sourceJSON := utils.JSONB[models.Source]{Data: lesson.Source}

				var lessonDateForDB *string
				var lessonTimeForDB *string

				if lesson.IsScheduled {
					if lesson.ScheduleDate == nil || strings.TrimSpace(*lesson.ScheduleDate) == "" ||
						lesson.ScheduleTime == nil || strings.TrimSpace(*lesson.ScheduleTime) == "" {
						return errors.New("scheduled lesson requires schedule_date and schedule_time")
					}
					cleanDate, err := NormalizeDate(strings.TrimSpace(*lesson.ScheduleDate))
					if err != nil {
						return errors.New("invalid schedule date format")
					}
					lessonDateForDB = &cleanDate

					timeParsed, err := time.Parse("03:04 PM", strings.TrimSpace(*lesson.ScheduleTime))
					if err != nil {
						return errors.New("invalid schedule time format")
					}
					timeStr := timeParsed.Format("15:04:05")
					lessonTimeForDB = &timeStr
				}

				newCourseLesson := models.CourseLesson{
					ChapterID:   newCourseChapter.ID,
					Title:       lesson.Title,
					Description: utils.EmptyStringToNil(lesson.Description),
					Position:    idx,
					LessonType:  lesson.LessonType,
					SourceType:  lesson.SourceType,
					Source:      sourceJSON,
					IsPublished: func() bool {
						if lesson.IsScheduled {
							return false
						}
						return lesson.IsPublished
					}(),
					IsPublic:        lesson.IsPublic,
					IsScheduled:     &lesson.IsScheduled,
					ScheduleDate:    lessonDateForDB,
					ScheduleTime:    lessonTimeForDB,
					ShowCommingSoon: lesson.ShowCommingSoon,
				}

				if err := tx.Create(&newCourseLesson).Error; err != nil {
					return err
				}

				if len(lesson.Resources) > 0 {
					for _, res := range lesson.Resources {
						newResource := models.LessonResource{
							CourseID: newCourseDetails.ID,
							LessonID: newCourseLesson.ID,
							MimeType: res.MimeType,
							Title:    res.FileName,
							FilePath: res.URL,
							FileSize: res.Size,
							Position: res.Position,
						}
						if err := tx.Create(&newResource).Error; err != nil {
							return err
						}
					}
				}

			}

			// course quizes
			for _, quiz := range chapter.Quizzes {
				newCourseQuiz := models.CourseQuiz{
					CourseID:              newCourseDetails.ID,
					ChapterID:             newCourseChapter.ID,
					Title:                 quiz.Title,
					Instructions:          quiz.Instructions,
					IsPublished:           quiz.IsPublished,
					RandomizeQuestions:    quiz.RandomizeQuestions,
					SingleQuizView:        quiz.SingleQuizView,
					TimeLimit:             quiz.TimeLimit,
					TimeLimitOption:       quiz.TimeLimitOption,
					TotalVisibleQuestions: utils.ZeroToNil(quiz.TotalVisibleQuestions),
					RevealAnswers:         quiz.RevealAnswers,
					EnableRetry:           quiz.EnableRetry,
					RetryAttempts:         quiz.RetryAttempts,
					MinimumPassPercentage: quiz.MinimumPassPercentage,
				}

				if err := tx.Create(&newCourseQuiz).Error; err != nil {
					return err
				}

				for _, question := range quiz.Questions {
					newQuizQuestion := models.QuizQuestion{
						QuizID:            newCourseQuiz.ID,
						Title:             question.Title,
						Details:           utils.ZeroToNil(question.Details),
						Media:             utils.ZeroToNil(question.Media),
						Type:              question.Type,
						Marks:             question.Marks,
						AnswerRequired:    question.AnswerRequired,
						AnswerExplanation: utils.ZeroToNil(question.AnswerExplanation),
					}

					if err := tx.Create(&newQuizQuestion).Error; err != nil {
						return err
					}
				}

			}

			// course assignments
			for _, assignment := range chapter.Assignments {
				newAssignment := models.CourseAssignment{
					CourseID:         newCourseDetails.ID,
					ChapterID:        newCourseChapter.ID,
					Title:            assignment.Title,
					Instructions:     assignment.Instructions,
					IsPublished:      assignment.IsPublished,
					TimeLimit:        assignment.TimeLimit,
					TimeLimitOption:  assignment.TimeLimitOption,
					Attachments:      utils.ZeroToNil(assignment.Attachments),
					FileUploadLimit:  assignment.FileUploadLimit,
					TotalMarks:       assignment.TotalMarks,
					MinimumPassMarks: assignment.MinimumPassMarks,
				}

				if err := tx.Create(&newAssignment).Error; err != nil {
					return err
				}
			}
		}

		// Create instructors
		for _, instructor := range input.Instructors {
			newCourseInstructor := models.CourseInstructor{
				CourseID:     newCourseDetails.ID,
				InstructorID: uint(instructor),
			}
			if err := tx.Create(&newCourseInstructor).Error; err != nil {
				return err
			}
		}

		// Create general settings
		var difficultyLevelPtr *models.DifficultyLevel
		if input.GeneralSettings.DifficultyLevel != "" {
			difficultyLevelPtr = &input.GeneralSettings.DifficultyLevel
		} else {
			defaultVal := models.All
			difficultyLevelPtr = &defaultVal
		}

		defaultLang := "english"

		newGeneralSettings := models.CourseGeneralSettings{
			CourseID:        newCourseDetails.ID,
			DifficultyLevel: difficultyLevelPtr,
			MaximumStudent:  utils.ZeroToNil(input.GeneralSettings.MaximumStudent),
			Language:        &defaultLang,
			CategoryID:      input.GeneralSettings.CategoryID,
			SubCategoryID:   utils.ZeroToNil(input.GeneralSettings.SubCategoryID),
			Duration:        utils.ZeroToNil(input.GeneralSettings.Duration),
		}

		if err := tx.Create(&newGeneralSettings).Error; err != nil {
			return err
		}

		// ✅ Everything succeeded, transaction will be committed
		return nil
	})
}

func (s *courseService) Update(courseID, tenantID, userID uint, input CourseDetailsInput) error {
	// Fetch the existing course
	var existing models.CourseDetails
	if err := s.db.Where("id = ? AND tenant_id = ?", courseID, tenantID).First(&existing).Error; err != nil {
		return err
	}

	if input.FeaturedImage != nil {
		if *input.FeaturedImage != "" { // new image uploaded
			if existing.FeaturedImage != nil && *existing.FeaturedImage != "" {
				if delErr := utils.DeleteFromBunny(*existing.FeaturedImage); delErr != nil {
					fmt.Println("Failed to delete old file:", delErr)
				}
			}
		} else {
			// input.FeaturedImage is empty string -> user cleared the image
			if existing.FeaturedImage != nil && *existing.FeaturedImage != "" {
				if delErr := utils.DeleteFromBunny(*existing.FeaturedImage); delErr != nil {
					fmt.Println("Failed to delete old file:", delErr)
				}
			}
			input.FeaturedImage = nil // ensure DB gets null
		}
	} else {
		// input.FeaturedImage is nil -> user didn't touch image, keep existing
		input.FeaturedImage = existing.FeaturedImage
	}

	// Normalize & assign values
	tagsJSON, _ := utils.NormalizeTags(input.Tags)
	overviewJSON, _ := utils.NormalizeTags(input.Overview)

	var introVideo *utils.JSONB[models.IntroVideo]
	if input.IntroVideo != nil {
		introVideo = &utils.JSONB[models.IntroVideo]{Data: *input.IntroVideo}
	}

	var scheduleDateForDB *string
	var scheduleTimeForDB *string
	visibility := input.Visibility
	if input.IsScheduled == "true" {
		if input.ScheduleDate != nil && *input.ScheduleDate != "" {
			cleanDate, err := NormalizeDate(*input.ScheduleDate)
			if err != nil {
				return errors.New("invalid schedule date format")
			}
			scheduleDateForDB = &cleanDate
		}

		if input.ScheduleTime != nil && *input.ScheduleTime != "" {
			timeParsed, err := time.Parse("03:04 PM", *input.ScheduleTime)
			if err != nil {
				return errors.New("invalid schedule time format")
			}
			timeStr := timeParsed.Format("15:04:05")
			scheduleTimeForDB = &timeStr
		}

		visibility = models.Protected
	} else {
		scheduleDateForDB = nil
		scheduleTimeForDB = nil
	}

	// Build update map (use concrete values, not pointers, so GORM gets clean data)
	updates := map[string]interface{}{
		"title":             input.Title,
		"slug":              utils.Slugify(input.Title) + "-" + strconv.Itoa(int(courseID)),
		"summary":           input.Summary,
		"description":       utils.ZeroToNil(input.Description),
		"visibility":        visibility,
		"is_scheduled":      input.IsScheduled == "true",
		"pricing_model":     input.PricingModel,
		"regular_price":     input.RegularPrice,
		"sale_price":        input.SalePrice,
		"show_comming_soon": input.ShowCommingSoon != nil && *input.ShowCommingSoon == "true",
		"tags":              tagsJSON,
		"overview":          overviewJSON,
		"featured_image":    input.FeaturedImage,
		"intro_video":       introVideo,
		"author_id":         userID,
		"tenant_id":         tenantID,
		"schedule_date":     scheduleDateForDB,
		"schedule_time":     scheduleTimeForDB,
	}

	// Apply course-level updates
	if err := s.db.Model(&models.CourseDetails{}).Where("id = ?", courseID).Updates(updates).Error; err != nil {
		return err
	}

	// Get existing chapters and lessons from DB
	var existingChaptersForUpdate []models.CourseChapter
	s.db.Preload("Lessons").Where("course_id = ?", courseID).Find(&existingChaptersForUpdate)

	// Map of existing lesson IDs for quick lookup
	existingLessonMap := make(map[uint]models.CourseLesson)
	for _, chapter := range existingChaptersForUpdate {
		for _, lesson := range chapter.Lessons {
			existingLessonMap[lesson.ID] = lesson
		}
	}

	// Map of existing assignment IDs for quick lookup
	existingAssignmentMap := make(map[uint]models.CourseAssignment)
	s.db.Preload("Assignments").Where("course_id = ?", courseID).Find(&existingChaptersForUpdate)
	for _, chapter := range existingChaptersForUpdate {
		for _, assignment := range chapter.Assignments {
			existingAssignmentMap[assignment.ID] = assignment
		}
	}

	// Map of existing quiz IDs for quick lookup
	existingQuizMap := make(map[uint]models.CourseQuiz)
	s.db.Preload("Quizzes").Preload("Quizzes.Questions").Where("course_id = ?", courseID).Find(&existingChaptersForUpdate)
	for _, chapter := range existingChaptersForUpdate {
		for _, quiz := range chapter.Quizzes {
			existingQuizMap[quiz.ID] = quiz
		}
	}

	// Maps to track incoming IDs
	incomingChapterIDs := make(map[uint]bool)
	incomingLessonIDs := make(map[uint]bool)
	incomingAssignmentIDs := make(map[uint]bool)
	incomingQuizIDs := make(map[uint]bool)
	incomingQuizQuestionIDs := make(map[uint]bool)

	// Fetch all existing chapters and their lessons
	var existingChapters []models.CourseChapter
	s.db.Preload("Lessons").
		Preload("Assignments").
		Preload("Quizzes").
		Preload("Quizzes.Questions").
		Where("course_id = ?", courseID).
		Find(&existingChapters)

	chapterMap := make(map[uint]models.CourseChapter)
	lessonMap := make(map[uint]models.CourseLesson)
	assignmentMap := make(map[uint]models.CourseAssignment)
	quizMap := make(map[uint]models.CourseQuiz)
	quizQuestionMap := make(map[uint]models.QuizQuestion)

	for _, ch := range existingChapters {
		chapterMap[ch.ID] = ch
		// Map existing lessons
		for _, lesson := range ch.Lessons {
			lessonMap[lesson.ID] = lesson
		}
		// Map existing assignments
		for _, assignment := range ch.Assignments {
			assignmentMap[assignment.ID] = assignment
		}
		// Map existing quizzes
		for _, quiz := range ch.Quizzes {
			quizMap[quiz.ID] = quiz
			// Map existing quiz questions
			for _, question := range quiz.Questions {
				quizQuestionMap[question.ID] = question
			}
		}
	}

	for chIdx, chapter := range input.CourseChapters {
		fmt.Println("🔄 Processing chapter:", chapter.Title)
		// Handle chapter update or create
		var chapterID uint
		if chapter.ID != nil && *chapter.ID != 0 {
			chapterID = uint(*chapter.ID)
			incomingChapterIDs[chapterID] = true

			if existingCh, found := chapterMap[chapterID]; found {
				// Update
				existingCh.Title = chapter.Title
				existingCh.Description = utils.EmptyStringToNil(chapter.Description)
				existingCh.Position = chIdx
				existingCh.Access = chapter.Access

				if err := s.db.Save(&existingCh).Error; err != nil {
					return err
				}
			}
		} else {
			// Create new chapter
			newChapter := models.CourseChapter{
				CourseID:    courseID,
				Title:       chapter.Title,
				Description: utils.EmptyStringToNil(chapter.Description),
				Position:    chIdx,
				Access:      chapter.Access,
			}
			if err := s.db.Create(&newChapter).Error; err != nil {
				return err
			}
			chapterID = newChapter.ID
			incomingChapterIDs[chapterID] = true
		}

		// Handle lessons inside chapter
		for lIdx, lesson := range chapter.CourseLessons {

			fmt.Println("🔄 Processing lesson:", lesson.Title)

			sourceJSON := utils.JSONB[models.Source]{Data: lesson.Source}

			var lessonDateForDB *string
			var lessonTimeForDB *string

			fmt.Println("✅ Processing lesson:", lesson.IsScheduled)

			// Normalize lesson date/time if scheduled (must not deref nil pointers)
			if lesson.IsScheduled {
				if lesson.ScheduleDate == nil || strings.TrimSpace(*lesson.ScheduleDate) == "" ||
					lesson.ScheduleTime == nil || strings.TrimSpace(*lesson.ScheduleTime) == "" {
					return errors.New("scheduled lesson requires schedule_date and schedule_time")
				}
				cleanDate, err := NormalizeDate(strings.TrimSpace(*lesson.ScheduleDate))
				if err != nil {
					return errors.New("invalid lesson schedule date format")
				}
				lessonDateForDB = &cleanDate
				fmt.Println("🎉 lesson cleaned date:", cleanDate)
				tp, err := time.Parse("03:04 PM", strings.TrimSpace(*lesson.ScheduleTime))
				if err != nil {
					return errors.New("invalid lesson schedule time format")
				}
				ts := tp.Format("15:04:05")
				lessonTimeForDB = &ts
			} else {
				lessonDateForDB = nil
				lessonTimeForDB = nil
			}

			// If lesson exists -> update
			if lesson.ID != nil && *lesson.ID != 0 {
				lessonID := uint(*lesson.ID)
				incomingLessonIDs[lessonID] = true

				if existingLesson, ok := lessonMap[lessonID]; ok {
					// Overwrite existingLesson fields explicitly with cleaned values
					existingLesson.Title = lesson.Title
					existingLesson.Description = utils.EmptyStringToNil(lesson.Description)
					existingLesson.Position = lIdx
					existingLesson.LessonType = lesson.LessonType
					existingLesson.SourceType = lesson.SourceType
					existingLesson.Source = sourceJSON
					if lesson.IsScheduled {
						existingLesson.IsPublished = false
					} else {
						existingLesson.IsPublished = lesson.IsPublished
					}
					existingLesson.IsPublic = lesson.IsPublic
					existingLesson.ChapterID = chapterID
					existingLesson.IsScheduled = &lesson.IsScheduled
					// CRITICAL: assign cleaned pointers (may be nil)
					existingLesson.ScheduleDate = lessonDateForDB
					existingLesson.ScheduleTime = lessonTimeForDB
					existingLesson.ShowCommingSoon = lesson.ShowCommingSoon

					// Save existing lesson
					if err := s.db.Save(&existingLesson).Error; err != nil {
						return err
					}

					// resources (unchanged)
					for _, res := range lesson.Resources {
						if res.FileName == "" {
							continue
						}
						newResource := models.LessonResource{
							CourseID: courseID,
							LessonID: existingLesson.ID,
							MimeType: res.MimeType,
							Title:    res.FileName,
							FilePath: res.URL,
							Position: res.Position,
							FileSize: res.Size,
						}
						if err := s.db.Create(&newResource).Error; err != nil {
							return err
						}
					}
				}
			} else {
				// Create new lesson - use cleaned fields (lessonDateForDB / lessonTimeForDB)
				newLesson := models.CourseLesson{
					ChapterID:       chapterID,
					Title:           lesson.Title,
					Description:     utils.EmptyStringToNil(lesson.Description),
					Position:        lIdx,
					LessonType:      lesson.LessonType,
					SourceType:      lesson.SourceType,
					Source:          sourceJSON,
					IsPublished:     !lesson.IsScheduled,
					IsPublic:        lesson.IsPublic,
					IsScheduled:     &lesson.IsScheduled,
					ScheduleDate:    lessonDateForDB,
					ScheduleTime:    lessonTimeForDB,
					ShowCommingSoon: lesson.ShowCommingSoon,
				}
				if err := s.db.Create(&newLesson).Error; err != nil {
					return err
				}

				// resources
				for _, res := range lesson.Resources {
					if res.FileName == "" {
						continue
					}
					newResource := models.LessonResource{
						CourseID: courseID,
						LessonID: newLesson.ID,
						MimeType: res.MimeType,
						Title:    res.FileName,
						FilePath: res.URL,
						Position: res.Position,
						FileSize: res.Size,
					}
					if err := s.db.Create(&newResource).Error; err != nil {
						return err
					}
				}

				incomingLessonIDs[newLesson.ID] = true
			}
		} // end lessons loop

		// Handle assignments inside chapter
		for _, assignment := range chapter.Assignments {
			if assignment.ID != nil && *assignment.ID != 0 {
				assignmentID := uint(*assignment.ID)
				incomingAssignmentIDs[assignmentID] = true

				if existingAssignment, found := assignmentMap[assignmentID]; found {
					// Update
					existingAssignment.Title = assignment.Title
					existingAssignment.Instructions = assignment.Instructions
					// existingAssignment.Position = lIdx
					existingAssignment.Attachments = assignment.Attachments
					existingAssignment.IsPublished = assignment.IsPublished
					existingAssignment.TimeLimit = assignment.TimeLimit
					existingAssignment.TimeLimitOption = assignment.TimeLimitOption
					existingAssignment.FileUploadLimit = assignment.FileUploadLimit
					existingAssignment.TotalMarks = assignment.TotalMarks
					existingAssignment.MinimumPassMarks = assignment.MinimumPassMarks

					if err := s.db.Save(&existingAssignment).Error; err != nil {
						return err
					}
				}
			} else {
				// Create new assignment
				newAssignment := models.CourseAssignment{
					ChapterID:        chapterID,
					CourseID:         courseID,
					Title:            assignment.Title,
					Instructions:     assignment.Instructions,
					IsPublished:      assignment.IsPublished,
					TimeLimit:        assignment.TimeLimit,
					TimeLimitOption:  assignment.TimeLimitOption,
					FileUploadLimit:  assignment.FileUploadLimit,
					TotalMarks:       assignment.TotalMarks,
					MinimumPassMarks: assignment.MinimumPassMarks,
					Attachments:      nil,
				}
				if err := s.db.Create(&newAssignment).Error; err != nil {
					return err
				}
				incomingAssignmentIDs[newAssignment.ID] = true
			}
		}

		// Handle quiz inside chapter
		for _, quiz := range chapter.Quizzes {
			if quiz.ID != nil && *quiz.ID != 0 {
				quizID := uint(*quiz.ID)
				incomingQuizIDs[quizID] = true

				if existingQuiz, found := quizMap[quizID]; found {
					// Update
					existingQuiz.Title = quiz.Title
					existingQuiz.Instructions = quiz.Instructions
					// existingQuiz.Position = lIdx
					existingQuiz.IsPublished = quiz.IsPublished
					existingQuiz.TimeLimit = quiz.TimeLimit
					existingQuiz.TimeLimitOption = quiz.TimeLimitOption
					existingQuiz.TotalVisibleQuestions = quiz.TotalVisibleQuestions
					existingQuiz.RevealAnswers = quiz.RevealAnswers
					existingQuiz.EnableRetry = quiz.EnableRetry
					existingQuiz.RetryAttempts = quiz.RetryAttempts
					existingQuiz.MinimumPassPercentage = quiz.MinimumPassPercentage

					// handle questions of each quiz
					for _, question := range quiz.Questions {
						if question.ID != nil && *question.ID != 0 {
							questionID := uint(*question.ID)
							incomingQuizQuestionIDs[questionID] = true

							if existingQuestion, found := quizQuestionMap[questionID]; found {
								// Update
								existingQuestion.Title = question.Title
								existingQuestion.Details = question.Details
								existingQuestion.Marks = question.Marks
								existingQuestion.AnswerRequired = question.AnswerRequired
								existingQuestion.AnswerExplanation = question.AnswerExplanation
								existingQuestion.Type = question.Type
								// existingQuestion.Media = question.Media

								if err := s.db.Save(&existingQuestion).Error; err != nil {
									return err
								}
							}
						} else {
							// Create new question
							newQuestion := models.QuizQuestion{
								QuizID:            quizID,
								Title:             question.Title,
								Details:           question.Details,
								Marks:             question.Marks,
								Type:              question.Type,
								Media:             nil,
								AnswerRequired:    question.AnswerRequired,
								AnswerExplanation: question.AnswerExplanation,
							}
							if err := s.db.Create(&newQuestion).Error; err != nil {
								return err
							}
							incomingQuizQuestionIDs[newQuestion.ID] = true
						}
					}

					if err := s.db.Save(&existingQuiz).Error; err != nil {
						return err
					}
				}
			} else {
				// Create new quiz
				newQuiz := models.CourseQuiz{
					ChapterID:             chapterID,
					Title:                 quiz.Title,
					Instructions:          quiz.Instructions,
					IsPublished:           quiz.IsPublished,
					TimeLimit:             quiz.TimeLimit,
					TimeLimitOption:       quiz.TimeLimitOption,
					RandomizeQuestions:    quiz.RandomizeQuestions,
					SingleQuizView:        quiz.SingleQuizView,
					TotalVisibleQuestions: quiz.TotalVisibleQuestions,
					RevealAnswers:         quiz.RevealAnswers,
					EnableRetry:           quiz.EnableRetry,
					RetryAttempts:         quiz.RetryAttempts,
					MinimumPassPercentage: quiz.MinimumPassPercentage,
				}
				if err := s.db.Create(&newQuiz).Error; err != nil {
					return err
				}

				// create questions
				for _, question := range quiz.Questions {
					newQuestion := models.QuizQuestion{
						QuizID:            newQuiz.ID,
						Title:             question.Title,
						Details:           question.Details,
						Type:              question.Type,
						Marks:             question.Marks,
						AnswerRequired:    question.AnswerRequired,
						AnswerExplanation: question.AnswerExplanation,
						// Media:     question.Media,
					}
					if err := s.db.Create(&newQuestion).Error; err != nil {
						return err
					}
				}

				incomingQuizIDs[newQuiz.ID] = true
			}
		}
	}

	// Delete removed lessons
	for id := range lessonMap {
		if !incomingLessonIDs[id] {
			_ = s.db.Where("id = ?", id).Delete(&models.CourseLesson{})
		}
	}

	// Delete removed assignments
	for id := range assignmentMap {
		if !incomingAssignmentIDs[id] {
			_ = s.db.Where("id = ?", id).Delete(&models.CourseAssignment{})
		}
	}

	// Delete removed quizes
	for id := range quizMap {
		if !incomingQuizIDs[id] {
			_ = s.db.Where("id = ?", id).Delete(&models.CourseQuiz{})
		}
	}

	// Delete removed questions
	for id := range quizQuestionMap {
		if !incomingQuizQuestionIDs[id] {
			_ = s.db.Where("id = ?", id).Delete(&models.QuizQuestion{})
		}
	}

	// Delete removed chapters
	for id := range chapterMap {
		if !incomingChapterIDs[id] {
			_ = s.db.Where("id = ?", id).Delete(&models.CourseChapter{})
		}
	}

	// Replace instructors
	// 1. Fetch existing instructor IDs for this course
	var existingInstructors []models.CourseInstructor
	if err := s.db.Where("course_id = ?", courseID).Find(&existingInstructors).Error; err != nil {
		return err
	}

	// Create a map for quick lookup of existing instructor IDs
	existingMap := make(map[uint]bool)
	for _, inst := range existingInstructors {
		existingMap[inst.InstructorID] = true
	}

	// Create a map for incoming instructors for quick lookup
	incomingMap := make(map[uint]bool)
	for _, instID := range input.Instructors {
		incomingMap[uint(instID)] = true
	}

	// 2. Delete instructors that exist in DB but not in input
	for _, inst := range existingInstructors {
		if !incomingMap[inst.InstructorID] {
			if err := s.db.Where("course_id = ? AND instructor_id = ?", courseID, inst.InstructorID).Delete(&models.CourseInstructor{}).Error; err != nil {
				return err
			}
		}
	}

	// 3. Add new instructors that are in input but not in DB
	for _, instID := range input.Instructors {
		uid := uint(instID)
		if !existingMap[uid] {
			newInst := models.CourseInstructor{
				CourseID:     courseID,
				InstructorID: uid,
			}
			if err := s.db.Create(&newInst).Error; err != nil {
				return err
			}
		}

	}

	// general settings
	var difficultyLevelPtr *models.DifficultyLevel
	if input.GeneralSettings.DifficultyLevel != "" {
		difficultyLevelPtr = &input.GeneralSettings.DifficultyLevel
	} else {
		defaultVal := models.All
		difficultyLevelPtr = &defaultVal
	}

	deafultLng := "english"

	updateGeneralSettings := models.CourseGeneralSettings{
		DifficultyLevel: difficultyLevelPtr,
		MaximumStudent:  utils.ZeroToNil(input.GeneralSettings.MaximumStudent),
		Language:        &deafultLng,
		CategoryID:      input.GeneralSettings.CategoryID,
		SubCategoryID:   utils.ZeroToNil(input.GeneralSettings.SubCategoryID),
		Duration:        utils.ZeroToNil(input.GeneralSettings.Duration),
	}

	// If exists, update
	return s.db.Where("course_id = ?", courseID).Select("difficulty_level", "maximum_student", "language", "category_id", "sub_category_id", "duration").Updates(&updateGeneralSettings).Error

}

func (s *courseService) Delete(id uint, tenantID uint) error {

	// 1. Check and delete CourseInstructors if any
	var existingInstructors []models.CourseInstructor
	if err := s.db.Where("course_id = ?", id).Find(&existingInstructors).Error; err != nil {
		return err
	}

	if len(existingInstructors) > 0 {
		if err := s.db.Where("course_id = ?", id).Delete(&models.CourseInstructor{}).Error; err != nil {
			return err
		}
	}

	// 2. Optionally delete GeneralSettings if it exists
	var generalSettings models.CourseGeneralSettings
	if err := s.db.Where("course_id = ?", id).First(&generalSettings).Error; err != nil {
		// Check if it's an actual error or just record not found
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	} else {
		// Found, so delete
		if err := s.db.Delete(&generalSettings).Error; err != nil {
			return err
		}
	}

	// 3. Check and delete the CourseDetails
	var existingCourseDetails models.CourseDetails
	if err := s.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&existingCourseDetails).Error; err != nil {
		return err
	}

	// Delete CDN image if present
	if existingCourseDetails.FeaturedImage != nil {
		if err := utils.DeleteFromBunny(*existingCourseDetails.FeaturedImage); err != nil {
			fmt.Printf("Failed to delete image from CDN: %v\n", err)
		}
	}

	// check and delete lesson resources
	var existingLessonsResources []models.LessonResource
	if err := s.db.Where("course_id = ?", id).Find(&existingLessonsResources).Error; err != nil {
		return err
	}

	if len(existingLessonsResources) > 0 {
		for _, res := range existingLessonsResources {
			// Delete from Bunny CDN
			if err := utils.DeleteFromBunny(res.FilePath); err != nil {
				fmt.Printf("Failed to delete resource from CDN: %v\n", err)
			}
		}
	}

	// Finally delete the course
	if err := s.db.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&models.CourseDetails{}).Error; err != nil {
		return err
	}

	return nil
}

func (s *courseService) DeleteLessonResource(courseID uint, resourceID uint) error {
	var existingLessonsResources models.LessonResource
	if err := s.db.Where("course_id = ? AND id = ?", courseID, resourceID).First(&existingLessonsResources).Error; err != nil {
		return err
	}

	// Delete from Bunny CDN
	if err := utils.DeleteFromBunny(existingLessonsResources.FilePath); err != nil {
		fmt.Printf("Failed to delete resource from CDN: %v\n", err)
	}

	return s.db.Where("course_id = ? AND id = ?", courseID, resourceID).Delete(&models.LessonResource{}).Error

}

func (s *courseService) ReorderCourses(tenantID uint, input ReorderRequest) error {
	var courses []models.CourseDetails

	// 1. Fetch courses ordered by index for the tenant
	if err := s.db.
		Where("tenant_id = ?", tenantID).
		Select("id, position").
		Order("position ASC").
		Find(&courses).Error; err != nil {
		return err
	}

	// 2. Find positions of active and over
	activeIndex := -1
	overIndex := -1

	for i, course := range courses {
		if course.ID == input.ActiveID {
			activeIndex = i
		}
		if course.ID == input.OverID {
			overIndex = i
		}
	}

	if activeIndex == -1 || overIndex == -1 {
		return fmt.Errorf("invalid course IDs")
	}

	// 3. Move the dragged item in the slice
	moved := courses[activeIndex]
	courses = append(courses[:activeIndex], courses[activeIndex+1:]...) // remove old position

	// insert at new position
	newList := make([]models.CourseDetails, 0, len(courses)+1)
	newList = append(newList, courses[:overIndex]...)
	newList = append(newList, moved)
	newList = append(newList, courses[overIndex:]...)

	courses = newList

	// 4. Update new indexes
	for i, c := range courses {
		if err := s.db.Model(&models.CourseDetails{}).
			Where("id = ?", c.ID).
			Update("position", i).Error; err != nil {
			return err
		}
	}

	return nil
}
