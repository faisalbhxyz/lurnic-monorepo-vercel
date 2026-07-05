package course

import (
	"testing"

	"dashlearn/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupQuizTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	stmts := []string{
		`CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			slug TEXT NOT NULL,
			description TEXT,
			thumbnail TEXT,
			tenant_id INTEGER DEFAULT 0,
			created_at DATETIME,
			updated_at DATETIME
		)`,
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
			regular_price REAL,
			sale_price REAL,
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
			retry_attempts INTEGER DEFAULT 1,
			minimum_pass_percentage REAL NOT NULL,
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
			type TEXT NOT NULL,
			marks REAL DEFAULT 1,
			answer_required INTEGER DEFAULT 0,
			answer_explanation TEXT,
			created_at DATETIME,
			updated_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS course_general_settings (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			course_id INTEGER NOT NULL,
			difficulty_level TEXT DEFAULT 'all',
			maximum_student INTEGER DEFAULT 0,
			language TEXT DEFAULT 'english',
			category_id INTEGER NOT NULL,
			sub_category_id INTEGER,
			duration TEXT,
			created_at DATETIME,
			updated_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS course_instructors (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			course_id INTEGER NOT NULL,
			instructor_id INTEGER NOT NULL,
			created_at DATETIME,
			updated_at DATETIME
		)`,
	}
	for _, stmt := range stmts {
		if err := db.Exec(stmt).Error; err != nil {
			t.Fatalf("exec ddl: %v", err)
		}
	}

	return db
}

func sampleQuizInput() CreateCourseQuizInput {
	visible := 1
	return CreateCourseQuizInput{
		Title:                 "Sample Quiz",
		Instructions:          "Answer all questions.",
		IsPublished:           true,
		RandomizeQuestions:    false,
		SingleQuizView:        true,
		TimeLimit:             1,
		TimeLimitOption:       models.CourseQuizTimeLimitOptionWeek,
		TotalVisibleQuestions: &visible,
		RevealAnswers:         false,
		EnableRetry:           true,
		RetryAttempts:         2,
		MinimumPassPercentage: 50,
		Questions: []CreateQuizQuestionInput{
			{
				Title:          "What is 2+2?",
				Type:           models.QuizQuestionTypeSingleChoice,
				Marks:          1,
				AnswerRequired: true,
			},
		},
	}
}

func seedExistingCourse(t *testing.T, db *gorm.DB) (courseID, chapterID uint) {
	t.Helper()

	category := models.Category{Name: "Test", Slug: "test"}
	if err := db.Create(&category).Error; err != nil {
		t.Fatalf("create category: %v", err)
	}

	course := models.CourseDetails{
		Title:        "Existing Course",
		Slug:         "existing-course-1",
		Summary:      "Summary",
		Visibility:   models.Public,
		PricingModel: models.CoursePricingModelFree,
		AuthorID:     1,
		TenantID:     1,
	}
	if err := db.Create(&course).Error; err != nil {
		t.Fatalf("create course: %v", err)
	}

	chapter := models.CourseChapter{
		CourseID: course.ID,
		Title:    "Chapter 1",
		Position: 0,
		Access:   models.Published,
	}
	if err := db.Create(&chapter).Error; err != nil {
		t.Fatalf("create chapter: %v", err)
	}

	all := models.All
	settings := models.CourseGeneralSettings{
		CourseID:        course.ID,
		DifficultyLevel: &all,
		CategoryID:      category.ID,
	}
	if err := db.Create(&settings).Error; err != nil {
		t.Fatalf("create settings: %v", err)
	}

	return course.ID, chapter.ID
}

func baseCourseInput(chapterID int64, quizzes []CreateCourseQuizInput) CourseDetailsInput {
	chapterIDCopy := chapterID
	return CourseDetailsInput{
		Title:        "Existing Course",
		Summary:      "Summary",
		Visibility:   models.Public,
		IsScheduled:  "false",
		PricingModel: models.CoursePricingModelFree,
		AuthorID:     1,
		CourseChapters: []CreateCourseChapter{
			{
				ID:       &chapterIDCopy,
				Position: 0,
				Title:    "Chapter 1",
				Access:   models.Published,
				Quizzes:  quizzes,
			},
		},
		GeneralSettings: CreateGeneralSettings{
			DifficultyLevel: models.All,
			CategoryID:      1,
		},
		Instructors: []int32{1},
	}
}

func TestCreateCourseWithQuizPersistsQuiz(t *testing.T) {
	db := setupQuizTestDB(t)
	svc := NewCourseService(db)

	category := models.Category{Name: "Create Cat", Slug: "create-cat"}
	if err := db.Create(&category).Error; err != nil {
		t.Fatalf("create category: %v", err)
	}

	input := CourseDetailsInput{
		Title:        "New Course With Quiz",
		Summary:      "Summary",
		Visibility:   models.Public,
		IsScheduled:  "false",
		PricingModel: models.CoursePricingModelFree,
		AuthorID:     1,
		CourseChapters: []CreateCourseChapter{
			{
				Position: 0,
				Title:    "Intro",
				Access:   models.Published,
				Quizzes:  []CreateCourseQuizInput{sampleQuizInput()},
			},
		},
		GeneralSettings: CreateGeneralSettings{
			DifficultyLevel: models.All,
			CategoryID:      category.ID,
		},
		Instructors: []int32{1},
	}

	if err := svc.Create(input, 1, 1); err != nil {
		t.Fatalf("Create() error: %v", err)
	}

	var quizzes []models.CourseQuiz
	if err := db.Preload("Questions").Find(&quizzes).Error; err != nil {
		t.Fatalf("load quizzes: %v", err)
	}
	if len(quizzes) != 1 {
		t.Fatalf("expected 1 quiz, got %d", len(quizzes))
	}
	if quizzes[0].CourseID == 0 {
		t.Fatal("expected quiz.course_id to be set on create")
	}
	if quizzes[0].ChapterID == 0 {
		t.Fatal("expected quiz.chapter_id to be set on create")
	}
	if len(quizzes[0].Questions) != 1 {
		t.Fatalf("expected 1 quiz question, got %d", len(quizzes[0].Questions))
	}
}

func TestUpdateExistingCourseAddsQuizWithCourseID(t *testing.T) {
	db := setupQuizTestDB(t)
	svc := NewCourseService(db)

	courseID, chapterID := seedExistingCourse(t, db)
	input := baseCourseInput(int64(chapterID), []CreateCourseQuizInput{sampleQuizInput()})

	if err := svc.Update(courseID, 1, 1, input); err != nil {
		t.Fatalf("Update() error: %v", err)
	}

	var quizzes []models.CourseQuiz
	if err := db.Preload("Questions").Where("course_id = ?", courseID).Find(&quizzes).Error; err != nil {
		t.Fatalf("load quizzes: %v", err)
	}
	if len(quizzes) != 1 {
		t.Fatalf("expected 1 quiz after update, got %d", len(quizzes))
	}
	if quizzes[0].CourseID != courseID {
		t.Fatalf("quiz.course_id = %d, want %d", quizzes[0].CourseID, courseID)
	}
	if quizzes[0].ChapterID != chapterID {
		t.Fatalf("quiz.chapter_id = %d, want %d", quizzes[0].ChapterID, chapterID)
	}
	if quizzes[0].Title != "Sample Quiz" {
		t.Fatalf("quiz title = %q", quizzes[0].Title)
	}
	if len(quizzes[0].Questions) != 1 {
		t.Fatalf("expected 1 quiz question, got %d", len(quizzes[0].Questions))
	}
}

func TestUpdateExistingCourseCanAddSecondQuiz(t *testing.T) {
	db := setupQuizTestDB(t)
	svc := NewCourseService(db)

	courseID, chapterID := seedExistingCourse(t, db)

	first := sampleQuizInput()
	if err := svc.Update(courseID, 1, 1, baseCourseInput(int64(chapterID), []CreateCourseQuizInput{first})); err != nil {
		t.Fatalf("first Update() error: %v", err)
	}

	var saved []models.CourseQuiz
	if err := db.Where("course_id = ?", courseID).Find(&saved).Error; err != nil {
		t.Fatalf("load first quiz: %v", err)
	}
	if len(saved) != 1 {
		t.Fatalf("expected 1 quiz after first update, got %d", len(saved))
	}

	quizID := int64(saved[0].ID)
	second := sampleQuizInput()
	second.Title = "Second Quiz"

	if err := svc.Update(courseID, 1, 1, baseCourseInput(int64(chapterID), []CreateCourseQuizInput{
		{
			ID:                    &quizID,
			Title:                 saved[0].Title,
			Instructions:          saved[0].Instructions,
			IsPublished:           saved[0].IsPublished,
			RandomizeQuestions:    saved[0].RandomizeQuestions,
			SingleQuizView:        saved[0].SingleQuizView,
			TimeLimit:             saved[0].TimeLimit,
			TimeLimitOption:       saved[0].TimeLimitOption,
			TotalVisibleQuestions: saved[0].TotalVisibleQuestions,
			RevealAnswers:         saved[0].RevealAnswers,
			EnableRetry:           saved[0].EnableRetry,
			RetryAttempts:         saved[0].RetryAttempts,
			MinimumPassPercentage: saved[0].MinimumPassPercentage,
			Questions:             first.Questions,
		},
		second,
	})); err != nil {
		t.Fatalf("second Update() error: %v", err)
	}

	if err := db.Where("course_id = ?", courseID).Find(&saved).Error; err != nil {
		t.Fatalf("load quizzes: %v", err)
	}
	if len(saved) != 2 {
		t.Fatalf("expected 2 quizzes, got %d", len(saved))
	}
}
