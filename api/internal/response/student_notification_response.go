package response

import "time"

type StudentNotificationResponse struct {
	ID        uint       `json:"id"`
	Title     string     `json:"title"`
	Body      *string    `json:"body,omitempty"`
	Type      string     `json:"type"`
	Data      []byte     `json:"data,omitempty"`
	ReadAt    *time.Time `json:"read_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
