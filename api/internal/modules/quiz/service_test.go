package quiz

import (
	"fmt"
	"testing"
	"time"

	"dashlearn/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupQuizServiceTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	stmts := []string{
		`CREATE TABLE IF NOT EXISTS course_details (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			slug TEXT NOT NULL,
			summary TEXT NOT NULL,
			description TEXT,
			visibility TEXT NOT NULL DEFAULT 'public',
			is_scheduled INTEGER DEFAULT 0,
			schedule_date TEXT,
			schedule_time TEXT,
			show_comming_soon INTEGER DEFAULT 0,
			featured_image TEXT,
			intro_video TEXT,
			pricing_model TEXT NOT NULL DEFAULT 'free',
			regular_price REAL DEFAULT 0,
			sale_price REAL DEFAULT 0,
			tags TEXT,
			overview TEXT,
			author_id INTEGER NOT NULL,
			position INTEGER DEFAULT 0,
			tenant_id INTEGER NOT NULL,
			created_at DATETIME,
			updated_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS course_chapters (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			position INTEGER DEFAULT 0,
			title TEXT NOT NULL,
			description TEXT,
			access TEXT NOT NULL DEFAULT 'published',
			course_id INTEGER NOT NULL,
			created_at DATETIME,
			updated_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS course_quizzes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			course_id INTEGER NOT NULL,
			chapter_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			instructions TEXT NOT NULL,
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
		`CREATE TABLE IF NOT EXISTS quiz_questions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			quiz_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			details TEXT,
			media TEXT,
			options TEXT,
			correct_answer TEXT,
			answer_explanation TEXT,
			type TEXT NOT NULL,
			marks REAL DEFAULT 1,
			answer_required INTEGER DEFAULT 0,
			created_at DATETIME,
			updated_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS students (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id TEXT,
			tenant_id INTEGER NOT NULL,
			first_name TEXT NOT NULL,
			last_name TEXT,
			phone TEXT,
			email TEXT NOT NULL,
			password TEXT NOT NULL,
			profile_image TEXT,
			status INTEGER DEFAULT 0,
			otp_code TEXT,
			created_at DATETIME,
			updated_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS enrollments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			tenant_id INTEGER NOT NULL,
			student_id INTEGER NOT NULL,
			course_id INTEGER NOT NULL,
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
			status TEXT DEFAULT 'submitted',
			submitted_at DATETIME,
			graded_at DATETIME,
			instructor_feedback TEXT,
			created_at DATETIME,
			updated_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS quiz_submission_answers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			submission_id INTEGER NOT NULL,
			question_id INTEGER NOT NULL,
			answer TEXT,
			is_correct INTEGER,
			marks_awarded REAL DEFAULT 0,
			created_at DATETIME,
			updated_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS quiz_attempt_sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			tenant_id INTEGER NOT NULL,
			student_id INTEGER NOT NULL,
			quiz_id INTEGER NOT NULL,
			attempt_number INTEGER NOT NULL,
			question_order TEXT NOT NULL,
			started_at DATETIME NOT NULL,
			expires_at DATETIME,
			submitted_at DATETIME,
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

func seedQuizFixture(t *testing.T, db *gorm.DB) (tenantID, courseID, chapterID, quizID, studentID, questionID uint) {
	t.Helper()

	tenantID = 1
	course := models.CourseDetails{
		Title:        "Test Course",
		Slug:         "test-course",
		Summary:      "summary",
		Visibility:   models.Public,
		PricingModel: models.CoursePricingModelFree,
		AuthorID:     1,
		TenantID:     tenantID,
	}
	if err := db.Create(&course).Error; err != nil {
		t.Fatalf("create course: %v", err)
	}
	courseID = course.ID

	chapter := models.CourseChapter{
		Title:    "Chapter 1",
		Access:   models.Published,
		CourseID: course.ID,
	}
	if err := db.Create(&chapter).Error; err != nil {
		t.Fatalf("create chapter: %v", err)
	}
	chapterID = chapter.ID

	quiz := models.CourseQuiz{
		CourseID:              course.ID,
		ChapterID:             chapter.ID,
		Title:                 "Quiz 1",
		Instructions:          "",
		IsPublished:           true,
		TimeLimit:             0,
		TimeLimitOption:       models.CourseQuizTimeLimitOptionMinute,
		RevealAnswers:         true,
		EnableRetry:           true,
		RetryAttempts:         2,
		MinimumPassPercentage: 60,
	}
	if err := db.Create(&quiz).Error; err != nil {
		t.Fatalf("create quiz: %v", err)
	}
	quizID = quiz.ID

	question := models.QuizQuestion{
		QuizID: quiz.ID,
		Title:  "Q1",
		Type:   models.QuizQuestionTypeSingleChoice,
		Marks:  1,
	}
	if err := db.Create(&question).Error; err != nil {
		t.Fatalf("create question: %v", err)
	}
	questionID = question.ID

	student := models.Student{
		TenantID:  tenantID,
		FirstName: "Test",
		Email:     "student@example.com",
		Password:  "secret",
	}
	if err := db.Create(&student).Error; err != nil {
		t.Fatalf("create student: %v", err)
	}
	studentID = student.ID

	enrollment := models.Enrollment{
		TenantID:  tenantID,
		StudentID: student.ID,
		CourseID:  course.ID,
	}
	if err := db.Create(&enrollment).Error; err != nil {
		t.Fatalf("create enrollment: %v", err)
	}

	return tenantID, courseID, chapterID, quizID, studentID, questionID
}

func TestSkipQuizWithoutActiveSessionCreatesForfeitSubmission(t *testing.T) {
	db := setupQuizServiceTestDB(t)
	tenantID, _, _, quizID, studentID, _ := seedQuizFixture(t, db)
	svc := NewQuizService(db)

	detail, err := svc.SkipQuiz(tenantID, studentID, "test-course", quizID)
	if err != nil {
		t.Fatalf("SkipQuiz: %v", err)
	}
	if detail.Score != 0 || detail.Passed {
		t.Fatalf("expected forfeit score 0 and failed, got score=%v passed=%v", detail.Score, detail.Passed)
	}
	if detail.Status != models.QuizSubmissionStatusGraded {
		t.Fatalf("expected graded status, got %s", detail.Status)
	}

	var count int64
	db.Model(&models.QuizSubmission{}).
		Where("tenant_id = ? AND student_id = ? AND quiz_id = ?", tenantID, studentID, quizID).
		Count(&count)
	if count != 1 {
		t.Fatalf("expected 1 submission, got %d", count)
	}
}

func TestSkipQuizWithActiveSessionForfeitsSession(t *testing.T) {
	db := setupQuizServiceTestDB(t)
	tenantID, _, _, quizID, studentID, _ := seedQuizFixture(t, db)
	svc := NewQuizService(db)

	_, err := svc.GetStudentQuiz(tenantID, studentID, "test-course", quizID)
	if err != nil {
		t.Fatalf("GetStudentQuiz: %v", err)
	}

	detail, err := svc.SkipQuiz(tenantID, studentID, "test-course", quizID)
	if err != nil {
		t.Fatalf("SkipQuiz: %v", err)
	}
	if detail.AttemptNumber != 1 {
		t.Fatalf("expected attempt 1, got %d", detail.AttemptNumber)
	}

	var session models.QuizAttemptSession
	if err := db.Where("student_id = ? AND quiz_id = ?", studentID, quizID).First(&session).Error; err != nil {
		t.Fatalf("load session: %v", err)
	}
	if session.SubmittedAt == nil {
		t.Fatal("expected session submitted_at to be set")
	}
}

func TestGetStudentSubmissionReturnsDetailForOwner(t *testing.T) {
	db := setupQuizServiceTestDB(t)
	tenantID, _, _, quizID, studentID, questionID := seedQuizFixture(t, db)
	svc := NewQuizService(db)

	submitted, err := svc.SubmitQuiz(tenantID, studentID, "test-course", quizID, SubmitQuizInput{
		Answers: []SubmitQuizAnswerInput{{QuestionID: questionID, Value: "a"}},
	})
	if err != nil {
		t.Fatalf("SubmitQuiz: %v", err)
	}

	detail, err := svc.GetStudentSubmission(tenantID, studentID, submitted.ID)
	if err != nil {
		t.Fatalf("GetStudentSubmission: %v", err)
	}
	if detail.ID != submitted.ID {
		t.Fatalf("expected submission id %d, got %d", submitted.ID, detail.ID)
	}
	if len(detail.Answers) == 0 {
		t.Fatal("expected answers in student submission detail")
	}
}

func TestGetStudentSubmissionNotFoundForOtherStudent(t *testing.T) {
	db := setupQuizServiceTestDB(t)
	tenantID, _, _, quizID, studentID, questionID := seedQuizFixture(t, db)
	svc := NewQuizService(db)

	submitted, err := svc.SubmitQuiz(tenantID, studentID, "test-course", quizID, SubmitQuizInput{
		Answers: []SubmitQuizAnswerInput{{QuestionID: questionID, Value: "a"}},
	})
	if err != nil {
		t.Fatalf("SubmitQuiz: %v", err)
	}

	other := models.Student{
		TenantID:  tenantID,
		FirstName: "Other",
		Email:     "other@example.com",
		Password:  "secret",
	}
	if err := db.Create(&other).Error; err != nil {
		t.Fatalf("create other student: %v", err)
	}

	if _, err := svc.GetStudentSubmission(tenantID, other.ID, submitted.ID); err == nil {
		t.Fatal("expected not found for other student")
	}
}

func TestUpdateSubmissionFeedback(t *testing.T) {
	db := setupQuizServiceTestDB(t)
	tenantID, courseID, _, quizID, studentID, questionID := seedQuizFixture(t, db)
	svc := NewQuizService(db)

	submitted, err := svc.SubmitQuiz(tenantID, studentID, "test-course", quizID, SubmitQuizInput{
		Answers: []SubmitQuizAnswerInput{{QuestionID: questionID, Value: "a"}},
	})
	if err != nil {
		t.Fatalf("SubmitQuiz: %v", err)
	}

	feedback := "<p>Well done!</p>"
	updated, err := svc.UpdateSubmissionFeedback(tenantID, courseID, submitted.ID, UpdateQuizSubmissionFeedbackInput{
		Feedback: &feedback,
	})
	if err != nil {
		t.Fatalf("UpdateSubmissionFeedback: %v", err)
	}
	if updated.InstructorFeedback == nil || *updated.InstructorFeedback != feedback {
		t.Fatalf("expected feedback %q, got %#v", feedback, updated.InstructorFeedback)
	}
}

func TestSkipQuizRespectsMaximumAttempts(t *testing.T) {
	db := setupQuizServiceTestDB(t)
	tenantID, _, _, quizID, studentID, _ := seedQuizFixture(t, db)

	var quiz models.CourseQuiz
	if err := db.First(&quiz, quizID).Error; err != nil {
		t.Fatalf("load quiz: %v", err)
	}
	quiz.EnableRetry = true
	quiz.RetryAttempts = 1
	if err := db.Save(&quiz).Error; err != nil {
		t.Fatalf("update quiz: %v", err)
	}

	now := time.Now()
	submission := models.QuizSubmission{
		TenantID:      tenantID,
		CourseID:      quiz.CourseID,
		ChapterID:     quiz.ChapterID,
		QuizID:        quiz.ID,
		StudentID:     studentID,
		AttemptNumber: 1,
		Status:        models.QuizSubmissionStatusGraded,
		SubmittedAt:   now,
		GradedAt:      &now,
	}
	if err := db.Create(&submission).Error; err != nil {
		t.Fatalf("seed submission: %v", err)
	}

	svc := NewQuizService(db)
	if _, err := svc.SkipQuiz(tenantID, studentID, "test-course", quizID); err == nil {
		t.Fatal("expected maximum attempts error")
	}
}
