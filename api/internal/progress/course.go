package progress

import (
	"dashlearn/internal/models"
	"errors"
	"math"

	"gorm.io/gorm"
)

type Options struct {
	CountLessons     bool
	CountQuizzes     bool
	CountAssignments bool
}

func DefaultOptions() Options {
	return Options{
		CountLessons:     true,
		CountQuizzes:     true,
		CountAssignments: true,
	}
}

func LoadOptions(db *gorm.DB, courseID uint) Options {
	var settings models.CourseCertificateSettings
	if err := db.Where("course_id = ?", courseID).First(&settings).Error; err != nil {
		return DefaultOptions()
	}
	opts := Options{
		CountLessons:     settings.CountLessons,
		CountQuizzes:     settings.CountQuizzes,
		CountAssignments: settings.CountAssignments,
	}
	if !opts.CountLessons && !opts.CountQuizzes && !opts.CountAssignments {
		return DefaultOptions()
	}
	return opts
}

type Breakdown struct {
	LessonsDone          int64   `json:"lessons_completed"`
	LessonsTotal         int64   `json:"lessons_total"`
	QuizzesDone          int64   `json:"quizzes_completed"`
	QuizzesTotal         int64   `json:"quizzes_total"`
	AssignmentsDone      int64   `json:"assignments_completed"`
	AssignmentsTotal     int64   `json:"assignments_total"`
	Percent              float32 `json:"progress_percent"`
	CountLessons         bool    `json:"count_lessons"`
	CountQuizzes         bool    `json:"count_quizzes"`
	CountAssignments     bool    `json:"count_assignments"`
	CompletedLessonIDs   []uint  `json:"completed_lesson_ids,omitempty"`
}

func CalcCourseProgress(db *gorm.DB, tenantID, studentID, courseID uint, opts Options) float32 {
	return CalcBreakdown(db, tenantID, studentID, courseID, opts, false).Percent
}

func CalcBreakdown(db *gorm.DB, tenantID, studentID, courseID uint, opts Options, includeLessonIDs bool) Breakdown {
	if !opts.CountLessons && !opts.CountQuizzes && !opts.CountAssignments {
		opts = DefaultOptions()
	}

	var lessonsDone, lessonsTotal int64
	var quizzesDone, quizzesTotal int64
	var assignmentsDone, assignmentsTotal int64

	if opts.CountLessons {
		db.Model(&models.CourseLesson{}).
			Joins("JOIN course_chapters ON course_chapters.id = course_lessons.chapter_id").
			Where("course_chapters.course_id = ? AND course_lessons.is_published = ?", courseID, true).
			Count(&lessonsTotal)

		db.Model(&models.StudentLessonCompletion{}).
			Where("tenant_id = ? AND student_id = ? AND course_id = ?", tenantID, studentID, courseID).
			Count(&lessonsDone)
	}

	if opts.CountQuizzes {
		db.Model(&models.CourseQuiz{}).
			Where("course_id = ? AND is_published = ?", courseID, true).
			Count(&quizzesTotal)

		db.Model(&models.QuizSubmission{}).
			Where("tenant_id = ? AND student_id = ? AND course_id = ?", tenantID, studentID, courseID).
			Distinct("quiz_id").
			Count(&quizzesDone)
	}

	if opts.CountAssignments {
		db.Model(&models.CourseAssignment{}).
			Where("course_id = ? AND is_published = ?", courseID, true).
			Count(&assignmentsTotal)

		db.Model(&models.AssignmentSubmission{}).
			Where("tenant_id = ? AND student_id = ? AND course_id = ?", tenantID, studentID, courseID).
			Distinct("assignment_id").
			Count(&assignmentsDone)
	}

	total := lessonsTotal + quizzesTotal + assignmentsTotal
	done := lessonsDone + quizzesDone + assignmentsDone

	percent := float32(0)
	if total > 0 {
		percent = float32(math.Round(float64(done)/float64(total)*100*10) / 10)
	}

	breakdown := Breakdown{
		LessonsDone:      lessonsDone,
		LessonsTotal:     lessonsTotal,
		QuizzesDone:      quizzesDone,
		QuizzesTotal:     quizzesTotal,
		AssignmentsDone:  assignmentsDone,
		AssignmentsTotal: assignmentsTotal,
		Percent:          percent,
		CountLessons:     opts.CountLessons,
		CountQuizzes:     opts.CountQuizzes,
		CountAssignments: opts.CountAssignments,
	}

	if includeLessonIDs {
		var rows []models.StudentLessonCompletion
		db.Select("lesson_id").
			Where("tenant_id = ? AND student_id = ? AND course_id = ?", tenantID, studentID, courseID).
			Find(&rows)
		breakdown.CompletedLessonIDs = make([]uint, 0, len(rows))
		for _, row := range rows {
			breakdown.CompletedLessonIDs = append(breakdown.CompletedLessonIDs, row.LessonID)
		}
	}

	return breakdown
}

func ValidateOptions(opts Options) error {
	if !opts.CountLessons && !opts.CountQuizzes && !opts.CountAssignments {
		return errors.New("at least one progress item type must be selected")
	}
	return nil
}
