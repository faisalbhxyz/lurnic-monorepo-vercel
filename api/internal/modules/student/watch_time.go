package student

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"dashlearn/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	maxWatchedSecondsPerEvent = 300
	maxBatchWatchEvents       = 50
	watchSourceEnrolled       = "enrolled"
	watchSourceFreeLesson     = "free_lesson"
	watchSourceOffline        = "offline"
)

var (
	errInvalidWatchSeconds = errors.New("watched_seconds must be greater than 0")
	errInvalidWatchDate    = errors.New("invalid watch_date")
	errFutureWatchDate     = errors.New("watch_date is too far in the future")
	errInvalidTimezone     = errors.New("invalid timezone")
	errInvalidClientEvent  = errors.New("client_event_id is required")
	errInvalidWatchSource  = errors.New("invalid source")
	errLessonNotPlayable   = errors.New("lesson not playable")
	errBatchTooLarge       = errors.New("batch exceeds 50 events")
	errEmptyBatch          = errors.New("events required")
)

type WatchTimeInput struct {
	ClientEventID  string  `json:"client_event_id" binding:"required"`
	WatchedSeconds int     `json:"watched_seconds" binding:"required"`
	WatchDate      string  `json:"watch_date" binding:"required"`
	Timezone       string  `json:"timezone" binding:"required"`
	WatchedAt      *string `json:"watched_at"`
	CourseID       *uint   `json:"course_id"`
	LessonID       *uint   `json:"lesson_id"`
	Source         string  `json:"source"`
	DevicePlatform *string `json:"device_platform"`
}

type WatchTimeBatchInput struct {
	Events []WatchTimeInput `json:"events" binding:"required"`
}

type WatchTimeAcceptResponse struct {
	Accepted         bool   `json:"accepted"`
	WatchDate        string `json:"watch_date"`
	DayVideoSeconds  int    `json:"day_video_seconds"`
	Duplicate        bool   `json:"duplicate"`
	ClientEventID    string `json:"client_event_id,omitempty"`
}

type WatchTimeBatchResponse struct {
	AcceptedCount  int                      `json:"accepted_count"`
	DuplicateCount int                      `json:"duplicate_count"`
	Results        []WatchTimeAcceptResponse `json:"results"`
	DailyTotals    []DailyWatchSecondsEntry  `json:"daily_totals"`
}

type LearningReportTotals struct {
	VideoSecondsPeriod  int64   `json:"video_seconds_period"`
	VideoSecondsAllTime int64   `json:"video_seconds_all_time"`
	LastWatchedAt       *string `json:"last_watched_at"`
}

type LearningReportByCourse struct {
	CourseID    uint   `json:"course_id"`
	CourseTitle string `json:"course_title"`
	VideoSeconds int64 `json:"video_seconds"`
}

type AdminLearningReportResponse struct {
	LearningReportResponse
	Totals   LearningReportTotals     `json:"totals"`
	ByCourse []LearningReportByCourse `json:"by_course"`
}

type LearningTimeSummary struct {
	VideoSeconds7d   int64   `json:"video_seconds_7d"`
	VideoSeconds30d  int64   `json:"video_seconds_30d"`
	StreakDays       int     `json:"streak_days"`
	LastWatchedAt    *string `json:"last_watched_at"`
}

type normalizedWatchEvent struct {
	ClientEventID  string
	WatchedSeconds int
	WatchDate      time.Time
	WatchDateStr   string
	Timezone       string
	WatchedAt      time.Time
	CourseID       *uint
	LessonID       *uint
	Source         string
	DevicePlatform *string
}

func (s *StorefrontService) IngestWatchTime(tenantID, studentID uint, input WatchTimeInput) (*WatchTimeAcceptResponse, error) {
	event, err := s.normalizeWatchInput(input)
	if err != nil {
		return nil, err
	}
	if err := s.ensureLessonPlayable(tenantID, studentID, event); err != nil {
		return nil, err
	}
	return s.acceptWatchEvent(tenantID, studentID, event)
}

