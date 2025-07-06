package course

import (
	"context"
	"dashlearn/models"
	"dashlearn/utils"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type CourseService interface {
	GetAll(tenantID uint) ([]CourseDetailsResponse, error)
	GetAllLite(tenantID uint) ([]struct {
		ID    uint   `json:"id"`
		Title string `json:"title"`
	}, error)
	GetAllPublic(tenantID uint) ([]CourseDetailsPublicResponse, error)
	GetByID(tenantID uint, courseID uint) (CourseDetailsResponse, error)
	Create(input CourseDetailsInput, tenantID uint, userID uint) error
	Update(courseID, tenantID, userID uint, input CourseDetailsInput) error
	Delete(id uint, tenantID uint) error
}

type courseService struct {
	db *gorm.DB
}

func NewCourseService(db *gorm.DB) CourseService {
	return &courseService{
		db: db,
	}
}

func (s *courseService) GetAll(tenantID uint) ([]CourseDetailsResponse, error) {
	var courses []CourseDetailsResponse

	err := s.db.Where("tenant_id = ?", tenantID).Preload("Author").Preload("Chapters").Preload("Chapters.Lessons").Preload("GeneralSettings").Preload("GeneralSettings.Category").Preload("Instructors").Preload("Instructors.Instructor").Find(&courses).Error

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

	err := s.db.Table("course_details").Where("tenant_id = ?", tenantID).Select("id", "title", "tenant_id").Find(&courses).Error

	return courses, err
}

func (s *courseService) GetAllPublic(tenantID uint) ([]CourseDetailsPublicResponse, error) {
	var courses []CourseDetailsPublicResponse

	err := s.db.Where("tenant_id = ?", tenantID).Preload("GeneralSettings").Preload("GeneralSettings.Category").Find(&courses).Error

	return courses, err
}

func (s *courseService) GetByID(tenantID uint, courseID uint) (CourseDetailsResponse, error) {
	var course CourseDetailsResponse

	err := s.db.
		Where("tenant_id = ? AND id = ?", tenantID, courseID).
		Preload("Author").
		Preload("Chapters").
		Preload("Chapters.Lessons").
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

func (s *courseService) Create(input CourseDetailsInput, tenantID uint, userID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var videoPtr *models.IntroVideo

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

		newCourseDetails := models.CourseDetails{
			Title:           input.Title,
			Summary:         input.Summary,
			Description:     utils.ZeroToNil(input.Description),
			Visibility:      input.Visibility,
			IsScheduled:     utils.ZeroToNil(input.IsScheduled),
			ScheduleDate:    utils.ZeroToNil(input.ScheduleDate),
			ScheduleTime:    utils.ZeroToNil(input.ScheduleTime),
			PricingModel:    input.PricingModel,
			RegularPrice:    input.RegularPrice,
			SalePrice:       input.SalePrice,
			ShowCommingSoom: input.ShowCommingSoom,
			Tags:            tagsJSON,
			Overview:        overviewJSON,
			FeaturedImage:   input.FeaturedImage,
			IntroVideo:      introVideo,
			AuthorID:        userID,
			TenantID:        tenantID,
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

				newCourseLesson := models.CourseLesson{
					ChapterID:   newCourseChapter.ID,
					Title:       lesson.Title,
					Description: utils.EmptyStringToNil(lesson.Description),
					Position:    idx,
					LessonType:  lesson.LessonType,
					SourceType:  lesson.SourceType,
					Source:      sourceJSON,
					IsPublished: lesson.IsPublished,
					IsPublic:    lesson.IsPublic,
				}

				if err := tx.Create(&newCourseLesson).Error; err != nil {
					return err
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

	if input.FeaturedImage != nil && *input.FeaturedImage != "" && existing.FeaturedImage != nil {
		if delErr := utils.DeleteCDNFile(context.Background(), *existing.FeaturedImage); delErr != nil {
			// You can log or ignore deletion errors as per your need
			fmt.Println("Failed to delete old file:", delErr)
		}
	}

	// Normalize & assign values
	tagsJSON, _ := utils.NormalizeTags(input.Tags)
	overviewJSON, _ := utils.NormalizeTags(input.Overview)

	var introVideo *utils.JSONB[models.IntroVideo]
	if input.IntroVideo != nil {
		introVideo = &utils.JSONB[models.IntroVideo]{Data: *input.IntroVideo}
	}

	// Update course
	updateData := models.CourseDetails{
		Title:           input.Title,
		Summary:         input.Summary,
		Description:     utils.ZeroToNil(input.Description),
		Visibility:      input.Visibility,
		IsScheduled:     utils.ZeroToNil(input.IsScheduled),
		ScheduleDate:    utils.ZeroToNil(input.ScheduleDate),
		ScheduleTime:    utils.ZeroToNil(input.ScheduleTime),
		PricingModel:    input.PricingModel,
		RegularPrice:    input.RegularPrice,
		SalePrice:       input.SalePrice,
		ShowCommingSoom: input.ShowCommingSoom,
		Tags:            tagsJSON,
		Overview:        overviewJSON,
		FeaturedImage:   input.FeaturedImage,
		IntroVideo:      introVideo,
		AuthorID:        userID,
		TenantID:        tenantID,
	}

	if err := s.db.Where("id = ?", courseID).Updates(&updateData).Error; err != nil {
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

	existingAssignmentMap := make(map[uint]models.CourseAssignment)
	s.db.Preload("Assignments").Where("course_id = ?", courseID).Find(&existingChaptersForUpdate)
	for _, chapter := range existingChaptersForUpdate {
		for _, assignment := range chapter.Assignments {
			existingAssignmentMap[assignment.ID] = assignment
		}
	}

	// Maps to track incoming IDs
	incomingChapterIDs := make(map[uint]bool)
	incomingLessonIDs := make(map[uint]bool)
	incomingAssignmentIDs := make(map[uint]bool)

	// Fetch all existing chapters and their lessons
	var existingChapters []models.CourseChapter
	s.db.Preload("Lessons").Preload("Assignments").Where("course_id = ?", courseID).Find(&existingChapters)

	chapterMap := make(map[uint]models.CourseChapter)
	lessonMap := make(map[uint]models.CourseLesson)
	assignmentMap := make(map[uint]models.CourseAssignment)

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
	}

	for chIdx, chapter := range input.CourseChapters {
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
			sourceJSON := utils.JSONB[models.Source]{Data: lesson.Source}

			if lesson.ID != nil && *lesson.ID != 0 {
				lessonID := uint(*lesson.ID)
				incomingLessonIDs[lessonID] = true

				if existingLesson, found := lessonMap[lessonID]; found {
					// Update
					existingLesson.Title = lesson.Title
					existingLesson.Description = utils.EmptyStringToNil(lesson.Description)
					existingLesson.Position = lIdx
					existingLesson.LessonType = lesson.LessonType
					existingLesson.SourceType = lesson.SourceType
					existingLesson.Source = sourceJSON
					existingLesson.IsPublished = lesson.IsPublished
					existingLesson.IsPublic = lesson.IsPublic
					existingLesson.ChapterID = chapterID

					if err := s.db.Save(&existingLesson).Error; err != nil {
						return err
					}
				}
			} else {
				// Create new lesson
				newLesson := models.CourseLesson{
					ChapterID:   chapterID,
					Title:       lesson.Title,
					Description: utils.EmptyStringToNil(lesson.Description),
					Position:    lIdx,
					LessonType:  lesson.LessonType,
					SourceType:  lesson.SourceType,
					Source:      sourceJSON,
					IsPublished: lesson.IsPublished,
					IsPublic:    lesson.IsPublic,
				}
				if err := s.db.Create(&newLesson).Error; err != nil {
					return err
				}
				incomingLessonIDs[newLesson.ID] = true
			}
		}

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
				incomingLessonIDs[newAssignment.ID] = true
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
		Duration:        utils.ZeroToNil(input.GeneralSettings.Duration),
	}

	// If exists, update
	return s.db.Where("course_id = ?", courseID).Updates(&updateGeneralSettings).Error

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
		if err := utils.DeleteCDNFile(context.Background(), *existingCourseDetails.FeaturedImage); err != nil {
			fmt.Printf("Failed to delete image from CDN: %v\n", err)
		}
	}

	// Finally delete the course
	if err := s.db.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&models.CourseDetails{}).Error; err != nil {
		return err
	}

	return nil
}
