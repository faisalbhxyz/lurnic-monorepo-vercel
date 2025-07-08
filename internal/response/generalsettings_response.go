package response

import "time"

type GeneralSettingsRes struct {
	ID            uint      `json:"id" form:"id"`
	OrgName       string    `json:"org_name" form:"org_name"`
	Logo          *string   `json:"logo" form:"logo"`
	Favicon       *string   `json:"favicon" form:"favicon"`
	StudentPrefix string    `json:"student_prefix" form:"student_prefix"`
	TeacherPrefix string    `json:"teacher_prefix" form:"teacher_prefix"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
