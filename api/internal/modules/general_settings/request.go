package generalsettings

type CreateOrUpdateGeneralSettingsInput struct {
	OrgName       string `json:"org_name" form:"org_name" binding:"required,min=3,max=100"`
	Logo          *string
	Favicon       *string
	StudentPrefix string `json:"student_prefix" form:"student_prefix" binding:"required,min=1,max=10"`
	TeacherPrefix string `json:"teacher_prefix" form:"teacher_prefix" binding:"required,min=1,max=10"`
}

type UpdateCurrencyInput struct {
	Currency string `json:"currency" binding:"required,min=3,max=10"`
}
