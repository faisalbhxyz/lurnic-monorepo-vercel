package response

import (
	"dashlearn/models"
	"time"
)

type CreateOrderResponse struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	CourseID     uint      `json:"course_id" gorm:"not null"`
	Total        float64   `json:"total" gorm:"not null"`
	InvoiceID    int64     `json:"invoice_id" gorm:"not null;index"`
	CustomerNote string    `json:"customer_note" gorm:"type:text;default:''"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (CreateOrderResponse) TableName() string {
	return "orders"
}

type OrderWithStudent struct {
	ID        uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    string  `gorm:"column:user_id;uniqueIndex" json:"user_id"`
	FirstName string  `json:"first_name"`
	LastName  *string `json:"last_name"`
	Phone     *string `json:"phone"`
	Email     string  `gorm:"uniqueIndex" json:"email"`
}

func (OrderWithStudent) TableName() string {
	return "students"
}

type OrderWithCourse struct {
	ID    uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Title string `json:"title"`
}

func (OrderWithCourse) TableName() string {
	return "course_details"
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

func (GetAllOrderResponse) TableName() string {
	return "orders"
}
