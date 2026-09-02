package models

import (
	"time"

	"gorm.io/datatypes"
)

type CourseReviewStatus string

const (
	CourseReviewStatusPublished CourseReviewStatus = "published"
	CourseReviewStatusHidden    CourseReviewStatus = "hidden"
)

type CourseReview struct {
	ID        uint               `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID  uint               `gorm:"column:tenant_id" json:"-"`
	CourseID  uint               `gorm:"column:course_id" json:"course_id"`
	StudentID uint               `gorm:"column:student_id" json:"student_id"`
	Rating    uint8              `json:"rating"`
	Comment   *string            `gorm:"type:text" json:"comment"`
	Tags      datatypes.JSON     `gorm:"type:json" json:"tags"`
	Status    CourseReviewStatus `gorm:"type:enum('published','hidden');default:'published'" json:"status"`
	CreatedAt time.Time          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time          `gorm:"autoUpdateTime" json:"updated_at"`
	Student   Student            `gorm:"foreignKey:StudentID;references:ID" json:"student,omitempty"`
}

func (CourseReview) TableName() string {
	return "course_reviews"
}
