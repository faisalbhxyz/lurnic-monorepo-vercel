package paymentmethod

type CreatePaymentMethodInput struct {
	Title       string  `json:"title" form:"title" binding:"required"`
	Image       *string `json:"image"`
	Instruction string  `json:"instruction" form:"instruction" binding:"required"`
}
