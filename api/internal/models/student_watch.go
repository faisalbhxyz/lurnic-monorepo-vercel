package models

import "time"

// StudentDailyWatch is the per-calendar-day aggregate of actual video play seconds.
type StudentDailyWatch struct {
	ID                uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID          uint      `gorm:"column:tenant_id" json:"-"`
	StudentID         uint      `gorm:"column:student_id" json:"-"`
	WatchDate         time.Time `gorm:"column:watch_date;type:date" json:"watch_date"`
	Timezone          string    `gorm:"column:timezone;size:64" json:"timezone"`
	VideoSeconds      int       `gorm:"column:video_seconds" json:"video_seconds"`
	LiveClassSeconds  int       `gorm:"column:live_class_seconds" json:"live_class_seconds"`
	QuizSeconds       int       `gorm:"column:quiz_seconds" json:"quiz_seconds"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (StudentDailyWatch) TableName() string {
	return "student_daily_watch"
}

// StudentWatchEvent is an idempotent delta log of actual play time.
type StudentWatchEvent struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID        uint      `gorm:"column:tenant_id" json:"-"`
	StudentID       uint      `gorm:"column:student_id" json:"-"`
	CourseID        *uint     `gorm:"column:course_id" json:"course_id"`
	LessonID        *uint     `gorm:"column:lesson_id" json:"lesson_id"`
	Source          string    `gorm:"column:source;size:32" json:"source"`
	WatchedSeconds  int       `gorm:"column:watched_seconds" json:"watched_seconds"`
	ClientEventID   string    `gorm:"column:client_event_id;size:64" json:"client_event_id"`
	WatchedAt       time.Time `gorm:"column:watched_at" json:"watched_at"`
	WatchDate       time.Time `gorm:"column:watch_date;type:date" json:"watch_date"`
	Timezone        string    `gorm:"column:timezone;size:64" json:"timezone"`
	DevicePlatform  *string   `gorm:"column:device_platform;size:16" json:"device_platform"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (StudentWatchEvent) TableName() string {
	return "student_watch_events"
}
