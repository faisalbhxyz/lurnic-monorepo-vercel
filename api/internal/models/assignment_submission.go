package models

import "time"

type AssignmentSubmissionStatus string

const (
	AssignmentSubmissionStatusSubmitted     AssignmentSubmissionStatus = "submitted"
	AssignmentSubmissionStatusGraded        AssignmentSubmissionStatus = "graded"
	AssignmentSubmissionStatusPendingReview AssignmentSubmissionStatus = "pending_review"
)

type AssignmentSubmission struct {
	ID                  uint                       `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID            uint                       `gorm:"column:tenant_id" json:"-"`
	CourseID            uint                       `gorm:"column:course_id" json:"course_id"`
	ChapterID           uint                       `gorm:"column:chapter_id" json:"chapter_id"`
	AssignmentID        uint                       `gorm:"column:assignment_id" json:"assignment_id"`
	StudentID           uint                       `gorm:"column:student_id" json:"student_id"`
	ResponseText        *string                    `gorm:"column:response_text;type:text" json:"response_text,omitempty"`
	Score               float32                    `json:"score"`
	MaxScore            float32                    `gorm:"column:max_score" json:"max_score"`
	Percentage          float32                    `json:"percentage"`
	Passed              bool                       `json:"passed"`
	Status              AssignmentSubmissionStatus `gorm:"type:enum('submitted','graded','pending_review');default:'pending_review'" json:"status"`
	InstructorFeedback  *string                    `gorm:"column:instructor_feedback;type:text" json:"instructor_feedback,omitempty"`
	SubmittedAt         time.Time                  `gorm:"column:submitted_at" json:"submitted_at"`
	GradedAt            *time.Time                 `gorm:"column:graded_at" json:"graded_at,omitempty"`
	CreatedAt           time.Time                  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time                  `gorm:"autoUpdateTime" json:"updated_at"`
	Assignment          CourseAssignment           `gorm:"foreignKey:AssignmentID;references:ID" json:"assignment,omitempty"`
	Student             Student                    `gorm:"foreignKey:StudentID;references:ID" json:"student,omitempty"`
	Files               []AssignmentSubmissionFile `gorm:"foreignKey:SubmissionID;references:ID" json:"files,omitempty"`
}

func (AssignmentSubmission) TableName() string {
	return "assignment_submissions"
}

type AssignmentSubmissionFile struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	SubmissionID uint      `gorm:"column:submission_id" json:"submission_id"`
	URL          string    `gorm:"column:url;type:varchar(512)" json:"url"`
	FileName     string    `gorm:"column:file_name;type:varchar(255)" json:"file_name"`
	MimeType     string    `gorm:"column:mime_type;type:varchar(128)" json:"mime_type"`
	Size         int64     `json:"size"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (AssignmentSubmissionFile) TableName() string {
	return "assignment_submission_files"
}
