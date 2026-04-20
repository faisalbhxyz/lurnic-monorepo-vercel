package models

import "time"

type GeneralSettings struct {
	ID            uint      `json:"id" form:"id" gorm:"primaryKey;autoIncrement"`
	OrgName       string    `json:"org_name" form:"org_name" gorm:"type:varchar(100);default:'Lurnic'"`
	Logo          *string   `json:"logo" form:"logo" gorm:"type:varchar(255);null"`
	Favicon       *string   `json:"favicon" form:"favicon" gorm:"type:varchar(255);null"`
	StudentPrefix string    `json:"student_prefix" form:"student_prefix" gorm:"type:varchar(10);default:'S-'"`
	TeacherPrefix string    `json:"teacher_prefix" form:"teacher_prefix" gorm:"type:varchar(10);default:'T-'"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	TenantID      uint      `gorm:"column:tenant_id" json:"-"`
}
