package course

import (
	"dashlearn/internal/models"

	"gorm.io/datatypes"
)

type CourseDetailsInput struct {
	Title           string                    `form:"title" json:"title" binding:"required"`
	Summary         string                    `form:"summary" json:"summary" binding:"required"`
	Description     *string                   `form:"description" json:"description" binding:"omitempty"`
	Visibility      models.Visibility         `form:"visibility" json:"visibility" binding:"required,oneof=public private protected"`
	IsScheduled     string                    `form:"is_scheduled" json:"is_scheduled" binding:"omitempty"`
	ScheduleDate    *string                   `form:"schedule_date" json:"schedule_date" binding:"omitempty"`
	ScheduleTime    *string                   `form:"schedule_time" json:"schedule_time" binding:"omitempty"`
	ShowCommingSoon *string                   `form:"show_comming_soon" json:"show_comming_soon" binding:"omitempty"`
	PricingModel    models.CoursePricingModel `form:"pricing_model" json:"pricing_model" binding:"omitempty"`
	RegularPrice    *float32                  `form:"regular_price" json:"regular_price" binding:"omitempty"`
	SalePrice       *float32                  `form:"sale_price" json:"sale_price" binding:"omitempty"`
	Tags            *[]string                 `form:"tags" json:"tags" binding:"omitempty"`
	AuthorID        uint                      `form:"author_id" json:"author_id" binding:"required"`
	FeaturedImage   *string                   `form:"featured_image" json:"featured_image" binding:"omitempty"`
	IntroVideo      *models.IntroVideo        `form:"intro_video" json:"intro_video" binding:"omitempty"`
	Overview        *[]string                 `form:"overview" json:"overview" binding:"omitempty"`
	CourseChapters  []CreateCourseChapter     `form:"course_chapters" json:"course_chapters"`
	GeneralSettings CreateGeneralSettings     `form:"general_settings" json:"general_settings" binding:"required"`
	Instructors     []int32                   `form:"course_instructors" json:"course_instructors"`
}

type CreateCourseDetailsInput struct {
	Title           string                    `form:"title" json:"title" binding:"required"`
	Summary         string                    `form:"summary" json:"summary" binding:"required"`
	Description     *string                   `form:"description" json:"description" binding:"omitempty"`
	Visibility      models.Visibility         `form:"visibility" json:"visibility" binding:"required,oneof=public private protected"`
	IsScheduled     string                    `form:"is_scheduled" json:"is_scheduled" binding:"omitempty"`
	ScheduleDate    *string                   `form:"schedule_date" json:"schedule_date" binding:"omitempty"`
	ScheduleTime    *string                   `form:"schedule_time" json:"schedule_time" binding:"omitempty"`
	PricingModel    models.CoursePricingModel `form:"pricing_model" json:"pricing_model" binding:"omitempty"`
	RegularPrice    *float32                  `form:"regular_price" json:"regular_price" binding:"omitempty"`
	SalePrice       *float32                  `form:"sale_price" json:"sale_price" binding:"omitempty"`
	ShowCommingSoon *string                   `form:"show_comming_soon" json:"show_comming_soon" binding:"omitempty"`
	Tags            *[]string                 `form:"tags" json:"tags" binding:"omitempty"`
	AuthorID        uint                      `form:"author_id" json:"author_id" binding:"required"`
	IntroVideo      *models.IntroVideo        `form:"intro_video" json:"intro_video" binding:"omitempty"`
	Overview        *[]string                 `form:"overview" json:"overview" binding:"omitempty"`
	// FeaturedImage   *string             `form:"featured_image" binding:"omitempty"`
}

type CreateCourseChapter struct {
	ID            *int64                  `json:"id" form:"id" binding:"omitempty"`
	Position      int32                   `json:"position" form:"position" binding:"required"`
	Title         string                  `json:"title" form:"title" binding:"required"`
	Description   *string                 `json:"description" form:"description" binding:"omitempty"`
	Access        models.Access           `json:"access" form:"access" binding:"required,oneof=draft published"`
	CourseLessons []CreateCourseLesson    `form:"course_lessons" json:"course_lessons"`
	Quizzes       []CreateCourseQuizInput `form:"quizzes" json:"quizzes"`
	Assignments   []CreateAssignmentInput `form:"assignments" json:"assignments"`
}

type LessonResourceInput struct {
	URL      string `json:"url"`
	FileName string `json:"file_name"`
	Position int    `json:"position"`
	MimeType string `json:"mime_type"`
	Size     int64  `json:"size"`
}

