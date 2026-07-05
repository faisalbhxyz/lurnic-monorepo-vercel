package student

import (
	"dashlearn/internal/models"
	"dashlearn/internal/progress"
	"dashlearn/internal/response"
	"math"

	"gorm.io/gorm"
)

func buildStudentAdminDetails(db *gorm.DB, tenantID, studentID uint) (*response.StudentDetailsAdminResponse, error) {
	var student models.Student
	if err := db.Where("id = ? AND tenant_id = ?", studentID, tenantID).First(&student).Error; err != nil {
		return nil, err
	}

	var enrollments []models.Enrollment
	if err := db.Preload("Course").
		Where("student_id = ? AND tenant_id = ?", studentID, tenantID).
		Order("created_at DESC").
		Find(&enrollments).Error; err != nil {
		return nil, err
	}

	var quizSubmissions int64
	db.Model(&models.QuizSubmission{}).
		Where("tenant_id = ? AND student_id = ?", tenantID, studentID).
		Count(&quizSubmissions)

	var assignmentSubmissions int64
	db.Model(&models.AssignmentSubmission{}).
		Where("tenant_id = ? AND student_id = ?", tenantID, studentID).
		Count(&assignmentSubmissions)

	enrollmentDetails := make([]response.StudentEnrollmentDetail, 0, len(enrollments))
	var progressSum float32

	for _, enrollment := range enrollments {
		opts := progress.LoadOptions(db, enrollment.CourseID)
		breakdown := progress.CalcBreakdown(db, tenantID, studentID, enrollment.CourseID, opts, false)
		progressSum += breakdown.Percent

		enrollmentDetails = append(enrollmentDetails, response.StudentEnrollmentDetail{
			ID:                   enrollment.ID,
			CourseID:             enrollment.CourseID,
			Title:                enrollment.Course.Title,
			Slug:                 enrollment.Course.Slug,
			FeaturedImage:        enrollment.Course.FeaturedImage,
			EnrolledAt:           enrollment.CreatedAt,
			ProgressPercent:      breakdown.Percent,
			LessonsCompleted:     breakdown.LessonsDone,
			LessonsTotal:         breakdown.LessonsTotal,
			QuizzesCompleted:     breakdown.QuizzesDone,
			QuizzesTotal:         breakdown.QuizzesTotal,
			AssignmentsCompleted: breakdown.AssignmentsDone,
			AssignmentsTotal:     breakdown.AssignmentsTotal,
		})
	}

	avgProgress := float32(0)
	if len(enrollmentDetails) > 0 {
		avgProgress = float32(math.Round(float64(progressSum/float32(len(enrollmentDetails)))*10) / 10)
	}

	orderDetails, orderStats := loadStudentOrders(db, tenantID, studentID)

	var activeDevice *response.StudentActiveDevice
	if session := getStudentSessionForAdmin(db, studentID); session != nil {
		deviceName := "Unknown device"
		if session.DeviceName != nil && *session.DeviceName != "" {
			deviceName = *session.DeviceName
		}
		activeDevice = &response.StudentActiveDevice{
			DeviceID:   session.DeviceID,
			DeviceName: deviceName,
			IPAddress:  session.IPAddress,
			UserAgent:  session.UserAgent,
			LoggedInAt: session.CreatedAt,
			LastSeenAt: session.LastSeenAt,
		}
	}

	return &response.StudentDetailsAdminResponse{
		ID:           student.ID,
		UserID:       student.UserID,
		FirstName:    student.FirstName,
		LastName:     student.LastName,
		Phone:        student.Phone,
		Email:        student.Email,
		ProfileImage: student.ProfileImage,
		Status:       student.Status,
		CreatedAt:    student.CreatedAt,
		UpdatedAt:    student.UpdatedAt,
		Enrollments: enrollmentDetails,
		Orders:      orderDetails,
		Stats: response.StudentDetailsStats{
			TotalEnrollments:     len(enrollmentDetails),
			AverageProgress:      avgProgress,
			QuizzesSubmitted:     quizSubmissions,
			AssignmentsSubmitted: assignmentSubmissions,
			TotalOrders:          orderStats.total,
			PaidOrders:           orderStats.paid,
			UnpaidOrders:         orderStats.unpaid,
			TotalSpent:           orderStats.spent,
		},
		ActiveDevice: activeDevice,
	}, nil
}

type studentOrderStats struct {
	total  int
	paid   int
	unpaid int
	spent  float64
}

func loadStudentOrders(db *gorm.DB, tenantID, studentID uint) ([]response.StudentOrderDetail, studentOrderStats) {
	var rows []struct {
		models.Order
		CourseTitle   string  `gorm:"column:course_title"`
		FeaturedImage *string `gorm:"column:featured_image"`
	}

	err := db.Table("orders").
		Select("orders.*, course_details.title AS course_title, course_details.featured_image").
		Joins("LEFT JOIN course_details ON course_details.id = orders.course_id").
		Where("orders.student_id = ? AND orders.tenant_id = ?", studentID, tenantID).
		Order("orders.created_at DESC").
		Scan(&rows).Error
	if err != nil {
		return []response.StudentOrderDetail{}, studentOrderStats{}
	}

	stats := studentOrderStats{total: len(rows)}
	orders := make([]response.StudentOrderDetail, 0, len(rows))

	for _, row := range rows {
		if row.PaymentStatus == models.OrderPaymentStatusPaid {
			stats.paid++
			stats.spent += row.Total
		} else {
			stats.unpaid++
		}

		orders = append(orders, response.StudentOrderDetail{
			ID:            row.ID,
			InvoiceID:     row.InvoiceID,
			CourseID:      row.CourseID,
			CourseTitle:   row.CourseTitle,
			FeaturedImage: row.FeaturedImage,
			Total:         row.Total,
			Discount:      row.Discount,
			DiscountType:  row.DiscountType,
			PaymentStatus: row.PaymentStatus,
			PaymentType:   row.PaymentType,
			PaymentMethod: row.PaymentMethod,
			TransactionID: row.TransactionID,
			CustomerNote:  row.CustomerNote,
			AdminNote:     row.AdminNote,
			OrderedAt:     row.CreatedAt,
			UpdatedAt:     row.UpdatedAt,
		})
	}

	return orders, stats
}
