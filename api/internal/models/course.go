package models

import (
	"dashlearn/internal/utils"
	"time"

	"gorm.io/datatypes"
)

type Visibility string
type CoursePricingModel string
type DifficultyLevel string
type Access string

const (
	Public    Visibility = "public"
	Private   Visibility = "private"
	Protected Visibility = "protected"
)

const (
	CoursePricingModelFree CoursePricingModel = "free"
	CoursePricingModelPaid CoursePricingModel = "paid"
)

const (
	All          DifficultyLevel = "all"
	Beginner     DifficultyLevel = "beginner"
	Intermediate DifficultyLevel = "intermediate"
	Expert       DifficultyLevel = "expert"
)

const (
	Draft     Access = "draft"
	Published Access = "published"
)

type IntroVideo struct {
	Type   string `json:"type"`
	Source string `json:"source"`
}

type CourseDetails struct {
	ID              uint                     `gorm:"primaryKey;autoIncrement" json:"id"`
	Title           string                   `json:"title"`
	Slug            string                   `json:"slug"`
	Summary         string                   `gorm:"type:text" json:"summary"`
	Description     *string                  `gorm:"type:text" json:"description"`
	Visibility      Visibility               `gorm:"type:enum('public','private','protected');default:'public'" json:"visibility"`
	IsScheduled     *bool                    `gorm:"column:is_scheduled;default:false" json:"is_scheduled"`
	ScheduleDate    *string                  `gorm:"type:date" json:"schedule_date"`
	ScheduleTime    *string                  `gorm:"type:time" json:"schedule_time"`
	ShowCommingSoon *bool                    `gorm:"default:false" json:"show_comming_soon"`
	FeaturedImage   *string                  `gorm:"column:featured_image" json:"featured_image"`
	IntroVideo      *utils.JSONB[IntroVideo] `gorm:"type:json;column:intro_video" json:"intro_video"`
	PricingModel    CoursePricingModel       `gorm:"column:pricing_model;enum('free','paid');default:'free'" json:"pricing_model"`
	RegularPrice    *float32                 `gorm:"column:regular_price;default:0" json:"regular_price"`
	SalePrice       *float32                 `gorm:"column:sale_price;default:0" json:"sale_price"`
	Tags            datatypes.JSON           `gorm:"type:json" json:"tags"`
	Overview        datatypes.JSON           `gorm:"type:json" json:"overview"`
	AuthorID        uint                     `gorm:"column:author_id" json:"author_id"`
	Author          User                     `gorm:"foreignKey:AuthorID;references:ID" json:"author"`
	Position        int64                    `gorm:"default:0" json:"position"`
	CreatedAt       time.Time                `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time                `gorm:"autoUpdateTime" json:"updated_at"`
	TenantID        uint                     `gorm:"column:tenant_id" json:"-"`
	Tenant          Tenant                   `gorm:"foreignKey:TenantID;references:ID" json:"-"`
	Chapters        []CourseChapter          `gorm:"foreignKey:CourseID;references:ID" json:"course_chapters"`
	GeneralSettings CourseGeneralSettings    `gorm:"foreignKey:CourseID;references:ID" json:"general_settings"`
	Instructors     []CourseInstructor       `gorm:"foreignKey:CourseID;references:ID" json:"course_instructors"`
	Enrollments     []Enrollment             `gorm:"foreignKey:CourseID;references:ID" json:"enrollments"`
}

type CourseChapter struct {
	ID          uint               `gorm:"primaryKey;autoIncrement" json:"id"`
	Position    int                `gorm:"default:0" json:"position"`
	Title       string             `json:"title"`
	Description *string            `gorm:"type:text" json:"description"`
	Access      Access             `gorm:"enum('draft','published');default:'published'" json:"access"`
	CreatedAt   time.Time          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time          `gorm:"autoUpdateTime" json:"updated_at"`
	CourseID    uint               `gorm:"column:course_id" json:"course_id"`
	Lessons     []CourseLesson     `gorm:"foreignKey:ChapterID;references:ID" json:"course_lessons"`
	Assignments []CourseAssignment `gorm:"foreignKey:ChapterID;references:ID" json:"assignments"`
	Quizzes     []CourseQuiz       `gorm:"foreignKey:ChapterID;references:ID" json:"quizzes"`
	// CourseDetails CourseDetails `gorm:"foreignKey:CourseID;references:ID" json:"course_details"`
}

type LessonType string

const (
	Video       LessonType = "video"
	LiveSession LessonType = "live_session"
	Audio       LessonType = "audio"
	Text        LessonType = "text"
)

type LessonSourceType string

const (
	Youtube    LessonSourceType = "youtube"
	Vimeo      LessonSourceType = "vimeo"
	CustomCode LessonSourceType = "custom_code"
	UploadFile LessonSourceType = "upload"
	SoundCloud LessonSourceType = "sound_cloud"
	Spotify    LessonSourceType = "spotify"
)

type Source struct {
	Data          string  `json:"data"`
	IsFile        bool    `json:"is_file"`
	PlaybackTimes *string `json:"playback_times"`
}

type CourseLesson struct {
	ID              uint                `gorm:"primaryKey;autoIncrement" json:"id"`
	Title           string              `json:"title"`
	Description     *string             `gorm:"type:text" json:"description"`
	LessonType      LessonType          `gorm:"enum('video','live_session','audio','text');default:'video'" json:"lesson_type"`
	SourceType      LessonSourceType    `gorm:"enum('youtube','vimeo', 'sound_cloud','spotify','custom_code','upload');default:'youtube'" json:"source_type"`
	Source          utils.JSONB[Source] `gorm:"type:json" json:"source"`
	IsPublished     bool                `gorm:"default:false" json:"is_published"`
	IsPublic        bool                `gorm:"default:false" json:"is_public"`
	IsScheduled     *bool               `gorm:"column:is_scheduled;default:false" json:"is_scheduled"`
	ScheduleDate    *string             `gorm:"type:date" json:"schedule_date"`
	ScheduleTime    *string             `gorm:"type:time" json:"schedule_time"`
	ShowCommingSoon *bool               `gorm:"default:false" json:"show_comming_soon"`
	Position        int                 `gorm:"default:0" json:"position"`
	CreatedAt       time.Time           `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time           `gorm:"autoUpdateTime" json:"updated_at"`
	ChapterID       uint                `gorm:"column:chapter_id" json:"chapter_id"`
	Resources       []LessonResource    `gorm:"foreignKey:LessonID;references:ID" json:"resources"`
	// CourseChapter CourseChapter      `gorm:"foreignKey:ChapterID;references:ID" json:"course_chapter"`

}

