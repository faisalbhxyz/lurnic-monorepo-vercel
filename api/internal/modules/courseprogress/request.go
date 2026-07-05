package courseprogress

type UpdateLessonVideoProgressRequest struct {
	MaxPositionSeconds float64 `json:"max_position_seconds" binding:"required,min=0"`
	DurationSeconds    float64 `json:"duration_seconds" binding:"min=0"`
	ProgressPercent    float64 `json:"progress_percent" binding:"min=0,max=100"`
}