func (s *StorefrontService) IngestWatchTimeBatch(tenantID, studentID uint, input WatchTimeBatchInput) (*WatchTimeBatchResponse, error) {
	if len(input.Events) == 0 {
		return nil, errEmptyBatch
	}
	if len(input.Events) > maxBatchWatchEvents {
		return nil, errBatchTooLarge
	}

	results := make([]WatchTimeAcceptResponse, 0, len(input.Events))
	accepted := 0
	duplicates := 0
	touchedDates := map[string]struct{}{}

	for _, raw := range input.Events {
		event, err := s.normalizeWatchInput(raw)
		if err != nil {
			return nil, err
		}
		if err := s.ensureLessonPlayable(tenantID, studentID, event); err != nil {
			return nil, err
		}
		res, err := s.acceptWatchEvent(tenantID, studentID, event)
		if err != nil {
			return nil, err
		}
		results = append(results, *res)
		touchedDates[res.WatchDate] = struct{}{}
		if res.Duplicate {
			duplicates++
		} else {
			accepted++
		}
	}

	dailyTotals, err := s.dayTotalsForDates(tenantID, studentID, touchedDates)
	if err != nil {
		return nil, err
	}

	return &WatchTimeBatchResponse{
		AcceptedCount:  accepted,
		DuplicateCount: duplicates,
		Results:        results,
		DailyTotals:    dailyTotals,
	}, nil
}

func (s *StorefrontService) GetAdminLearningReport(tenantID, studentID uint, period string) (*AdminLearningReportResponse, error) {
	var student models.Student
	if err := s.db.Where("id = ? AND tenant_id = ?", studentID, tenantID).First(&student).Error; err != nil {
		return nil, err
	}

	base, err := s.buildLearningReport(tenantID, studentID, period, true)
	if err != nil {
		return nil, err
	}

	totals, err := s.calcLearningTotals(tenantID, studentID, base.DailyWatchSeconds)
	if err != nil {
		return nil, err
	}

	byCourse, err := s.calcWatchByCourse(tenantID, studentID, period)
	if err != nil {
		return nil, err
	}

	return &AdminLearningReportResponse{
		LearningReportResponse: *base,
		Totals:                 *totals,
		ByCourse:               byCourse,
	}, nil
}

func (s *StorefrontService) GetLearningTimeSummary(tenantID, studentID uint) (*LearningTimeSummary, error) {
	loc := reportLocation()
	today := startOfLocalDay(time.Now(), loc)

	sum7, err := s.sumVideoSecondsSince(tenantID, studentID, today.AddDate(0, 0, -6))
	if err != nil {
		return nil, err
	}
	sum30, err := s.sumVideoSecondsSince(tenantID, studentID, today.AddDate(0, 0, -29))
	if err != nil {
		return nil, err
	}

	daily7, err := s.aggregateDailyWatchSeconds(tenantID, studentID, today.AddDate(0, 0, -6), 7, loc)
	if err != nil {
		return nil, err
	}

	lastAt, err := s.lastWatchedAt(tenantID, studentID)
	if err != nil {
		return nil, err
	}

	return &LearningTimeSummary{
		VideoSeconds7d:  sum7,
		VideoSeconds30d: sum30,
		StreakDays:      calcStreakDays(daily7),
		LastWatchedAt:   lastAt,
	}, nil
}

func (s *StorefrontService) normalizeWatchInput(input WatchTimeInput) (*normalizedWatchEvent, error) {
	clientEventID := strings.TrimSpace(input.ClientEventID)
	if clientEventID == "" || len(clientEventID) > 64 {
		return nil, errInvalidClientEvent
	}

	if input.WatchedSeconds <= 0 {
		return nil, errInvalidWatchSeconds
	}
	watched := input.WatchedSeconds
	if watched > maxWatchedSecondsPerEvent {
		watched = maxWatchedSecondsPerEvent
	}

	tzName := strings.TrimSpace(input.Timezone)
	if tzName == "" {
		return nil, errInvalidTimezone
	}
	loc, err := time.LoadLocation(tzName)
	if err != nil {
		return nil, errInvalidTimezone
	}

	watchDateStr := normalizeReportDate(input.WatchDate)
	watchDate, err := time.ParseInLocation("2006-01-02", watchDateStr, loc)
	if err != nil {
		return nil, errInvalidWatchDate
	}

	todayLocal := startOfLocalDay(time.Now(), loc)
	maxAllowed := todayLocal.AddDate(0, 0, 1)
	if watchDate.After(maxAllowed) {
		return nil, errFutureWatchDate
	}

	watchedAt := time.Now().UTC()
	if input.WatchedAt != nil && strings.TrimSpace(*input.WatchedAt) != "" {
		parsed, parseErr := time.Parse(time.RFC3339, strings.TrimSpace(*input.WatchedAt))
		if parseErr != nil {
			return nil, fmt.Errorf("invalid watched_at")
		}
		watchedAt = parsed.UTC()
	}

	source := strings.TrimSpace(input.Source)
	if source == "" {
		source = watchSourceEnrolled
	}
	switch source {
	case watchSourceEnrolled, watchSourceFreeLesson, watchSourceOffline:
	default:
		return nil, errInvalidWatchSource
	}

	var platform *string
	if input.DevicePlatform != nil {
		p := strings.ToLower(strings.TrimSpace(*input.DevicePlatform))
		if p != "" {
			if len(p) > 16 {
				p = p[:16]
			}
			platform = &p
		}
	}

	return &normalizedWatchEvent{
		ClientEventID:  clientEventID,
		WatchedSeconds: watched,
		WatchDate:      watchDate,
		WatchDateStr:   watchDateStr,
		Timezone:       tzName,
		WatchedAt:      watchedAt,
		CourseID:       input.CourseID,
		LessonID:       input.LessonID,
		Source:         source,
		DevicePlatform: platform,
	}, nil
}

