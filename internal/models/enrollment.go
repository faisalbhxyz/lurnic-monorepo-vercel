package models

import (
	"time"
)

type Enrollment struct {
	ID        uint          `gorm:"primaryKey;autoIncrement" json:"id"`
	CourseID  uint          `gorm:"column:course_id" json:"course_id"`
	Course    CourseDetails `gorm:"foreignKey:ID;references:CourseID" json:"course"`
	StudentID uint          `gorm:"column:student_id" json:"student_id"`
	Student   Student       `gorm:"foreignKey:ID;references:StudentID" json:"student"`
	CreatedAt time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	TenantID  uint          `gorm:"column:tenant_id" json:"-"`
	Tenant    Tenant        `gorm:"foreignKey:TenantID;references:ID" json:"-"`
}
