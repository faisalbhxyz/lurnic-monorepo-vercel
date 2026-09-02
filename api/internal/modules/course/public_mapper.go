package course

import (
	"dashlearn/internal/models"
	quizmodule "dashlearn/internal/modules/quiz"
	"dashlearn/internal/response"
	"dashlearn/internal/utils"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func preloadPublicCourse(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Author").
		Preload("Chapters", "access = 'published'").
		Preload("Chapters.Lessons", "is_published = true").
		Preload("Chapters.Lessons.Resources").
		Preload("Chapters.Assignments", "is_published = true").
		Preload("Chapters.Quizzes", "is_published = true").
		Preload("Chapters.Quizzes.Questions").
		Preload("GeneralSettings").
		Preload("GeneralSettings.Category").
		Preload("CertificateSettings").
		Preload("Instructors").
		Preload("Instructors.Instructor").
		Preload("Enrollments")
}

// BuildCourseDetailsPublicResponse maps a preloaded CourseDetails model to the storefront response.
func BuildCourseDetailsPublicResponse(modelCourse *models.CourseDetails) *response.CourseDetailsPublicResponse {
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
				IsPublished: lesson.IsPublished,
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
				questions = append(questions, quizmodule.SanitizeQuestionResponse(question, false))
			}

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

	enrollments := make([]response.EnrolledCourseRes, len(modelCourse.Enrollments))
	for i, enrollment := range modelCourse.Enrollments {
		enrollments[i] = response.EnrolledCourseRes{
			ID:        enrollment.ID,
			CourseID:  enrollment.CourseID,
			StudentID: enrollment.StudentID,
			CreatedAt: enrollment.CreatedAt,
			UpdatedAt: enrollment.UpdatedAt,
		}
	}

	isScheduled := false
	if modelCourse.IsScheduled != nil {
		isScheduled = *modelCourse.IsScheduled
	}

	res := &response.CourseDetailsPublicResponse{
		ID:            modelCourse.ID,
		Title:         modelCourse.Title,
		Slug:          modelCourse.Slug,
		Summary:       modelCourse.Summary,
		Description:   modelCourse.Description,
		Visibility:    modelCourse.Visibility,
		IsScheduled:   isScheduled,
		FeaturedImage: modelCourse.FeaturedImage,
		IntroVideo:    modelCourse.IntroVideo,
		PricingModel:  modelCourse.PricingModel,
		RegularPrice:  modelCourse.RegularPrice,
		SalePrice:     modelCourse.SalePrice,
		Tags:          modelCourse.Tags,
		Overview:      modelCourse.Overview,
		Chapters:      chapters,
		Instructors:   instructors,
		Enrollments:   enrollments,
	}

	if modelCourse.GeneralSettings.ID != 0 {
		res.GeneralSettings = &response.CourseGeneralSettingsResponse{
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
		}
	}

	return res
}

// LoadPublicCoursesByIDs loads full storefront course payloads keyed by course ID.
func LoadPublicCoursesByIDs(db *gorm.DB, tenantID uint, courseIDs []uint) (map[uint]*response.CourseDetailsPublicResponse, error) {
	result := make(map[uint]*response.CourseDetailsPublicResponse, len(courseIDs))
	if len(courseIDs) == 0 {
		return result, nil
	}

	var courses []models.CourseDetails
	if err := preloadPublicCourse(db).
		Where("tenant_id = ? AND id IN ?", tenantID, courseIDs).
		Find(&courses).Error; err != nil {
		return nil, fmt.Errorf("failed to load courses: %w", err)
	}

	for i := range courses {
		course := &courses[i]
		result[course.ID] = BuildCourseDetailsPublicResponse(course)
	}
	return result, nil
}

// LoadPublicCourseBySlug loads a single course by slug for the storefront.
func LoadPublicCourseBySlug(db *gorm.DB, tenantID uint, slug string) (*response.CourseDetailsPublicResponse, error) {
	var modelCourse models.CourseDetails
	err := preloadPublicCourse(db).
		Where("tenant_id = ? AND slug = ?", tenantID, slug).
		First(&modelCourse).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to retrieve course: %w", err)
	}
	return BuildCourseDetailsPublicResponse(&modelCourse), nil
}
