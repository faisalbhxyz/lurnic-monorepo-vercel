package models

import "time"

type OrderPaymentStatus string

const (
	OrderPaymentStatusPaid   OrderPaymentStatus = "paid"
	OrderPaymentStatusUnpaid OrderPaymentStatus = "unpaid"
	// OrderPaymentStatusPartial   OrderPaymentStatus = "partial"
	// OrderPaymentStatusRefunded  OrderPaymentStatus = "refunded"
	// OrderPaymentStatusFailed    OrderPaymentStatus = "failed"
	// OrderPaymentStatusCancelled OrderPaymentStatus = "cancelled"
)

type Order struct {
	ID            uint               `json:"id" gorm:"primaryKey;autoIncrement"`
	StudentID     uint               `json:"student_id" gorm:"not null"`
	CourseID      uint               `json:"course_id" gorm:"not null"`
	DiscountType  string             `json:"discount_type" gorm:"type:varchar(50);default:'none'"`
	Discount      float64            `json:"discount" gorm:"default:0"`
	Total         float64            `json:"total" gorm:"not null"`
	PaymentStatus OrderPaymentStatus `json:"payment_status" gorm:"type:varchar(20);enum:paid,unpaid;default:unpaid"`
	InvoiceID     int64              `json:"invoice_id" gorm:"not null;index"`
	PaymentType   string             `json:"payment_type" gorm:"type:varchar(50);default:'manual'"`
	CustomerNote  *string            `json:"customer_note" gorm:"type:text;default:''"`
	AdminNote     *string            `json:"admin_note" gorm:"type:text;default:''"`
	PaymentMethod *string            `json:"payment_method" gorm:"type:varchar(255);null"`
	TransactionID *string            `json:"transaction_id" gorm:"type:varchar(255);null"`
	CreatedAt     time.Time          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time          `gorm:"autoUpdateTime" json:"updated_at"`
	TenantID      uint               `gorm:"column:tenant_id" json:"-"`
	Tenant        Tenant             `gorm:"foreignKey:TenantID;references:ID" json:"-"`
}
