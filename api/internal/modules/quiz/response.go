package quiz

import (
	"dashlearn/internal/models"
	"dashlearn/internal/response"
)

type SubmitQuizInput struct {
	Answers []SubmitQuizAnswerInput `json:"answers" binding:"required,min=1,dive"`
}

type SubmitQuizAnswerInput struct {
	QuestionID uint        `json:"question_id" binding:"required"`
	Value      interface{} `json:"value" binding:"required"`
}

type UpdateQuizSubmissionFeedbackInput struct {
	Feedback *string `json:"feedback" binding:"omitempty"`
}

type QuizSubmissionListItem struct {
	ID            uint                        `json:"id"`
	QuizID        uint                        `json:"quiz_id"`
	QuizTitle     string                      `json:"quiz_title"`
	ChapterID     uint                        `json:"chapter_id"`
	ChapterTitle  string                      `json:"chapter_title"`
	StudentID     uint                        `json:"student_id"`
	StudentName   string                      `json:"student_name"`
	StudentEmail  string                      `json:"student_email"`
	AttemptNumber int                         `json:"attempt_number"`
	Score         float32                     `json:"score"`
	MaxScore      float32                     `json:"max_score"`
	Percentage    float32                     `json:"percentage"`
	Passed        bool                        `json:"passed"`
	Status        models.QuizSubmissionStatus `json:"status"`
	SubmittedAt   string                      `json:"submitted_at"`
}

type QuizSubmissionDetail struct {
	QuizSubmissionListItem
	RevealAnswers           bool                           `json:"reveal_answers"`
	TotalQuestions          int                            `json:"total_questions"`
	CorrectCount            int                            `json:"correct_count"`
	IncorrectCount          int                            `json:"incorrect_count"`
	UnansweredCount         int                            `json:"unanswered_count"`
	PassMarks               float32                        `json:"pass_marks"`
	MinimumPassPercentage   float32                        `json:"minimum_pass_percentage"`
	QuizTimeSeconds         int                            `json:"quiz_time_seconds"`
	AttemptTimeSeconds      *int                           `json:"attempt_time_seconds,omitempty"`
	AttemptStartedAt        string                         `json:"attempt_started_at,omitempty"`
	InstructorFeedback      *string                        `json:"instructor_feedback,omitempty"`
	Answers                 []QuizSubmissionAnswerResponse `json:"answers"`
}

type QuizSubmissionAnswerResponse struct {
	QuestionID        uint                    `json:"question_id"`
	QuestionTitle     string                  `json:"question_title"`
	QuestionType      models.QuizQuestionType `json:"question_type"`
	Options           interface{}             `json:"options,omitempty"`
	SubmittedAnswer   interface{}             `json:"submitted_answer"`
	IsCorrect         *bool                   `json:"is_correct"`
	MarksAwarded      float32                 `json:"marks_awarded"`
	QuestionMarks     float32                 `json:"question_marks"`
	AnswerExplanation *string                 `json:"answer_explanation,omitempty"`
	CorrectAnswer     interface{}             `json:"correct_answer,omitempty"`
}

type StudentQuizResponse struct {
	response.CourseQuizResponse
	AttemptsUsed        int    `json:"attempts_used"`
	CanRetry            bool   `json:"can_retry"`
	DisplayMode         string `json:"display_mode"`
	AttemptNumber       int    `json:"attempt_number"`
	TotalQuestions      int    `json:"total_questions"`
	StartedAt           string `json:"started_at,omitempty"`
	ExpiresAt           string `json:"expires_at,omitempty"`
	SecondsRemaining    *int   `json:"seconds_remaining,omitempty"`
	CurrentQuestionIndex *int  `json:"current_question_index,omitempty"`
}

type StudentQuizQuestionResponse struct {
	response.CourseQuizQuestionsResponse
	AttemptNumber        int    `json:"attempt_number"`
	QuestionIndex        int    `json:"question_index"`
	TotalQuestions       int    `json:"total_questions"`
	DisplayMode          string `json:"display_mode"`
	StartedAt            string `json:"started_at,omitempty"`
	ExpiresAt            string `json:"expires_at,omitempty"`
	SecondsRemaining     *int   `json:"seconds_remaining,omitempty"`
}
