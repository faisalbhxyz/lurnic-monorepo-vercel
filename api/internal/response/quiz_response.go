package response

import (
	"dashlearn/internal/models"
	"time"

	"gorm.io/datatypes"
)

type CourseQuizResponse struct {
	ID                    uint                             `json:"id"`
	CourseID              uint                             `json:"course_id"`
	ChapterID             uint                             `json:"chapter_id"`
	Title                 string                           `json:"title"`
	Instructions          string                           `json:"instructions"`
	IsPublished           bool                             `json:"is_published"`
	RandomizeQuestions    bool                             `json:"randomize_questions"`
	SingleQuizView        bool                             `json:"single_quiz_view"`
	TimeLimit             int                              `json:"time_limit"`
	TimeLimitOption       models.CourseQuizTimeLimitOption `json:"time_limit_option"`
	TotalVisibleQuestions *int                             `json:"total_visible_questions"`
	RevealAnswers         bool                             `json:"reveal_answers"`
	EnableRetry           bool                             `json:"enable_retry"`
	RetryAttempts         int                              `json:"retry_attempts"`
	MinimumPassPercentage float32                          `json:"minimum_pass_percentage"`
	CreatedAt             time.Time                        `json:"created_at"`
	UpdatedAt             time.Time                        `json:"updated_at"`
	Questions             []CourseQuizQuestionsResponse    `json:"questions"`
}

type CourseQuizQuestionsResponse struct {
	ID                uint                    `json:"id"`
	QuizID            uint                    `json:"quiz_id"`
	Title             string                  `json:"title" gorm:"type:varchar(255)"`
	Details           *string                 `json:"details" gorm:"type:text"`
	Media             *datatypes.JSON         `json:"media"`
	Type              models.QuizQuestionType `json:"type"`
	Marks             float32                 `json:"marks"`
	AnswerRequired    bool                    `json:"answer_required"`
	AnswerExplanation *string                 `json:"answer_explanation"`
	CreatedAt         time.Time               `json:"created_at"`
	UpdatedAt         time.Time               `json:"updated_at"`
}
