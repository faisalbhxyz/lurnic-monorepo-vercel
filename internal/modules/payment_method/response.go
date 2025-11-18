package paymentmethod

import "time"

type PaymentMethodResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Image       *string   `json:"image"`
	Instruction string    `json:"instruction"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
