package course

import (
	"dashlearn/models"
	"time"

	"gorm.io/datatypes"
)

type CourseDetailsResponse struct {
	ID              uint                         `gorm:"primaryKey;autoIncrement" json:"id"`
	Title           string                       `json:"title"`
	Summary         string                       `gorm:"type:text" json:"summary"`
	Description     *string                      `gorm:"type:text" json:"description"`
	Visibility      models.Visibility            `gorm:"type:enum('public','private','protected');default:'public'" json:"visibility"`
	IsScheduled     *bool                        `gorm:"default:false" json:"is_scheduled"`
	ScheduleDate    *time.Time                   `gorm:"type:date" json:"schedule_date"`
	ScheduleTime    *time.Time                   `gorm:"type:time" json:"schedule_time"`
	FeaturedImage   *string                      `gorm:"column:featured_image" json:"featured_image"`
	IntroVideo      datatypes.JSON               `gorm:"type:json;column:intro_video" json:"intro_video"`
	PricingModel    models.CoursePricingModel    `gorm:"column:pricing_model;enum('free','paid');default:'free'" json:"pricing_model"`
	RegularPrice    *float32                     `gorm:"column:regular_price;default:0" json:"regular_price"`
	SalePrice       *float32                     `gorm:"column:sale_price;default:0" json:"sale_price"`
	ShowCommingSoon *bool                        `gorm:"column:show_comming_soom;default:false" json:"show_comming_soon"`
	Tags            datatypes.JSON               `gorm:"type:json" json:"tags"`
	Overview        datatypes.JSON               `gorm:"type:json" json:"overview"`
	AuthorID        uint                         `gorm:"column:author_id" json:"author_id"`
	Author          models.UserWithoutRole       `gorm:"foreignKey:AuthorID;references:ID" json:"author"`
	CreatedAt       time.Time                    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time                    `gorm:"autoUpdateTime" json:"updated_at"`
	Chapters        []CourseChapterResponse      `gorm:"foreignKey:CourseID;references:ID" json:"course_chapters"`
	GeneralSettings models.CourseGeneralSettings `gorm:"foreignKey:CourseID;references:ID" json:"general_settings"`
	Instructors     []models.CourseInstructor    `gorm:"foreignKey:CourseID;references:ID" json:"course_instructors"`
	Enrollments     []models.EnrollmentResponse  `gorm:"foreignKey:CourseID;references:ID" json:"enrollments"`
}

func (CourseDetailsResponse) TableName() string {
	return "course_details"
}

type CourseChapterResponse struct {
	ID          uint                       `gorm:"primaryKey;autoIncrement" json:"id"`
	Position    int                        `gorm:"default:0" json:"position"`
	Title       string                     `json:"title"`
	Description *string                    `gorm:"type:text" json:"description"`
	Access      models.Access              `gorm:"enum('draft','published');default:'published'" json:"access"`
	CreatedAt   time.Time                  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time                  `gorm:"autoUpdateTime" json:"updated_at"`
	CourseID    uint                       `gorm:"column:course_id" json:"course_id"`
	Lessons     []CourseLessonResponse     `gorm:"foreignKey:ChapterID;references:ID" json:"course_lessons"`
	Assignments []CourseAssignmentResponse `gorm:"foreignKey:ChapterID;references:ID" json:"assignments"`
	Quizzes     []CourseQuizResponse       `gorm:"foreignKey:ChapterID;references:ID" json:"quizzes"`
}

func (CourseChapterResponse) TableName() string {
	return "course_chapters"
}

type CourseLessonResponse struct {
	ID          uint                    `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string                  `json:"title"`
	Description *string                 `gorm:"type:text" json:"description"`
	LessonType  models.LessonType       `gorm:"enum('video','live_session','audio','text');default:'video'" json:"lesson_type"`
	SourceType  models.LessonSourceType `gorm:"enum('youtube','vimeo', 'sound_cloud','spotify','custom_code','upload');default:'youtube'" json:"source_type"`
	Source      datatypes.JSON          `gorm:"type:json" json:"source"`
	IsPublished bool                    `gorm:"default:false" json:"is_published"`
	IsPublic    bool                    `gorm:"default:false" json:"is_public"`
	Resources   datatypes.JSON          `gorm:"type:json" json:"resources"` // filename, mimetype, url, size
	Position    int                     `gorm:"default:0" json:"position"`
	CreatedAt   time.Time               `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time               `gorm:"autoUpdateTime" json:"updated_at"`
	ChapterID   uint                    `gorm:"column:chapter_id" json:"chapter_id"`
}

func (CourseLessonResponse) TableName() string {
	return "course_lessons"
}

type CourseAssignmentResponse struct {
	ID               uint                             `gorm:"primaryKey;autoIncrement" json:"id"`
	CourseID         uint                             `json:"course_id" gorm:"column:course_id"`
	ChapterID        uint                             `json:"chapter_id" gorm:"column:chapter_id"`
	Title            string                           `json:"title" gorm:"type:varchar(255)"`
	Instructions     string                           `json:"instructions" gorm:"type:text"`
	Attachments      *datatypes.JSON                  `gorm:"type:json" json:"attachments"`
	IsPublished      bool                             `json:"is_published" gorm:"column:is_published;default:false"`
	TimeLimit        int                              `json:"time_limit" gorm:"column:time_limit;default:1"`
	TimeLimitOption  models.CourseQuizTimeLimitOption `json:"time_limit_option" gorm:"column:time_limit_option;type:enum('minutes','hours','days','weeks','months');default:'weeks'"`
	FileUploadLimit  int                              `json:"file_upload_limit" gorm:"column:file_upload_limit;default:1"`
	TotalMarks       float32                          `json:"total_marks" gorm:"column:total_marks;default:1"`
	MinimumPassMarks float32                          `json:"minimum_pass_marks" gorm:"column:minimum_pass_marks;default:0"`
	CreatedAt        time.Time                        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time                        `gorm:"autoUpdateTime" json:"updated_at"`
}

