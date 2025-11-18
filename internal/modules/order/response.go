package order

import (
	"dashlearn/internal/models"
	"time"
)

type OrderWithStudent struct {
	ID        uint    `json:"id"`
	UserID    string  `json:"user_id"`
	FirstName string  `json:"first_name"`
	LastName  *string `json:"last_name"`
	Phone     *string `json:"phone"`
	Email     string  `json:"email"`
}

type OrderWithCourse struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

type GetAllOrderResponse struct {
	ID            uint                      `json:"id" gorm:"primaryKey;autoIncrement"`
	StudentID     uint                      `json:"student_id" gorm:"not null"`
	Student       OrderWithStudent          `gorm:"foreignKey:StudentID;references:ID" json:"student"`
	CourseID      uint                      `json:"course_id" gorm:"not null"`
	Course        OrderWithCourse           `gorm:"foreignKey:CourseID;references:ID" json:"course"`
	DiscountType  string                    `json:"discount_type" gorm:"type:varchar(50);default:'none'"`
	Discount      float64                   `json:"discount" gorm:"default:0"`
	Total         float64                   `json:"total" gorm:"not null"`
	PaymentStatus models.OrderPaymentStatus `json:"payment_status" gorm:"type:varchar(20);enum:paid,unpaid;default:unpaid"`
	InvoiceID     int64                     `json:"invoice_id" gorm:"not null;index"`
	PaymentType   string                    `json:"payment_type" gorm:"type:varchar(50);default:'manual'"`
	CustomerNote  string                    `json:"customer_note" gorm:"type:text;default:''"`
	AdminNote     string                    `json:"admin_note" gorm:"type:text;default:''"`
	CreatedAt     time.Time                 `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time                 `gorm:"autoUpdateTime" json:"updated_at"`
}
