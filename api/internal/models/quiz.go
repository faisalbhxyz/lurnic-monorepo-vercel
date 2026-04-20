package models

import "time"

type CourseQuizTimeLimitOption string

const (
	CourseQuizTimeLimitOptionMinute CourseQuizTimeLimitOption = "minutes"
	CourseQuizTimeLimitOptionHour   CourseQuizTimeLimitOption = "hours"
	CourseQuizTimeLimitOptionDay    CourseQuizTimeLimitOption = "days"
	CourseQuizTimeLimitOptionWeek   CourseQuizTimeLimitOption = "weeks"
	CourseQuizTimeLimitOptionMonth  CourseQuizTimeLimitOption = "months"
)

type CourseQuiz struct {
	ID                    uint                      `gorm:"primaryKey;autoIncrement" json:"id"`
	CourseID              uint                      `json:"course_id" gorm:"column:course_id"`
	ChapterID             uint                      `json:"chapter_id" gorm:"column:chapter_id"`
	Title                 string                    `json:"title" gorm:"type:varchar(255)"`
	Instructions          string                    `json:"instructions" gorm:"type:text"`
	IsPublished           bool                      `json:"is_published" gorm:"column:is_published;default:false"`
	RandomizeQuestions    bool                      `json:"randomize_questions" gorm:"column:randomize_questions;default:false"`
	SingleQuizView        bool                      `json:"single_quiz_view" gorm:"column:single_quiz_view;default:false"`
	TimeLimit             int                       `json:"time_limit" gorm:"column:time_limit;default:1"`
	TimeLimitOption       CourseQuizTimeLimitOption `json:"time_limit_option" gorm:"column:time_limit_option;type:enum('minutes','hours','days','weeks','months');default:'weeks'"`
	TotalVisibleQuestions *int                      `json:"total_visible_questions" gorm:"column:total_visible_questions;default:0"`
	RevealAnswers         bool                      `json:"reveal_answers" gorm:"column:reveal_answers;default:false"`
	EnableRetry           bool                      `json:"enable_retry" gorm:"column:enable_retry;default:false"`
	RetryAttempts         int                       `json:"retry_attempts" gorm:"column:retry_attempts;default:1"`
	MinimumPassPercentage float32                   `json:"minimum_pass_percentage" gorm:"column:minimum_pass_percentage;default:0"`
	CreatedAt             time.Time                 `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt             time.Time                 `gorm:"autoUpdateTime" json:"updated_at"`
	Questions             []QuizQuestion            `gorm:"foreignKey:QuizID;references:ID" json:"questions"`
}

func (CourseQuiz) TableName() string {
	return "course_quizzes"
}