func (s *StorefrontService) ensureLessonPlayable(tenantID, studentID uint, event *normalizedWatchEvent) error {
	if event.LessonID == nil {
		return nil
	}

	lessonID := *event.LessonID
	var lesson struct {
		ID        uint
		CourseID  uint
		IsPublic  bool
		Published bool
	}
	err := s.db.Table("course_lessons AS l").
		Select("l.id, ch.course_id AS course_id, l.is_public, l.is_published AS published").
		Joins("JOIN course_chapters AS ch ON ch.id = l.chapter_id").
		Joins("JOIN course_details AS c ON c.id = ch.course_id").
		Where("l.id = ? AND c.tenant_id = ?", lessonID, tenantID).
		Scan(&lesson).Error
	if err != nil {
		return err
	}
	if lesson.ID == 0 || !lesson.Published {
		return errLessonNotPlayable
	}

	if event.CourseID != nil && *event.CourseID != lesson.CourseID {
		return errLessonNotPlayable
	}
	courseID := lesson.CourseID
	event.CourseID = &courseID

	var enrolled int64
	if err := s.db.Model(&models.Enrollment{}).
		Where("tenant_id = ? AND student_id = ? AND course_id = ?", tenantID, studentID, lesson.CourseID).
		Count(&enrolled).Error; err != nil {
		return err
	}
	if enrolled > 0 {
		return nil
	}
	if lesson.IsPublic {
		return nil
	}
	return errLessonNotPlayable
}

func (s *StorefrontService) acceptWatchEvent(tenantID, studentID uint, event *normalizedWatchEvent) (*WatchTimeAcceptResponse, error) {
	var existing models.StudentWatchEvent
	err := s.db.Where("tenant_id = ? AND student_id = ? AND client_event_id = ?", tenantID, studentID, event.ClientEventID).
		First(&existing).Error
	if err == nil {
		daySeconds, dayErr := s.dayVideoSeconds(tenantID, studentID, event.WatchDate)
		if dayErr != nil {
			return nil, dayErr
		}
		return &WatchTimeAcceptResponse{
			Accepted:        true,
			WatchDate:       event.WatchDateStr,
			DayVideoSeconds: daySeconds,
			Duplicate:       true,
			ClientEventID:   event.ClientEventID,
		}, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	var daySeconds int
	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		row := models.StudentWatchEvent{
			TenantID:       tenantID,
			StudentID:      studentID,
			CourseID:       event.CourseID,
			LessonID:       event.LessonID,
			Source:         event.Source,
			WatchedSeconds: event.WatchedSeconds,
			ClientEventID:  event.ClientEventID,
			WatchedAt:      event.WatchedAt,
			WatchDate:      event.WatchDate,
			Timezone:       event.Timezone,
			DevicePlatform: event.DevicePlatform,
		}
		if createErr := tx.Create(&row).Error; createErr != nil {
			if isDuplicateKey(createErr) {
				return gorm.ErrDuplicatedKey
			}
			return createErr
		}

		daily := models.StudentDailyWatch{
			TenantID:     tenantID,
			StudentID:    studentID,
			WatchDate:    event.WatchDate,
			Timezone:     event.Timezone,
			VideoSeconds: event.WatchedSeconds,
		}
		if upsertErr := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "tenant_id"}, {Name: "student_id"}, {Name: "watch_date"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"video_seconds": gorm.Expr("video_seconds + ?", event.WatchedSeconds),
				"timezone":      event.Timezone,
				"updated_at":    time.Now(),
			}),
		}).Create(&daily).Error; upsertErr != nil {
			return upsertErr
		}

		var updated models.StudentDailyWatch
		if findErr := tx.Where("tenant_id = ? AND student_id = ? AND watch_date = ?", tenantID, studentID, event.WatchDate.Format("2006-01-02")).
			First(&updated).Error; findErr != nil {
			return findErr
		}
		daySeconds = updated.VideoSeconds
		return nil
	})

	if errors.Is(txErr, gorm.ErrDuplicatedKey) || isDuplicateKey(txErr) {
		daySeconds, dayErr := s.dayVideoSeconds(tenantID, studentID, event.WatchDate)
		if dayErr != nil {
			return nil, dayErr
		}
		return &WatchTimeAcceptResponse{
			Accepted:        true,
			WatchDate:       event.WatchDateStr,
			DayVideoSeconds: daySeconds,
			Duplicate:       true,
			ClientEventID:   event.ClientEventID,
		}, nil
	}
	if txErr != nil {
		return nil, txErr
	}

	return &WatchTimeAcceptResponse{
		Accepted:        true,
		WatchDate:       event.WatchDateStr,
		DayVideoSeconds: daySeconds,
		Duplicate:       false,
		ClientEventID:   event.ClientEventID,
	}, nil
}

