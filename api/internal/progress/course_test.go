package progress

import (
	"fmt"
	"testing"
	"time"

	"dashlearn/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupProgressTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	stmts := []string{
		`CREATE TABLE IF NOT EXISTS course_lessons (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			chapter_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			is_published INTEGER DEFAULT 0,
			created_at DATETIME,
			updated_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS course_chapters (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			position INTEGER DEFAULT 0,
			course_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			description TEXT,
			access TEXT NOT NULL DEFAULT 'published',
			created_at DATETIME,
			updated_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS course_quizzes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			course_id INTEGER NOT NULL,
			chapter_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			instructions TEXT NOT NULL DEFAULT '',
			is_published INTEGER DEFAULT 0,
			randomize_questions INTEGER DEFAULT 0,
			single_quiz_view INTEGER DEFAULT 0,
			time_limit INTEGER DEFAULT 1,
			time_limit_option TEXT DEFAULT 'weeks',
			total_visible_questions INTEGER DEFAULT 0,
			reveal_answers INTEGER DEFAULT 0,
			enable_retry INTEGER DEFAULT 0,
			retry_attempts INTEGER DEFAULT 0,
			minimum_pass_percentage REAL DEFAULT 0,
			position INTEGER DEFAULT 0,
			created_at DATETIME,
			updated_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS student_lesson_completions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			tenant_id INTEGER NOT NULL,
			student_id INTEGER NOT NULL,
			course_id INTEGER NOT NULL,
			lesson_id INTEGER NOT NULL,
			completed_at DATETIME,
			created_at DATETIME,
			updated_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS quiz_submissions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			tenant_id INTEGER NOT NULL,
			course_id INTEGER NOT NULL,
			chapter_id INTEGER NOT NULL,
			quiz_id INTEGER NOT NULL,
			student_id INTEGER NOT NULL,
			attempt_number INTEGER NOT NULL,
			score REAL DEFAULT 0,
			max_score REAL DEFAULT 0,
			percentage REAL DEFAULT 0,
			passed INTEGER DEFAULT 0,
			status TEXT DEFAULT 'graded',
			submitted_at DATETIME,
			graded_at DATETIME,
			instructor_feedback TEXT,
			created_at DATETIME,
			updated_at DATETIME
		)`,
	}
	for _, stmt := range stmts {
		if err := db.Exec(stmt).Error; err != nil {
			t.Fatalf("exec schema: %v", err)
		}
	}

	return db
}

func TestCalcBreakdownIncludesCompletedQuizIDs(t *testing.T) {
	db := setupProgressTestDB(t)

	const tenantID, studentID, courseID, chapterID uint = 1, 2, 10, 20
	chapter := models.CourseChapter{ID: chapterID, CourseID: courseID, Title: "Ch1"}
	if err := db.Create(&chapter).Error; err != nil {
		t.Fatalf("create chapter: %v", err)
	}

	quiz := models.CourseQuiz{
		ID:          101,
		CourseID:    courseID,
		ChapterID:   chapterID,
		Title:       "Quiz A",
		IsPublished: true,
	}
	if err := db.Create(&quiz).Error; err != nil {
		t.Fatalf("create quiz: %v", err)
	}

	now := time.Now()
	submission := models.QuizSubmission{
		TenantID:      tenantID,
		CourseID:      courseID,
		ChapterID:     chapterID,
		QuizID:        quiz.ID,
		StudentID:     studentID,
		AttemptNumber: 1,
		Status:        models.QuizSubmissionStatusGraded,
		SubmittedAt:   now,
		GradedAt:      &now,
	}
	if err := db.Create(&submission).Error; err != nil {
		t.Fatalf("create submission: %v", err)
	}

	breakdown := CalcBreakdown(db, tenantID, studentID, courseID, DefaultOptions(), true)
	if breakdown.QuizzesDone != 1 {
		t.Fatalf("expected quizzes_done 1, got %d", breakdown.QuizzesDone)
	}
	if len(breakdown.CompletedQuizIDs) != 1 || breakdown.CompletedQuizIDs[0] != quiz.ID {
		t.Fatalf("expected completed_quiz_ids [%d], got %#v", quiz.ID, breakdown.CompletedQuizIDs)
	}
}
