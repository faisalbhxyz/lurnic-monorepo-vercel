package response

type InstructorResponse struct {
	ID          uint    `json:"id"`
	FirstName   string  `json:"first_name"`
	LastName    *string `json:"last_name"`
	Phone       *string `json:"phone"`
	Email       string  `json:"email"`
	Image       *string `json:"image"`
	Role        *string `json:"role"`
	Designation *string `json:"designation"`
}

type CourseInstructorResponse struct {
	ID           uint               `json:"id"`
	CourseID     uint               `json:"course_id"`
	InstructorID uint               `json:"instructor_id"`
	Instructor   InstructorResponse `json:"instructor"`
}

type InstructorDetailsResponse struct {
	ID          uint    `json:"id"`
	FirstName   string  `json:"first_name"`
	LastName    *string `json:"last_name"`
	Email       string  `json:"email"`
	Phone       *string `json:"phone"`
	Role        *string `json:"role"`
	Designation *string `json:"designation"`
	Image       *string `json:"image"`
}
