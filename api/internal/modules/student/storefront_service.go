package student

import (
	"errors"
	"math"
	"time"

	"dashlearn/internal/models"
	"dashlearn/internal/progress"
	"dashlearn/internal/response"

	"gorm.io/gorm"
)

type StorefrontService struct {
	db *gorm.DB
}

func NewStorefrontService(db *gorm.DB) *StorefrontService {
	return &StorefrontService{db: db}
}

type LearningReportResponse struct {
	Period              string                    `json:"period"`
	DailyWatchSeconds   []DailyWatchSecondsEntry  `json:"daily_watch_seconds"`
	StreakDays          int                       `json:"streak_days"`
	QuizAccuracyPercent float64                   `json:"quiz_accuracy_percent"`
	CoursesInProgress   int                       `json:"courses_in_progress"`
	CoursesCompleted    int                       `json:"courses_completed"`
}

type DailyWatchSecondsEntry struct {
	Date    string `json:"date"`
	Seconds int64  `json:"seconds"`
}

func (s *StorefrontService) GetLearningReport(tenantID, studentID uint, period string) (*LearningReportResponse, error) {
	days, err := parseReportPeriod(period)
	if err != nil {
		return nil, err
	}

	startDate := time.Now().AddDate(0, 0, -(days - 1)).Truncate(24 * time.Hour)

	dailyWatch, err := s.aggregateDailyWatchSeconds(tenantID, studentID, startDate)
	if err != nil {
		return nil, err
	}

	streakDays := calcStreakDays(dailyWatch)

	quizAccuracy, err := s.calcQuizAccuracy(tenantID, studentID, startDate)
	if err != nil {
		return nil, err
	}

	inProgress, completed, err := s.countCourseStates(tenantID, studentID)
	if err != nil {
		return nil, err
	}

	return &LearningReportResponse{
		Period:              period,
		DailyWatchSeconds:   dailyWatch,
		StreakDays:          streakDays,
		QuizAccuracyPercent: quizAccuracy,
		CoursesInProgress:   inProgress,
		CoursesCompleted:    completed,
	}, nil
}

func parseReportPeriod(period string) (int, error) {
	switch period {
	case "", "7d":
		return 7, nil
	case "30d":
		return 30, nil
	case "90d":
		return 90, nil
	default:
		return 0, errors.New("invalid period")
	}
}

