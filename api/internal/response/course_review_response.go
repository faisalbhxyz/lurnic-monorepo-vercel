package response

import "time"

type CourseReviewStudent struct {
	ID           uint    `json:"id"`
	FirstName    string  `json:"first_name"`
	LastName     *string `json:"last_name,omitempty"`
	ProfileImage *string `json:"profile_image,omitempty"`
}

type CourseReviewResponse struct {
	ID        uint                 `json:"id"`
	CourseID  uint                 `json:"course_id"`
	StudentID uint                 `json:"student_id"`
	Rating    uint8                `json:"rating"`
	Comment   *string              `json:"comment,omitempty"`
	Tags      []string             `json:"tags,omitempty"`
	Student   *CourseReviewStudent `json:"student,omitempty"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
}

type CourseReviewsSummary struct {
	AverageRating float64               `json:"average_rating"`
	TotalReviews  int                   `json:"total_reviews"`
	Reviews       []CourseReviewResponse `json:"reviews"`
	StudentReview *CourseReviewResponse `json:"student_review,omitempty"`
	CanReview     *bool                 `json:"can_review,omitempty"`
}
