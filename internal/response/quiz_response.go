package response

import (
	"dashlearn/models"
	"time"

	"gorm.io/datatypes"
)

type CourseQuizResponse struct {
	ID                    uint                             `gorm:"primaryKey;autoIncrement" json:"id"`
	CourseID              uint                             `json:"course_id" gorm:"column:course_id"`
	ChapterID             uint                             `json:"chapter_id" gorm:"column:chapter_id"`
	Title                 string                           `json:"title" gorm:"type:varchar(255)"`
	Instructions          string                           `json:"instructions" gorm:"type:text"`
	IsPublished           bool                             `json:"is_published" gorm:"column:is_published;default:false"`
	RandomizeQuestions    bool                             `json:"randomize_questions" gorm:"column:randomize_questions;default:false"`
	SingleQuizView        bool                             `json:"single_quiz_view" gorm:"column:single_quiz_view;default:false"`
	TimeLimit             int                              `json:"time_limit" gorm:"column:time_limit;default:1"`
	TimeLimitOption       models.CourseQuizTimeLimitOption `json:"time_limit_option" gorm:"column:time_limit_option;type:enum('minutes','hours','days','weeks','months');default:'weeks'"`
	TotalVisibleQuestions *int                             `json:"total_visible_questions" gorm:"column:total_visible_questions;default:0"`
	RevealAnswers         bool                             `json:"reveal_answers" gorm:"column:reveal_answers;default:false"`
	EnableRetry           bool                             `json:"enable_retry" gorm:"column:enable_retry;default:false"`
	RetryAttempts         int                              `json:"retry_attempts" gorm:"column:retry_attempts;default:1"`
	MinimumPassPercentage float32                          `json:"minimum_pass_percentage" gorm:"column:minimum_pass_percentage;default:0"`
	CreatedAt             time.Time                        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt             time.Time                        `gorm:"autoUpdateTime" json:"updated_at"`
	Questions             []CourseQuizQuestionsResponse    `gorm:"foreignKey:QuizID;references:ID" json:"questions"`
}

func (CourseQuizResponse) TableName() string {
	return "course_quizzes"
}

type CourseQuizQuestionsResponse struct {
	ID                uint                    `gorm:"primaryKey;autoIncrement" json:"id"`
	QuizID            uint                    `gorm:"column:quiz_id" json:"quiz_id"`
	Title             string                  `json:"title" gorm:"type:varchar(255)"`
	Details           *string                 `json:"details" gorm:"type:text"`
	Media             *datatypes.JSON         `gorm:"type:json" json:"media"`
	Type              models.QuizQuestionType `gorm:"type:enum('multiple_choice','single_choice','true_false')" json:"type"`
	Marks             float32                 `json:"marks" gorm:"column:marks;default:1"`
	AnswerRequired    bool                    `json:"answer_required" gorm:"column:answer_required;default:false"`
	AnswerExplanation *string                 `json:"answer_explanation" gorm:"type:text"`
	CreatedAt         time.Time               `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time               `gorm:"autoUpdateTime" json:"updated_at"`
}

func (CourseQuizQuestionsResponse) TableName() string {
	return "quiz_questions"
}
