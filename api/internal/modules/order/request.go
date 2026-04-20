package order

type CreateOrderInput struct {
	CourseID      uint    `json:"course_id" binding:"required"`
	CustomerNote  string  `json:"customer_note" binding:"omitempty"`
	PaymentMethod *string `json:"payment_method" binding:"omitempty"`
	TransactionID *string `json:"transaction_id" binding:"omitempty"`
}