func (s *StorefrontService) dayVideoSeconds(tenantID, studentID uint, watchDate time.Time) (int, error) {
	var row models.StudentDailyWatch
	err := s.db.Where("tenant_id = ? AND student_id = ? AND watch_date = ?", tenantID, studentID, watchDate.Format("2006-01-02")).
		First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return row.VideoSeconds, nil
}

func (s *StorefrontService) dayTotalsForDates(tenantID, studentID uint, dates map[string]struct{}) ([]DailyWatchSecondsEntry, error) {
	if len(dates) == 0 {
		return []DailyWatchSecondsEntry{}, nil
	}
	keys := make([]string, 0, len(dates))
	for d := range dates {
		keys = append(keys, d)
	}

	var rows []models.StudentDailyWatch
	if err := s.db.Where("tenant_id = ? AND student_id = ? AND watch_date IN ?", tenantID, studentID, keys).
		Order("watch_date ASC").
		Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]DailyWatchSecondsEntry, 0, len(rows))
	for _, row := range rows {
		result = append(result, DailyWatchSecondsEntry{
			Date:    row.WatchDate.Format("2006-01-02"),
			Seconds: int64(row.VideoSeconds),
		})
	}
	return result, nil
}

func (s *StorefrontService) buildLearningReport(tenantID, studentID uint, period string, allowAll bool) (*LearningReportResponse, error) {
	loc := reportLocation()
	today := startOfLocalDay(time.Now(), loc)

	var (
		dailyWatch []DailyWatchSecondsEntry
		startDate  time.Time
		err        error
	)

	if allowAll && period == "all" {
		dailyWatch, startDate, err = s.aggregateAllDailyWatchSeconds(tenantID, studentID, loc)
		if err != nil {
			return nil, err
		}
		period = "all"
	} else {
		days, parseErr := parseReportPeriod(period)
		if parseErr != nil {
			return nil, parseErr
		}
		if period == "" {
			period = "7d"
		}
		startDate = today.AddDate(0, 0, -(days - 1))
		dailyWatch, err = s.aggregateDailyWatchSeconds(tenantID, studentID, startDate, days, loc)
		if err != nil {
			return nil, err
		}
	}

	quizStart := startDate
	if period == "all" {
		quizStart = time.Time{}
	}

	quizAccuracy, err := s.calcQuizAccuracy(tenantID, studentID, quizStart)
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
		StreakDays:          calcStreakDaysEndingToday(dailyWatch, today, loc),
		QuizAccuracyPercent: quizAccuracy,
		CoursesInProgress:   inProgress,
		CoursesCompleted:    completed,
	}, nil
}

