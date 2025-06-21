package course

import (
	"dashlearn/models"
	"dashlearn/utils"

	"gorm.io/gorm"
)

type CourseService interface {
	GetAll(tenantID uint) ([]models.CourseDetailsResponse, error)
	GetAllLite(tenantID uint) ([]struct {
		ID    uint   `json:"id"`
		Title string `json:"title"`
	}, error)
	GetAllPublic(tenantID uint) ([]models.CourseDetailsPublicResponse, error)
	GetByID(tenantID uint, courseID uint) (models.CourseDetailsResponse, error)
	Create(input CourseDetailsInput, tenantID uint, userID uint) error
	// Update (id uint, input CreateCourseInput, tenantID uint) error
	// Delete (id uint, tenantID uint) error
}

type courseService struct {
	db *gorm.DB
}

func NewCourseService(db *gorm.DB) CourseService {
	return &courseService{
		db: db,
	}
}

func (s *courseService) GetAll(tenantID uint) ([]models.CourseDetailsResponse, error) {
	var courses []models.CourseDetailsResponse

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

func (s *courseService) GetAllPublic(tenantID uint) ([]models.CourseDetailsPublicResponse, error) {
	var courses []models.CourseDetailsPublicResponse

	err := s.db.Where("tenant_id = ?", tenantID).Preload("GeneralSettings").Preload("GeneralSettings.Category").Find(&courses).Error

	return courses, err
}

func (s *courseService) GetByID(tenantID uint, courseID uint) (models.CourseDetailsResponse, error) {
	var course models.CourseDetailsResponse

	err := s.db.
		Where("tenant_id = ? AND id = ?", tenantID, courseID).
		Preload("Author").
		Preload("Chapters").
		Preload("Chapters.Lessons").
		Preload("GeneralSettings").
		Preload("GeneralSettings.Category").
		Preload("Instructors").
		Preload("Instructors.Instructor").
		Preload("Enrollments").
		First(&course).Error

	return course, err
}

func (s *courseService) Create(input CourseDetailsInput, tenantID uint, userID uint) error {
	// // Step 3: Debug log the final parsed object
	// if output, err := json.MarshalIndent(input, "", "  "); err == nil {
	// 	fmt.Println("Parsed Input:\n", string(output))
	// }

	//create course details

	var videoPtr *models.IntroVideo

	if input.IntroVideo == nil ||
		(input.IntroVideo.Type == "" && input.IntroVideo.Source == "") {
		// Set pointer to nil to store NULL in DB
		videoPtr = nil
	} else {
		// Valid data, assign normally
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
	} else {
		introVideo = nil
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

	if err := s.db.Create(&newCourseDetails).Error; err != nil {
		return err
	}

	// craete course chapter
	for idx, chapter := range input.CourseChapters {
		newCourseChapter := models.CourseChapter{
			CourseID:    newCourseDetails.ID,
			Title:       chapter.Title,
			Description: utils.EmptyStringToNil(chapter.Description),
			Position:    idx,
			Access:      chapter.Access,
		}

		if err := s.db.Create(&newCourseChapter).Error; err != nil {
			return err
		}

		if len(chapter.CourseLessons) > 0 {
			for idx, lesson := range chapter.CourseLessons {

				// sourceJSON, err := json.Marshal(lesson.Source)
				// if err != nil {
				// 	sourceJSON = nil
				// }

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
				if err := s.db.Create(&newCourseLesson).Error; err != nil {
					return err
				}

			}
		}
	}

	// course instructors
	for _, instructor := range input.Instructors {
		newCourseInstructor := models.CourseInstructor{
			CourseID:     newCourseDetails.ID,
			InstructorID: uint(instructor),
		}
		if err := s.db.Create(&newCourseInstructor).Error; err != nil {
			return err
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

	newGeneralSettings := models.CourseGeneralSettings{
		CourseID:        newCourseDetails.ID,
		DifficultyLevel: difficultyLevelPtr,
		MaximumStudent:  utils.ZeroToNil(input.GeneralSettings.MaximumStudent),
		Language:        &deafultLng,
		CategoryID:      input.GeneralSettings.CategoryID,
		Duration:        utils.ZeroToNil(input.GeneralSettings.Duration),
	}

	return s.db.Create(&newGeneralSettings).Error
}