type LessonResource struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CourseID  uint      `gorm:"column:course_id" json:"course_id"`
	LessonID  uint      `gorm:"column:lesson_id" json:"lesson_id"`
	MimeType  string    `gorm:"column:mime_type" json:"mine_type"`
	Title     string    `json:"title"`
	FilePath  string    `gorm:"column:file_path" json:"file_path"`
	Position  int       `json:"position"`
	FileSize  int64     `gorm:"column:file_size" json:"file_size"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type CourseGeneralSettings struct {
	ID              uint             `gorm:"primaryKey;autoIncrement" json:"id"`
	CourseID        uint             `gorm:"column:course_id" json:"course_id"`
	DifficultyLevel *DifficultyLevel `gorm:"column:difficulty_level;enum('all','beginner','intermediate','expert');default:'all'" json:"difficulty_level"`
	MaximumStudent  *int32           `gorm:"column:maximum_student;default:0" json:"maximum_student"`
	Language        *string          `gorm:"default:'english'" json:"language"`
	CategoryID      uint             `gorm:"column:category_id" json:"category_id"`
	Category        Category         `gorm:"foreignKey:CategoryID;references:ID" json:"category"`
	SubCategoryID   *uint            `gorm:"column:sub_category_id" json:"sub_category_id"`
	SubCategory     SubCategory      `gorm:"foreignKey:SubCategoryID;references:ID" json:"sub_category"`
	Duration        *string          `json:"duration"`
	CreatedAt       time.Time        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time        `gorm:"autoUpdateTime" json:"updated_at"`
	// Course            CourseDetails   `gorm:"foreignKey:CourseID;references:ID" json:"course"`
}

type CourseInstructor struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	CourseID     uint       `gorm:"column:course_id" json:"course_id"`
	InstructorID uint       `gorm:"column:instructor_id" json:"instructor_id"`
	Instructor   Instructor `gorm:"foreignKey:InstructorID;references:ID" json:"instructor"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}
