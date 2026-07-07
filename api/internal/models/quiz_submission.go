package models

import "time"

type QuizSubmissionStatus string

const (
	QuizSubmissionStatusSubmitted     QuizSubmissionStatus = "submitted"
	QuizSubmissionStatusGraded        QuizSubmissionStatus = "graded"
	QuizSubmissionStatusPendingReview QuizSubmissionStatus = "pending_review"
)

type QuizSubmission struct {
	ID            uint                 `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID      uint                 `gorm:"column:tenant_id" json:"-"`
	CourseID      uint                 `gorm:"column:course_id" json:"course_id"`
	ChapterID     uint                 `gorm:"column:chapter_id" json:"chapter_id"`
	QuizID        uint                 `gorm:"column:quiz_id" json:"quiz_id"`
	StudentID     uint                 `gorm:"column:student_id" json:"student_id"`
	AttemptNumber int                  `gorm:"column:attempt_number" json:"attempt_number"`
	Score         float32              `json:"score"`
	MaxScore      float32              `gorm:"column:max_score" json:"max_score"`
	Percentage    float32              `json:"percentage"`
	Passed        bool                 `json:"passed"`
	Status              QuizSubmissionStatus `gorm:"type:enum('submitted','graded','pending_review');default:'submitted'" json:"status"`
	SubmittedAt         time.Time            `gorm:"column:submitted_at" json:"submitted_at"`
	GradedAt            *time.Time           `gorm:"column:graded_at" json:"graded_at,omitempty"`
	InstructorFeedback  *string              `gorm:"column:instructor_feedback;type:text" json:"instructor_feedback,omitempty"`
	CreatedAt     time.Time            `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time            `gorm:"autoUpdateTime" json:"updated_at"`
	Quiz          CourseQuiz           `gorm:"foreignKey:QuizID;references:ID" json:"quiz,omitempty"`
	Student       Student              `gorm:"foreignKey:StudentID;references:ID" json:"student,omitempty"`
	Answers       []QuizSubmissionAnswer `gorm:"foreignKey:SubmissionID;references:ID" json:"answers,omitempty"`
}

func (QuizSubmission) TableName() string {
	return "quiz_submissions"
}

type QuizSubmissionAnswer struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	SubmissionID uint      `gorm:"column:submission_id" json:"submission_id"`
	QuestionID   uint      `gorm:"column:question_id" json:"question_id"`
	Answer       []byte    `gorm:"type:json" json:"answer"`
	IsCorrect    *bool     `gorm:"column:is_correct" json:"is_correct"`
	MarksAwarded float32   `gorm:"column:marks_awarded" json:"marks_awarded"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Question     QuizQuestion `gorm:"foreignKey:QuestionID;references:ID" json:"question,omitempty"`
}

func (QuizSubmissionAnswer) TableName() string {
	return "quiz_submission_answers"
}
