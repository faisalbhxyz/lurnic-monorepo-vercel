package assignment

import (
	"dashlearn/internal/models"
	"dashlearn/internal/response"
)

type GradeAssignmentInput struct {
	Score    float32 `json:"score" binding:"required,gte=0"`
	Feedback *string `json:"feedback" binding:"omitempty"`
}

type AssignmentSubmissionListItem struct {
	ID               uint                              `json:"id"`
	AssignmentID     uint                              `json:"assignment_id"`
	AssignmentTitle  string                            `json:"assignment_title"`
	ChapterID        uint                              `json:"chapter_id"`
	ChapterTitle     string                            `json:"chapter_title"`
	StudentID        uint                              `json:"student_id"`
	StudentName      string                            `json:"student_name"`
	StudentEmail     string                            `json:"student_email"`
	Score            float32                           `json:"score"`
	MaxScore         float32                           `json:"max_score"`
	Percentage       float32                           `json:"percentage"`
	Passed           bool                              `json:"passed"`
	Status           models.AssignmentSubmissionStatus `json:"status"`
	SubmittedAt      string                            `json:"submitted_at"`
	FileCount        int                               `json:"file_count"`
}

type AssignmentSubmissionFileResponse struct {
	ID       uint   `json:"id"`
	URL      string `json:"url"`
	FileName string `json:"file_name"`
	MimeType string `json:"mime_type"`
	Size     int64  `json:"size"`
}

type AssignmentSubmissionDetail struct {
	AssignmentSubmissionListItem
	ResponseText       *string                            `json:"response_text,omitempty"`
	InstructorFeedback *string                            `json:"instructor_feedback,omitempty"`
	Files              []AssignmentSubmissionFileResponse `json:"files"`
}

type StudentAssignmentResponse struct {
	response.CourseAssignmentResponse
	HasSubmitted bool                         `json:"has_submitted"`
	CanSubmit    bool                         `json:"can_submit"`
	Submission   *AssignmentSubmissionSummary `json:"submission,omitempty"`
}

type AssignmentSubmissionSummary struct {
	ID          uint                              `json:"id"`
	Score       float32                           `json:"score"`
	MaxScore    float32                           `json:"max_score"`
	Percentage  float32                           `json:"percentage"`
	Passed      bool                              `json:"passed"`
	Status      models.AssignmentSubmissionStatus `json:"status"`
	SubmittedAt string                            `json:"submitted_at"`
}