func (s *StorefrontService) aggregateDailyWatchSeconds(tenantID, studentID uint, startDate time.Time) ([]DailyWatchSecondsEntry, error) {
	type row struct {
		Date    string
		Seconds float64
	}

	var rows []row
	err := s.db.Model(&models.StudentLessonVideoProgress{}).
		Select("DATE(updated_at) as date, SUM(max_position_seconds) as seconds").
		Where("tenant_id = ? AND student_id = ? AND updated_at >= ?", tenantID, studentID, startDate).
		Group("DATE(updated_at)").
		Order("date ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	byDate := make(map[string]int64, len(rows))
	for _, row := range rows {
		byDate[row.Date] = int64(math.Round(row.Seconds))
	}

	days := int(time.Since(startDate).Hours()/24) + 1
	if days < 1 {
		days = 1
	}

	result := make([]DailyWatchSecondsEntry, 0, days)
	for i := 0; i < days; i++ {
		day := startDate.AddDate(0, 0, i).Format("2006-01-02")
		result = append(result, DailyWatchSecondsEntry{
			Date:    day,
			Seconds: byDate[day],
		})
	}

	return result, nil
}

func calcStreakDays(daily []DailyWatchSecondsEntry) int {
	streak := 0
	for i := len(daily) - 1; i >= 0; i-- {
		if daily[i].Seconds > 0 {
			streak++
			continue
		}
		break
	}
	return streak
}

func (s *StorefrontService) calcQuizAccuracy(tenantID, studentID uint, startDate time.Time) (float64, error) {
	var submissions []models.QuizSubmission
	err := s.db.
		Where("tenant_id = ? AND student_id = ? AND submitted_at >= ?", tenantID, studentID, startDate).
		Find(&submissions).Error
	if err != nil {
		return 0, err
	}
	if len(submissions) == 0 {
		return 0, nil
	}

	var totalPercent float64
	for _, sub := range submissions {
		totalPercent += float64(sub.Percentage)
	}
	return math.Round(totalPercent/float64(len(submissions))*10) / 10, nil
}

func (s *StorefrontService) countCourseStates(tenantID, studentID uint) (inProgress int, completed int, err error) {
	var enrollments []models.Enrollment
	if err = s.db.Where("tenant_id = ? AND student_id = ?", tenantID, studentID).Find(&enrollments).Error; err != nil {
		return 0, 0, err
	}

	for _, enrollment := range enrollments {
		opts := progress.LoadOptions(s.db, enrollment.CourseID)
		percent := progress.CalcCourseProgress(s.db, tenantID, studentID, enrollment.CourseID, opts)
		if percent >= 100 {
			completed++
		} else if percent > 0 {
			inProgress++
		}
	}

	return inProgress, completed, nil
}

func (s *StorefrontService) ListNotifications(tenantID, studentID uint) ([]response.StudentNotificationResponse, error) {
	var rows []models.StudentNotification
	if err := s.db.
		Where("tenant_id = ? AND student_id = ?", tenantID, studentID).
		Order("created_at DESC").
		Limit(100).
		Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]response.StudentNotificationResponse, 0, len(rows))
	for _, row := range rows {
		result = append(result, response.StudentNotificationResponse{
			ID:        row.ID,
			Title:     row.Title,
			Body:      row.Body,
			Type:      row.Type,
			Data:      row.Data,
			ReadAt:    row.ReadAt,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
		})
	}
	return result, nil
}

func (s *StorefrontService) MarkNotificationRead(tenantID, studentID, notificationID uint) error {
	res := s.db.Model(&models.StudentNotification{}).
		Where("id = ? AND tenant_id = ? AND student_id = ?", notificationID, tenantID, studentID).
		Updates(map[string]interface{}{"read_at": time.Now()})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (s *StorefrontService) ListOrders(tenantID, studentID uint) ([]response.StudentOrderDetail, error) {
	var orders []models.Order
	if err := s.db.
		Where("tenant_id = ? AND student_id = ?", tenantID, studentID).
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return []response.StudentOrderDetail{}, nil
	}

	courseIDs := make([]uint, 0, len(orders))
	for _, order := range orders {
		courseIDs = append(courseIDs, order.CourseID)
	}

	var courses []models.CourseDetails
	if err := s.db.Select("id", "title", "featured_image").
		Where("tenant_id = ? AND id IN ?", tenantID, courseIDs).
		Find(&courses).Error; err != nil {
		return nil, err
	}

	courseByID := make(map[uint]models.CourseDetails, len(courses))
	for _, course := range courses {
		courseByID[course.ID] = course
	}

	result := make([]response.StudentOrderDetail, 0, len(orders))
	for _, order := range orders {
		course := courseByID[order.CourseID]
		result = append(result, response.StudentOrderDetail{
			ID:            order.ID,
			InvoiceID:     order.InvoiceID,
			CourseID:      order.CourseID,
			CourseTitle:   course.Title,
			FeaturedImage: course.FeaturedImage,
			Total:         order.Total,
			Discount:      order.Discount,
			DiscountType:  order.DiscountType,
			PaymentStatus: order.PaymentStatus,
			PaymentType:   order.PaymentType,
			PaymentMethod: order.PaymentMethod,
			TransactionID: order.TransactionID,
			CustomerNote:  order.CustomerNote,
			AdminNote:     order.AdminNote,
			OrderedAt:     order.CreatedAt,
			UpdatedAt:     order.UpdatedAt,
		})
	}

	return result, nil
}