type CreateCourseLesson struct {
	ID              *int64                  `json:"id" form:"id" binding:"omitempty"`
	Title           string                  `json:"title" form:"title" binding:"required"`
	Description     *string                 `json:"description" form:"description" binding:"omitempty"`
	LessonType      models.LessonType       `json:"lesson_type" form:"lesson_type" binding:"required,oneof=video live_session audio text"`
	SourceType      models.LessonSourceType `json:"source_type" form:"source_type" binding:"required,oneof=youtube vimeo sound_cloud spotify custom_code upload"`
	Source          models.Source           `json:"source" form:"source" binding:"omitempty"`
	IsPublished     bool                    `json:"is_published" form:"is_published" binding:"required"`
	IsPublic        bool                    `json:"is_public" form:"is_public" binding:"required"`
	Resources       []LessonResourceInput   `json:"resources"`
	IsScheduled     bool                    `form:"is_scheduled" json:"is_scheduled" binding:"omitempty"`
	ScheduleDate    *string                 `form:"schedule_date" json:"schedule_date" binding:"omitempty"`
	ScheduleTime    *string                 `form:"schedule_time" json:"schedule_time" binding:"omitempty"`
	ShowCommingSoon *bool                   `form:"show_comming_soon" json:"show_comming_soon" binding:"omitempty"`
}

type CreateGeneralSettings struct {
	DifficultyLevel models.DifficultyLevel `json:"difficulty_level" form:"difficulty_level" binding:"required,oneof=all beginner intermediate expert"`
	MaximumStudent  *int32                 `json:"maximum_student" form:"maximum_student" binding:"omitempty"`
	Language        *string                `json:"language" form:"language" binding:"omitempty"`
	CategoryID      uint                   `json:"category_id" form:"category_id" binding:"required"`
	SubCategoryID   *uint                  `json:"sub_category_id" form:"sub_category_id" binding:"omitempty"`
	Duration        *string                `json:"duration" form:"duration" binding:"omitempty"`
}

type CreateCourseQuizInput struct {
	ID                    *int64                           `json:"id" form:"id" binding:"omitempty"`
	Title                 string                           `json:"title" form:"title" binding:"required"`
	Instructions          string                           `json:"instructions" form:"instructions" binding:"required"`
	IsPublished           bool                             `json:"is_published" form:"is_published"`
	RandomizeQuestions    bool                             `json:"randomize_questions" form:"randomize_questions"`
	SingleQuizView        bool                             `json:"single_quiz_view" form:"single_quiz_view"`
	TimeLimit             int                              `json:"time_limit" form:"time_limit" binding:"required"`
	TimeLimitOption       models.CourseQuizTimeLimitOption `json:"time_limit_option" form:"time_limit_option"`
	TotalVisibleQuestions *int                             `json:"total_visible_questions" form:"total_visible_questions"`
	RevealAnswers         bool                             `json:"reveal_answers" form:"reveal_answers"`
	EnableRetry           bool                             `json:"enable_retry" form:"enable_retry"`
	RetryAttempts         int                              `json:"retry_attempts" form:"retry_attempts" binding:"required"`
	MinimumPassPercentage float32                          `json:"minimum_pass_percentage" form:"minimum_pass_percentage" binding:"required"`
	Questions             []CreateQuizQuestionInput        `json:"questions" form:"questions" binding:"required"`
}

type CreateQuizQuestionInput struct {
	ID                *int64                  `json:"id" form:"id" binding:"omitempty"`
	Title             string                  `json:"title" form:"title" binding:"required"`
	Details           *string                 `json:"details" form:"details"`
	Media             *datatypes.JSON         `json:"media" form:"media"`
	Type              models.QuizQuestionType `json:"type" form:"type"`
	Marks             float32                 `json:"marks" form:"marks" binding:"required"`
	AnswerRequired    bool                    `json:"answer_required" form:"answer_required"`
	AnswerExplanation *string                 `json:"answer_explanation" form:"answer_explanation"`
}

type CreateAssignmentInput struct {
	ID               *int64                           `json:"id" form:"id" binding:"omitempty"`
	Title            string                           `json:"title" form:"title" binding:"required"`
	Instructions     string                           `json:"instructions" form:"instructions"`
	Attachments      *datatypes.JSON                  `json:"attachments" form:"attachments"`
	IsPublished      bool                             `json:"is_published" form:"is_published"`
	TimeLimit        int                              `json:"time_limit" form:"time_limit" binding:"required"`
	TimeLimitOption  models.CourseQuizTimeLimitOption `json:"time_limit_option" form:"time_limit_option"`
	FileUploadLimit  int                              `json:"file_upload_limit" form:"file_upload_limit"`
	TotalMarks       float32                          `json:"total_marks" form:"total_marks" binding:"required"`
	MinimumPassMarks float32                          `json:"minimum_pass_marks" form:"minimum_pass_marks" binding:"required"`
}

type ReorderRequest struct {
	ActiveID uint `json:"activeID"`
	OverID   uint `json:"overID"`
}
