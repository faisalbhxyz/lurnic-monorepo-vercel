package course

import (
	"dashlearn/internal/models"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// func NormalizeDate(dateStr string) (string, error) {
// 	// Case 1: ISO format like "2025-11-28T00:00:00Z"
// 	if t, err := time.Parse(time.RFC3339, dateStr); err == nil {
// 		return t.Format("2006-01-02"), nil
// 	}

// 	// Case 2: simple date like "2025-11-28"
// 	if t, err := time.Parse("2006-01-02", dateStr); err == nil {
// 		return t.Format("2006-01-02"), nil
// 	}

// 	return "", errors.New("invalid date format")
// }

func NormalizeDate(input string) (string, error) {
	// Hard clean: remove Z, timezone, milliseconds, \r, \n
	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, "\r", "")
	input = strings.ReplaceAll(input, "\n", "")
	input = strings.ReplaceAll(input, "Z", "")

	// If there's a T, keep only the date part
	if idx := strings.Index(input, "T"); idx > 0 {
		input = input[:idx]
	}

	// Final clean: ensure only YYYY-MM-DD remains
	input = strings.TrimSpace(input)

	// Validate
	t, err := time.Parse("2006-01-02", input)
	if err != nil {
		return "", fmt.Errorf("invalid date: %s", input)
	}

	return t.Format("2006-01-02"), nil
}

// RunCronJobForCourses checks all scheduled courses and makes them public
func CronJobForCoursesSchedule(db *gorm.DB) error {
	// Load Bangladesh timezone (GMT+6)
	loc, err := time.LoadLocation("Asia/Dhaka")
	if err != nil {
		return fmt.Errorf("failed to load timezone Asia/Dhaka: %w", err)
	}

	now := time.Now().In(loc)

	var courses []models.CourseDetails

	// Fetch only scheduled courses
	if err := db.Model(&models.CourseDetails{}).
		Select("id", "visibility", "is_scheduled", "schedule_date", "schedule_time").
		Where("is_scheduled = ?", true).
		Find(&courses).Error; err != nil {
		return err
	}

	for _, course := range courses {
		if course.ScheduleDate == nil || course.ScheduleTime == nil {
			fmt.Println("⚠️ Schedule missing for course ID", course.ID)
			continue
		}

		// Parse schedule_date safely (your DB probably stores date without timezone)
		dateParsed, err := time.ParseInLocation(time.RFC3339, *course.ScheduleDate, loc)
		if err != nil {
			// Try fallback RFC3339 (if DB actually stored timezone)
			dateParsed, err = time.Parse("2006-01-02", *course.ScheduleDate)
			if err != nil {
				fmt.Println("⚠️ Invalid schedule_date for course ID", course.ID, err)
				continue
			}
		}

		// Parse schedule_time (HH:MM:SS) in Bangladesh timezone
		timeParsed, err := time.ParseInLocation("15:04:05", *course.ScheduleTime, loc)
		if err != nil {
			fmt.Println("⚠️ Invalid schedule_time for course ID", course.ID, err)
			continue
		}

		// Combine into a single datetime in BD timezone
		scheduledTime := time.Date(
			dateParsed.Year(), dateParsed.Month(), dateParsed.Day(),
			timeParsed.Hour(), timeParsed.Minute(), timeParsed.Second(),
			0,
			loc,
		)

		// If current time >= scheduled time → make course public
		if !scheduledTime.After(now) {
			isScheduled := false

			err := db.Model(&models.CourseDetails{}).
				Where("id = ?", course.ID).
				Updates(map[string]interface{}{
					"visibility":    models.Public,
					"is_scheduled":  isScheduled,
					"schedule_date": nil,
					"schedule_time": nil,
				}).Error

			if err != nil {
				fmt.Println("🔥 Failed to update course ID", course.ID, err)
			} else {
				fmt.Println("🎉 Course made public:", course.ID, "at", now)
			}
		}
	}

	return nil
}

// RunCronJobForCourses checks all scheduled courses and makes them public
func CronJobForCourseLessonsSchedule(db *gorm.DB) error {
	// Load Bangladesh timezone (GMT+6)
	loc, err := time.LoadLocation("Asia/Dhaka")
	if err != nil {
		return fmt.Errorf("failed to load timezone Asia/Dhaka: %w", err)
	}

	now := time.Now().In(loc)

	var lessons []models.CourseLesson
	if err := db.Model(&models.CourseLesson{}).
		Select("id", "is_published", "is_scheduled", "schedule_date", "schedule_time").
		Where("is_scheduled = ? AND is_published = ?", true, false).
		Find(&lessons).Error; err != nil {
		return err
	}

	for _, lesson := range lessons {
		if lesson.ScheduleDate == nil || lesson.ScheduleTime == nil {
			fmt.Println("⚠️ Schedule missing for lesson ID", lesson.ID)
			continue
		}

		// Parse schedule_date safely (your DB probably stores date without timezone)
		dateParsed, err := time.ParseInLocation(time.RFC3339, *lesson.ScheduleDate, loc)
		if err != nil {
			// Try fallback RFC3339 (if DB actually stored timezone)
			dateParsed, err = time.Parse("2006-01-02", *lesson.ScheduleDate)
			if err != nil {
				fmt.Println("⚠️ Invalid schedule_date for lesson ID", lesson.ID, err)
				continue
			}
		}

		// Parse schedule_time (HH:MM:SS) in Bangladesh timezone
		timeParsed, err := time.ParseInLocation("15:04:05", *lesson.ScheduleTime, loc)
		if err != nil {
			fmt.Println("⚠️ Invalid schedule_time for lesson ID", lesson.ID, err)
			continue
		}

		// Combine into a single datetime in BD timezone
		scheduledTime := time.Date(
			dateParsed.Year(), dateParsed.Month(), dateParsed.Day(),
			timeParsed.Hour(), timeParsed.Minute(), timeParsed.Second(),
			0,
			loc,
		)

		// If current time >= scheduled time → make course public
		if !scheduledTime.After(now) {
			isScheduled := false

			err := db.Model(&models.CourseLesson{}).
				Where("id = ?", lesson.ID).
				Updates(map[string]any{
					"is_published":   true,
					"is_scheduled":   isScheduled,
					"schedule_date":  nil,
					"schedule_time":  nil,
				}).Error

			if err != nil {
				fmt.Println("🔥 Failed to update lesson ID", lesson.ID, err)
			} else {
				fmt.Println("🎉 lesson made public:", lesson.ID, "at", now)
			}
		}
	}

	return nil
}
