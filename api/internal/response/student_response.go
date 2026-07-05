package response

import (
	"dashlearn/internal/models"
	"time"
)

type StudentEnrollmentDetail struct {
	ID                   uint      `json:"id"`
	CourseID             uint      `json:"course_id"`
	Title                string    `json:"title"`
	Slug                 string    `json:"slug"`
	FeaturedImage        *string   `json:"featured_image"`
	EnrolledAt           time.Time `json:"enrolled_at"`
	ProgressPercent      float32   `json:"progress_percent"`
	LessonsCompleted     int64     `json:"lessons_completed"`
	LessonsTotal         int64     `json:"lessons_total"`
	QuizzesCompleted     int64     `json:"quizzes_completed"`
	QuizzesTotal         int64     `json:"quizzes_total"`
	AssignmentsCompleted int64     `json:"assignments_completed"`
	AssignmentsTotal     int64     `json:"assignments_total"`
}

type StudentDetailsStats struct {
	TotalEnrollments       int     `json:"total_enrollments"`
	AverageProgress        float32 `json:"average_progress"`
	QuizzesSubmitted       int64   `json:"quizzes_submitted"`
	AssignmentsSubmitted   int64   `json:"assignments_submitted"`
	TotalOrders            int     `json:"total_orders"`
	PaidOrders             int     `json:"paid_orders"`
	UnpaidOrders           int     `json:"unpaid_orders"`
	TotalSpent             float64 `json:"total_spent"`
}

type StudentOrderDetail struct {
	ID              uint                      `json:"id"`
	InvoiceID       int64                     `json:"invoice_id"`
	CourseID        uint                      `json:"course_id"`
	CourseTitle     string                    `json:"course_title"`
	FeaturedImage   *string                   `json:"featured_image"`
	Total           float64                   `json:"total"`
	Discount        float64                   `json:"discount"`
	DiscountType    string                    `json:"discount_type"`
	PaymentStatus   models.OrderPaymentStatus `json:"payment_status"`
	PaymentType     string                    `json:"payment_type"`
	PaymentMethod   *string                   `json:"payment_method"`
	TransactionID   *string                   `json:"transaction_id"`
	CustomerNote    *string                   `json:"customer_note"`
	AdminNote       *string                   `json:"admin_note"`
	OrderedAt       time.Time                 `json:"ordered_at"`
	UpdatedAt       time.Time                 `json:"updated_at"`
}

type StudentDetailsAdminResponse struct {
	ID          uint                      `json:"id"`
	UserID      string                    `json:"user_id"`
	FirstName   string                    `json:"first_name"`
	LastName    *string                   `json:"last_name"`
	Phone        *string                   `json:"phone"`
	Email        string                    `json:"email"`
	ProfileImage *string                   `json:"profile_image"`
	Status       bool                      `json:"status"`
	CreatedAt   time.Time                 `json:"created_at"`
	UpdatedAt   time.Time                 `json:"updated_at"`
	Enrollments []StudentEnrollmentDetail `json:"enrollments"`
	Orders      []StudentOrderDetail      `json:"orders"`
	Stats       StudentDetailsStats       `json:"stats"`
}
