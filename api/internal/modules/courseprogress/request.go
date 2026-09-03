package courseprogress

type UpdateLessonVideoProgressRequest struct {
	// Use pointer so 0 is valid and omitted optional fields don't fail validation.
	MaxPositionSeconds float64  `json:"max_position_seconds" binding:"required,min=0"`
	DurationSeconds    float64  `json:"duration_seconds" binding:"omitempty,min=0"`
	ProgressPercent    *float64 `json:"progress_percent" binding:"omitempty,min=0,max=100"`
}
