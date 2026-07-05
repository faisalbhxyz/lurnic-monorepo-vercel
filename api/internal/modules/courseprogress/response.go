package courseprogress

import "time"

type LessonVideoProgressResponse struct {
	LessonID           uint      `json:"lesson_id"`
	MaxPositionSeconds float64   `json:"max_position_seconds"`
	DurationSeconds    float64   `json:"duration_seconds"`
	ProgressPercent    float64   `json:"progress_percent"`
	Completed          bool      `json:"completed"`
	UpdatedAt          time.Time `json:"updated_at"`
}
