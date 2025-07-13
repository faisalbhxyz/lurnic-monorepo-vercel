package response

import (
	"dashlearn/internal/models"
	"dashlearn/internal/utils"
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
	Enrollments     []EnrollmentResponse         `gorm:"foreignKey:CourseID;references:ID" json:"enrollments"`
}

func (CourseDetailsResponse) TableName() string {
	return "course_details"
}

// COURSES RESPONSE without chapters, instructors and enrollments
type CourseDetailsPublicResponse struct {
	ID              uint                            `json:"id,omitempty"`
	Title           string                          `json:"title,omitempty"`
	Slug            string                          `json:"slug,omitempty"`
	Summary         string                          `json:"summary,omitempty"`
	Description     *string                         `json:"description,omitempty"`
	Visibility      models.Visibility               `json:"visibility,omitempty"`
	IsScheduled     *bool                           `json:"is_scheduled,omitempty"`
	ScheduleDate    *time.Time                      `json:"schedule_date,omitempty"`
	ScheduleTime    *time.Time                      `json:"schedule_time,omitempty"`
	FeaturedImage   *string                         `json:"featured_image,omitempty"`
	IntroVideo      *utils.JSONB[models.IntroVideo] `json:"intro_video,omitempty"`
	PricingModel    models.CoursePricingModel       `json:"pricing_model,omitempty"`
	RegularPrice    *float32                        `json:"regular_price,omitempty"`
	SalePrice       *float32                        `json:"sale_price,omitempty"`
	ShowCommingSoom *bool                           `json:"show_comming_soom,omitempty"`
	Tags            datatypes.JSON                  `json:"tags,omitempty"`
	Overview        datatypes.JSON                  `json:"overview,omitempty"`
	GeneralSettings *CourseGeneralSettingsResponse  `json:"general_settings,omitempty"`
	Chapters        []CourseChapterResponse         `json:"course_chapters,omitempty"`
	Instructors     []CourseInstructorResponse      `json:"course_instructors,omitempty"`
	Enrollments     []EnrolledCourseRes             `json:"enrollments,omitempty"`
}
