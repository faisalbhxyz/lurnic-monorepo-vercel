package instructor

type CreateInstructorInput struct {
	FirstName   string  `json:"first_name" form:"first_name" binding:"required"`
	LastName    *string `json:"last_name" form:"last_name" binding:"omitempty"`
	Phone       *string `json:"phone" form:"phone" binding:"omitempty"`
	Role        *string `json:"role" form:"role" binding:"omitempty"`
	Designation *string `json:"designation" form:"designation" binding:"omitempty"`
	Email       string  `json:"email" form:"email" binding:"required,email"`
	Password    string  `json:"password" form:"password" binding:"required,min=6"`
}

type LoginInstructorInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}
