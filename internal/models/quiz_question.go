package models

import (
	"time"

	"gorm.io/datatypes"
)

type QuizQuestionType string

const (
	QuizQuestionTypeMultipleChoice QuizQuestionType = "multiple_choice"
	QuizQuestionTypeSingleChoice   QuizQuestionType = "single_choice"
	QuizQuestionTypeTrueFalse      QuizQuestionType = "true_false"
)

type QuizQuestion struct {
	ID                uint             `gorm:"primaryKey;autoIncrement" json:"id"`
	QuizID            uint             `gorm:"column:quiz_id" json:"quiz_id"`
	Quiz              CourseQuiz       `gorm:"foreignKey:QuizID;references:ID" json:"quiz"`
	Title             string           `json:"title" gorm:"type:varchar(255)"`
	Details           *string          `json:"details" gorm:"type:text"`
	Media             *datatypes.JSON  `gorm:"type:json" json:"media"`
	Type              QuizQuestionType `gorm:"type:enum('multiple_choice','single_choice','true_false')" json:"type"`
	Marks             float32          `json:"marks" gorm:"column:marks;default:1"`
	AnswerRequired    bool             `json:"answer_required" gorm:"column:answer_required;default:false"`
	AnswerExplanation *string          `json:"answer_explanation" gorm:"type:text"`
	CreatedAt         time.Time        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time        `gorm:"autoUpdateTime" json:"updated_at"`
}
