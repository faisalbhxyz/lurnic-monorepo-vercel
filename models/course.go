package models

import (
	"time"
)

type Visibility string
type PricingModel string
type DifficultyLevel string

const (
	Public    Visibility = "public"
	Private   Visibility = "private"
	Protected Visibility = "protected"
)

const (
	Free PricingModel = "free"
	Paid PricingModel = "paid"
)

const (
	All          DifficultyLevel = "all"
	Beginner     DifficultyLevel = "beginner"
	Intermediate DifficultyLevel = "intermediate"
	Expert       DifficultyLevel = "expert"
)

type Course struct {
	ID              uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	Title           string          `json:"title"`
	Description     *string         `gorm:"type:text" json:"description"`
	Visibility      Visibility      `gorm:"type:enum('public','private','protected');default:'public'" json:"visibility"`
	IsScheduled     bool            `gorm:"default:false" json:"is_scheduled"`
	ScheduleDate    *time.Time      `gorm:"type:date" json:"schedule_date"`
	ScheduleTime    *time.Time      `gorm:"type:time" json:"schedule_time"`
	ShowCommingSoon bool            `gorm:"column:show_comming_soon;default:false" json:"show_comming_soon"`
	FeaturedImage   *string         `gorm:"column:featured_image" json:"featured_image"`
	IntroVideo      *string         `gorm:"column:intro_video" json:"intro_video"`
	PricingModel    PricingModel    `gorm:"column:pricing_model;enum('free','paid');default:'free'" json:"pricing_model"`
	Tags            *[]string       `gorm:"type:json" json:"tags"`
	AuthorID        uint            `gorm:"column:author_id" json:"author_id"`
	Author          User            `gorm:"foreignKey:AuthorID;references:ID" json:"author"`
	DifficultyLevel DifficultyLevel `gorm:"column:difficulty_level;enum('all','beginner','intermediate','expert');default:'all'" json:"difficulty_level"`
	IsPublicCourse  bool            `gorm:"column:is_public_course;default:false" json:"is_public_course"`
	MaximumStudent  int32           `gorm:"column:maximum_student;default:0" json:"maximum_student"`
	TenantID        uint            `gorm:"column:tenant_id" json:"tenant_id"`
	Tenant          Tenant          `gorm:"foreignKey:TenantID;references:ID" json:"Tenant"`
	// CourseCurriculums []CourseCurriculum `gorm:"foreignKey:CourseID;references:ID" json:"course_curriculums"`
}

// type CourseCurriculum struct {
// 	ID          uint    `gorm:"primaryKey;autoIncrement" json:"id"`
// 	Position    int32   `gorm:"default:0" json:"position"`
// 	CourseID    uint    `gorm:"column:course_id" json:"course_id"`
// 	Course      Course  `gorm:"foreignKey:CourseID;references:ID" json:"course"`
// 	Title       string  `json:"title"`
// 	Description *string `gorm:"type:text" json:"description"`
// }

type CourseOverview struct {
	ID                uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	CourseID          uint    `gorm:"column:course_id" json:"course_id"`
	Course            Course  `gorm:"foreignKey:CourseID;references:ID" json:"course"`
	Outcomes          string  `json:"outcomes"`
	TargetAudience    *string `json:"target_audience"`
	DurationHours     *int32  `json:"duration_hours"`
	DurationMins      *int32  `json:"duration_mins"`
	MaterialsIncluded *string `json:"materials_included"`
	Requirements      *string `json:"requirements"`
}

// create course struct
type CreateCourseInput struct {
	Title           string              `json:"title" binding:"required"`
	Description     *string             `json:"description" binding:"omitempty"`
	Visibility      Visibility          `json:"visibility" binding:"required,oneof=public private protected"`
	IsScheduled     bool                `json:"is_scheduled"`
	ScheduleDate    *time.Time          `json:"schedule_date"`
	ScheduleTime    *time.Time          `json:"schedule_time"`
	ShowCommingSoon bool                `json:"show_comming_soon"`
	PricingModel    PricingModel        `json:"pricing_model" binding:"omitempty"`
	Tags            *[]string           `json:"tags" binding:"omitempty"`
	AuthorID        uint                `json:"author_id"`
	DifficultyLevel DifficultyLevel     `json:"difficulty_level" binding:"omitempty"`
	IsPublicCourse  bool                `json:"is_public_course" binding:"omitempty"`
	MaximumStudent  int32               `json:"maximum_student" binding:"omitempty"`
	CourseOverview  CourseOverviewInput `json:"course_overview"`
	// CourseCurriculums []CourseCurriculumInput `json:"course_curriculums"`
	// Author            user.User               `json:"author"`
	// FeaturedImage     *string                 `json:"featured_image"`
	// IntroVideo        *string                 `json:"intro_video"`
}

// type CourseCurriculumInput struct {
// 	Title       string  `json:"title" binding:"required"`
// 	Description *string `json:"description" binding:"omitempty"`
// }

type CourseOverviewInput struct {
	Outcomes          string  `json:"outcomes" binding:"required"`
	TargetAudience    *string `json:"target_audience" binding:"omitempty"`
	DurationHours     *int32  `json:"duration_hours" binding:"omitempty"`
	DurationMins      *int32  `json:"duration_mins" binding:"omitempty"`
	MaterialsIncluded *string `json:"materials_included" binding:"omitempty"`
	Requirements      *string `json:"requirements" binding:"omitempty"`
}
