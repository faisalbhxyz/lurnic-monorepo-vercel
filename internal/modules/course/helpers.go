package course

import (
	"dashlearn/internal/models"
	"fmt"
	"time"

	"gorm.io/gorm"
)

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
		dateParsed, err := time.ParseInLocation("2006-01-02", *course.ScheduleDate, loc)
		if err != nil {
			// Try fallback RFC3339 (if DB actually stored timezone)
			dateParsed, err = time.Parse(time.RFC3339, *course.ScheduleDate)
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
					"visibility":   models.Public,
					"is_scheduled": isScheduled,
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
