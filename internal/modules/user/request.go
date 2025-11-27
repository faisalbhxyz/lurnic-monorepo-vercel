package user

type CreateUserInput struct {
	Name     string  `json:"name" binding:"required"`
	Phone    *string `json:"phone" binding:"omitempty"`
	Email    string  `json:"email" binding:"required,email"`
	Password string  `json:"password" binding:"required,min=6"`
}

type LoginUserInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type CreateTeamMemberInput struct {
	UserID   string  `json:"user_id" binding:"required"`
	Name     string  `json:"name" binding:"required"`
	Phone    *string `json:"phone" binding:"omitempty"`
	Email    string  `json:"email" binding:"required,email"`
	Role     int32   `json:"role" binding:"required"`
	Password string  `json:"password" binding:"required,min=6"`
}

type UpdateTeamMemberInput struct {
	Name  *string `json:"name"`
	Phone *string `json:"phone"`
	Role  *int32  `json:"role"`
}
