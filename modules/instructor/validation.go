package instructor

type CreateInstructorInput struct {
	FirstName string  `json:"first_name" binding:"required"`
	LastName  *string `json:"last_name" binding:"omitempty"`
	Phone     *string `json:"phone" binding:"omitempty"`
	Email     string  `json:"email" binding:"required,email"`
	Password  string  `json:"password" binding:"required,min=6"`
}

type LoginInstructorInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}
