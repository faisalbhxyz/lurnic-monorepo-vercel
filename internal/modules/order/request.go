package order

type CreateOrderInput struct {
	CourseID     uint   `json:"course_id" binding:"required"`
	CustomerNote string `json:"customer_note" binding:"omitempty"`
}