func (s *StorefrontService) aggregateAllDailyWatchSeconds(tenantID, studentID uint, loc *time.Location) ([]DailyWatchSecondsEntry, time.Time, error) {
	type row struct {
		Date    string
		Seconds int64
	}
	var rows []row
	err := s.db.Model(&models.StudentDailyWatch{}).
		Select("DATE_FORMAT(watch_date, '%Y-%m-%d') as date, video_seconds as seconds").
		Where("tenant_id = ? AND student_id = ?", tenantID, studentID).
		Order("watch_date ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, time.Time{}, err
	}

	result := make([]DailyWatchSecondsEntry, 0, len(rows))
	var earliest time.Time
	for i, r := range rows {
		key := normalizeReportDate(r.Date)
		result = append(result, DailyWatchSecondsEntry{Date: key, Seconds: r.Seconds})
		if i == 0 {
			if parsed, parseErr := time.ParseInLocation("2006-01-02", key, loc); parseErr == nil {
				earliest = parsed
			}
		}
	}
	if earliest.IsZero() {
		earliest = startOfLocalDay(time.Now(), loc)
	}
	return result, earliest, nil
}

func calcStreakDaysEndingToday(daily []DailyWatchSecondsEntry, today time.Time, loc *time.Location) int {
	if len(daily) == 0 {
		return 0
	}
	byDate := make(map[string]int64, len(daily))
	for _, d := range daily {
		byDate[d.Date] = d.Seconds
	}

	streak := 0
	for i := 0; ; i++ {
		day := today.In(loc).AddDate(0, 0, -i).Format("2006-01-02")
		if byDate[day] > 0 {
			streak++
			continue
		}
		break
	}
	return streak
}

func (s *StorefrontService) calcLearningTotals(tenantID, studentID uint, periodDays []DailyWatchSecondsEntry) (*LearningReportTotals, error) {
	var periodSum int64
	for _, d := range periodDays {
		periodSum += d.Seconds
	}

	var allTime int64
	if err := s.db.Model(&models.StudentDailyWatch{}).
		Select("COALESCE(SUM(video_seconds), 0)").
		Where("tenant_id = ? AND student_id = ?", tenantID, studentID).
		Scan(&allTime).Error; err != nil {
		return nil, err
	}

	lastAt, err := s.lastWatchedAt(tenantID, studentID)
	if err != nil {
		return nil, err
	}

	return &LearningReportTotals{
		VideoSecondsPeriod:  periodSum,
		VideoSecondsAllTime: allTime,
		LastWatchedAt:       lastAt,
	}, nil
}

func (s *StorefrontService) lastWatchedAt(tenantID, studentID uint) (*string, error) {
	var event models.StudentWatchEvent
	err := s.db.Where("tenant_id = ? AND student_id = ?", tenantID, studentID).
		Order("watched_at DESC").
		First(&event).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	formatted := event.WatchedAt.UTC().Format(time.RFC3339)
	return &formatted, nil
}

func (s *StorefrontService) sumVideoSecondsSince(tenantID, studentID uint, startDate time.Time) (int64, error) {
	var total int64
	err := s.db.Model(&models.StudentDailyWatch{}).
		Select("COALESCE(SUM(video_seconds), 0)").
		Where("tenant_id = ? AND student_id = ? AND watch_date >= ?", tenantID, studentID, startDate.Format("2006-01-02")).
		Scan(&total).Error
	return total, err
}

func (s *StorefrontService) calcWatchByCourse(tenantID, studentID uint, period string) ([]LearningReportByCourse, error) {
	loc := reportLocation()
	today := startOfLocalDay(time.Now(), loc)

	query := s.db.Table("student_watch_events AS e").
		Select("e.course_id AS course_id, COALESCE(c.title, '') AS course_title, COALESCE(SUM(e.watched_seconds), 0) AS video_seconds").
		Joins("LEFT JOIN course_details AS c ON c.id = e.course_id").
		Where("e.tenant_id = ? AND e.student_id = ? AND e.course_id IS NOT NULL", tenantID, studentID).
		Group("e.course_id, c.title").
		Order("video_seconds DESC")

	if period != "all" {
		days, err := parseReportPeriod(period)
		if err != nil {
			return nil, err
		}
		startDate := today.AddDate(0, 0, -(days - 1))
		query = query.Where("e.watch_date >= ?", startDate.Format("2006-01-02"))
	}

	var rows []LearningReportByCourse
	if err := query.Scan(&rows).Error; err != nil {
		return nil, err
	}
	if rows == nil {
		rows = []LearningReportByCourse{}
	}
	return rows, nil
}

func isDuplicateKey(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "duplicate") || strings.Contains(msg, "unique constraint")
}
