package models

import (
	"time"
)

type StudentWithEnrollmentResponse struct {
	ID        uint `gorm:"primaryKey;autoIncrement" json:"id"`
	StudentID uint `gorm:"column:student_id" json:"student_id"`
	CourseID  uint `gorm:"column:course_id" json:"course_id"`
}

func (StudentWithEnrollmentResponse) TableName() string {
	return "enrollments"
}

type Student struct {
	ID          uint                            `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      string                          `gorm:"column:user_id;uniqueIndex" json:"user_id"`
	FirstName   string                          `json:"first_name"`
	LastName    *string                         `json:"last_name"`
	Phone       *string                         `json:"phone"`
	Email       string                          `gorm:"uniqueIndex" json:"email"`
	Password    string                          `json:"-"`
	Status      bool                            `json:"status"`
	OTPCode     *string                         `json:"otp_code,omitempty"`
	CreatedAt   time.Time                       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time                       `gorm:"autoUpdateTime" json:"updated_at"`
	TenantID    uint                            `gorm:"column:tenant_id" json:"-"`
	Tenant      Tenant                          `gorm:"foreignKey:TenantID;references:ID" json:"-"`
	Enrollments []StudentWithEnrollmentResponse `gorm:"foreignKey:StudentID" json:"enrollments"`
}
