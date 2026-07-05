package student

type CreateStudentInput struct {
	FirstName        string  `json:"first_name" form:"first_name" binding:"required"`
	LastName         *string `json:"last_name" form:"last_name" binding:"omitempty"`
	Phone            *string `json:"phone" form:"phone" binding:"omitempty"`
	Email            string  `json:"email" form:"email" binding:"required,email"`
	Password         string  `json:"password" form:"password" binding:"required,min=6"`
	ProfileImageURL  *string
}

type UpdateStudentInput struct {
	FirstName       string  `json:"first_name" form:"first_name" binding:"required"`
	LastName        *string `json:"last_name" form:"last_name" binding:"omitempty"`
	Phone           *string `json:"phone" form:"phone" binding:"omitempty"`
	ProfileImageURL *string
}

type LoginStudentInput struct {
	Email      string  `json:"email" binding:"required,email"`
	Password   string  `json:"password" binding:"required,min=6"`
	DeviceID   string  `json:"device_id" binding:"required,min=8,max=128"`
	DeviceName *string `json:"device_name" binding:"omitempty,max=255"`
}

type ForgotPasswordInput struct {
	Email    string `json:"email" binding:"required,email"`
	ResetURL string `json:"reset_url" binding:"required,url"`
}

type ResetPasswordInput struct {
	Email    string `json:"email" binding:"required,email"`
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}