func (CourseAssignmentResponse) TableName() string {
	return "course_assignments"
}

type CourseQuizResponse struct {
	ID                    uint                             `gorm:"primaryKey;autoIncrement" json:"id"`
	CourseID              uint                             `json:"course_id" gorm:"column:course_id"`
	ChapterID             uint                             `json:"chapter_id" gorm:"column:chapter_id"`
	Title                 string                           `json:"title" gorm:"type:varchar(255)"`
	Instructions          string                           `json:"instructions" gorm:"type:text"`
	IsPublished           bool                             `json:"is_published" gorm:"column:is_published;default:false"`
	RandomizeQuestions    bool                             `json:"randomize_questions" gorm:"column:randomize_questions;default:false"`
	SingleQuizView        bool                             `json:"single_quiz_view" gorm:"column:single_quiz_view;default:false"`
	TimeLimit             int                              `json:"time_limit" gorm:"column:time_limit;default:1"`
	TimeLimitOption       models.CourseQuizTimeLimitOption `json:"time_limit_option" gorm:"column:time_limit_option;type:enum('minutes','hours','days','weeks','months');default:'weeks'"`
	TotalVisibleQuestions *int                             `json:"total_visible_questions" gorm:"column:total_visible_questions;default:0"`
	RevealAnswers         bool                             `json:"reveal_answers" gorm:"column:reveal_answers;default:false"`
	EnableRetry           bool                             `json:"enable_retry" gorm:"column:enable_retry;default:false"`
	RetryAttempts         int                              `json:"retry_attempts" gorm:"column:retry_attempts;default:1"`
	MinimumPassPercentage float32                          `json:"minimum_pass_percentage" gorm:"column:minimum_pass_percentage;default:0"`
	CreatedAt             time.Time                        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt             time.Time                        `gorm:"autoUpdateTime" json:"updated_at"`
	Questions             []CourseQuizQuestionsResponse    `gorm:"foreignKey:QuizID;references:ID" json:"questions"`
}

func (CourseQuizResponse) TableName() string {
	return "course_quizzes"
}

type CourseQuizQuestionsResponse struct {
	ID                uint                    `gorm:"primaryKey;autoIncrement" json:"id"`
	QuizID            uint                    `gorm:"column:quiz_id" json:"quiz_id"`
	Quiz              models.CourseQuiz       `gorm:"foreignKey:QuizID;references:ID" json:"quiz"`
	Title             string                  `json:"title" gorm:"type:varchar(255)"`
	Details           *string                 `json:"details" gorm:"type:text"`
	Media             *datatypes.JSON         `gorm:"type:json" json:"media"`
	Type              models.QuizQuestionType `gorm:"type:enum('multiple_choice','single_choice','true_false')" json:"type"`
	Marks             float32                 `json:"marks" gorm:"column:marks;default:1"`
	AnswerRequired    bool                    `json:"answer_required" gorm:"column:answer_required;default:false"`
	AnswerExplanation *string                 `json:"answer_explanation" gorm:"type:text"`
	CreatedAt         time.Time               `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time               `gorm:"autoUpdateTime" json:"updated_at"`
}

func (CourseQuizQuestionsResponse) TableName() string {
	return "quiz_questions"
}

// ! COURSES RESPONSE without chapters, instructors'
type CourseDetailsPublicResponse struct {
	ID              uint                         `gorm:"primaryKey;autoIncrement" json:"id"`
	Title           string                       `json:"title"`
	Summary         string                       `gorm:"type:text" json:"summary"`
	Visibility      models.Visibility            `gorm:"type:enum('public','private','protected');default:'public'" json:"visibility"`
	IsScheduled     *bool                        `gorm:"default:false" json:"is_scheduled"`
	ScheduleDate    *time.Time                   `gorm:"type:date" json:"schedule_date"`
	ScheduleTime    *time.Time                   `gorm:"type:time" json:"schedule_time"`
	FeaturedImage   *string                      `gorm:"column:featured_image" json:"featured_image"`
	IntroVideo      datatypes.JSON               `gorm:"type:json;column:intro_video" json:"intro_video"`
	PricingModel    models.CoursePricingModel    `gorm:"column:pricing_model;enum('free','paid');default:'free'" json:"pricing_model"`
	RegularPrice    *float32                     `gorm:"column:regular_price;default:0" json:"regular_price"`
	SalePrice       *float32                     `gorm:"column:sale_price;default:0" json:"sale_price"`
	ShowCommingSoom *bool                        `gorm:"default:false" json:"show_comming_soom"`
	Tags            datatypes.JSON               `gorm:"type:json" json:"tags"`
	GeneralSettings models.CourseGeneralSettings `gorm:"foreignKey:CourseID;references:ID" json:"general_settings"`
}

func (CourseDetailsPublicResponse) TableName() string {
	return "course_details"
}
