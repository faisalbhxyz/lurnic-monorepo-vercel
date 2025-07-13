package response

import (
	"time"
)

type EnrollmentWithCourse struct {
	ID    uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Title string `json:"title"`
}

func (EnrollmentWithCourse) TableName() string {
	return "course_details"
}

type EnrollmentWithStudent struct {
	ID        uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName string  `json:"first_name"`
	LastName  *string `json:"last_name"`
	Email     string  `json:"email"`
}

func (EnrollmentWithStudent) TableName() string {
	return "students"
}

type EnrollmentResponse struct {
	ID        uint                  `gorm:"primaryKey;autoIncrement" json:"id"`
	CourseID  uint                  `gorm:"column:course_id" json:"course_id"`
	Course    EnrollmentWithCourse  `gorm:"foreignKey:ID;references:CourseID" json:"course"`
	StudentID uint                  `gorm:"column:student_id" json:"student_id"`
	Student   EnrollmentWithStudent `gorm:"foreignKey:ID;references:StudentID" json:"student"`
	CreatedAt time.Time             `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time             `gorm:"autoUpdateTime" json:"updated_at"`
}

func (EnrollmentResponse) TableName() string {
	return "enrollments"
}

type EnrolledCourseRes struct {
	ID        uint                        `json:"id"`
	CourseID  uint                        `json:"course_id"`
	Course    CourseDetailsPublicResponse `json:"course" gorm:"foreignKey:ID;references:CourseID"`
	StudentID uint                        `json:"student_id"`
	CreatedAt time.Time                   `json:"created_at"`
}
