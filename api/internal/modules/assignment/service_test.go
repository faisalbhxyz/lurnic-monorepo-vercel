package assignment

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"strings"
	"testing"
	"time"

	"dashlearn/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupAssignmentTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	stmts := []string{
		`CREATE TABLE IF NOT EXISTS tenants (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			app_key TEXT NOT NULL,
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
		`CREATE TABLE IF NOT EXISTS course_assignments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			course_id INTEGER NOT NULL,
			chapter_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			instructions TEXT NOT NULL,
			attachments TEXT,
			is_published INTEGER DEFAULT 0,
			time_limit INTEGER DEFAULT 1,
			time_limit_option TEXT DEFAULT 'weeks',
			file_upload_limit INTEGER DEFAULT 1,
			total_marks REAL DEFAULT 10,
			minimum_pass_marks REAL DEFAULT 6,
			position INTEGER DEFAULT 0,
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
		`CREATE TABLE IF NOT EXISTS assignment_submissions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			tenant_id INTEGER NOT NULL,
			course_id INTEGER NOT NULL,
			chapter_id INTEGER NOT NULL,
			assignment_id INTEGER NOT NULL,
			student_id INTEGER NOT NULL,
			response_text TEXT,
			score REAL DEFAULT 0,
			max_score REAL DEFAULT 0,
			percentage REAL DEFAULT 0,
			passed INTEGER DEFAULT 0,
			status TEXT DEFAULT 'pending_review',
			instructor_feedback TEXT,
			submitted_at DATETIME,
			graded_at DATETIME,
			created_at DATETIME,
			updated_at DATETIME,
			UNIQUE(assignment_id, student_id)
		)`,
		`CREATE TABLE IF NOT EXISTS assignment_submission_files (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			submission_id INTEGER NOT NULL,
			url TEXT NOT NULL,
			file_name TEXT NOT NULL,
			mime_type TEXT NOT NULL,
			size INTEGER DEFAULT 0,
			created_at DATETIME,
			updated_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS assignment_attempt_sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			tenant_id INTEGER NOT NULL,
			student_id INTEGER NOT NULL,
			assignment_id INTEGER NOT NULL,
			started_at DATETIME NOT NULL,
			expires_at DATETIME,
			created_at DATETIME,
			updated_at DATETIME,
			UNIQUE(tenant_id, student_id, assignment_id)
		)`,
	}
	for _, stmt := range stmts {
		if err := db.Exec(stmt).Error; err != nil {
			t.Fatalf("exec schema: %v", err)
		}
	}

	return db
}

func seedAssignmentFixture(t *testing.T, db *gorm.DB) (tenantID, courseID, chapterID, assignmentID, studentID uint) {
	t.Helper()

	course := models.CourseDetails{
		Title:        "Test Course",
		Slug:         "test-course",
		Summary:      "summary",
		Visibility:   models.Public,
		PricingModel: models.CoursePricingModelFree,
		AuthorID:     1,
		TenantID:     1,
	}
	if err := db.Create(&course).Error; err != nil {
		t.Fatalf("create course: %v", err)
	}

	chapter := models.CourseChapter{
		Title:    "Chapter 1",
		Access:   models.Published,
		CourseID: course.ID,
	}
	if err := db.Create(&chapter).Error; err != nil {
		t.Fatalf("create chapter: %v", err)
	}

	assignment := models.CourseAssignment{
		CourseID:         course.ID,
		ChapterID:        chapter.ID,
		Title:            "Homework 1",
		Instructions:     "<p>Submit your work</p>",
		IsPublished:      true,
		TimeLimit:        1,
		TimeLimitOption:  models.CourseAssignmentTimeLimitOptionWeek,
		FileUploadLimit:  2,
		TotalMarks:       10,
		MinimumPassMarks: 6,
	}
	if err := db.Create(&assignment).Error; err != nil {
		t.Fatalf("create assignment: %v", err)
	}

	lastName := "Student"
	student := models.Student{
		UserID:    "student-1",
		TenantID:  1,
		FirstName: "Test",
		LastName:  &lastName,
		Email:     "student@test.com",
		Password:  "hash",
	}
	if err := db.Create(&student).Error; err != nil {
		t.Fatalf("create student: %v", err)
	}

	enrollment := models.Enrollment{
		TenantID:  1,
		StudentID: student.ID,
		CourseID:  course.ID,
	}
	if err := db.Create(&enrollment).Error; err != nil {
		t.Fatalf("create enrollment: %v", err)
	}

	return 1, course.ID, chapter.ID, assignment.ID, student.ID
}

func mockUpload(_ multipart.File, header *multipart.FileHeader) (string, error) {
	return "https://cdn.example.com/" + header.Filename, nil
}

func makeFileHeader(t *testing.T, name, contentType string, body []byte) *multipart.FileHeader {
	t.Helper()
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("files", name)
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := part.Write(body); err != nil {
		t.Fatalf("write form file: %v", err)
	}
	_ = writer.Close()

	reader := multipart.NewReader(&buf, writer.Boundary())
	form, err := reader.ReadForm(1 << 20)
	if err != nil {
		t.Fatalf("read form: %v", err)
	}
	fileHeaders := form.File["files"]
	if len(fileHeaders) == 0 {
		t.Fatal("expected file header")
	}
	fileHeaders[0].Header.Set("Content-Type", contentType)
	return fileHeaders[0]
}

func TestGetStudentAssignment_Enrolled(t *testing.T) {
	db := setupAssignmentTestDB(t)
	tenantID, _, _, assignmentID, studentID := seedAssignmentFixture(t, db)
	svc := NewAssignmentService(db, mockUpload)

	resp, err := svc.GetStudentAssignment(tenantID, studentID, "test-course", assignmentID)
	if err != nil {
		t.Fatalf("GetStudentAssignment: %v", err)
	}
	if resp.Title != "Homework 1" {
		t.Fatalf("expected title Homework 1, got %q", resp.Title)
	}
	if resp.HasSubmitted {
		t.Fatal("expected has_submitted false")
	}
	if !resp.CanSubmit {
		t.Fatal("expected can_submit true")
	}
	if resp.DeadlineAt == "" {
		t.Fatal("expected deadline_at to be set")
	}
	if resp.SecondsRemaining == nil {
		t.Fatal("expected seconds_remaining to be set")
	}
	if resp.MaxFileSizeBytes != MaxAssignmentFileSizeBytes {
		t.Fatalf("unexpected max_file_size_bytes: %d", resp.MaxFileSizeBytes)
	}
}

func TestSubmitAndGradeAssignment_Flow(t *testing.T) {
	db := setupAssignmentTestDB(t)
	tenantID, courseID, _, assignmentID, studentID := seedAssignmentFixture(t, db)
	svc := NewAssignmentService(db, mockUpload)

	responseText := "<p>My answer</p>"
	file := makeFileHeader(t, "work.pdf", "application/pdf", []byte("%PDF-1.4 test"))

	submitted, err := svc.SubmitAssignment(
		tenantID,
		studentID,
		"test-course",
		assignmentID,
		&responseText,
		[]*multipart.FileHeader{file},
	)
	if err != nil {
		t.Fatalf("SubmitAssignment: %v", err)
	}
	if submitted.Status != models.AssignmentSubmissionStatusPendingReview {
		t.Fatalf("expected pending_review, got %s", submitted.Status)
	}
	if submitted.ResponseText == nil || *submitted.ResponseText != responseText {
		t.Fatal("expected response text on submission detail")
	}
	if len(submitted.Files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(submitted.Files))
	}
	if submitted.Files[0].FileName != "work.pdf" {
		t.Fatalf("unexpected file name %q", submitted.Files[0].FileName)
	}

	afterSubmit, err := svc.GetStudentAssignment(tenantID, studentID, "test-course", assignmentID)
	if err != nil {
		t.Fatalf("GetStudentAssignment after submit: %v", err)
	}
	if !afterSubmit.HasSubmitted || !afterSubmit.CanSubmit || !afterSubmit.CanEdit {
		t.Fatal("expected submitted with can_submit and can_edit true while pending review")
	}
	if afterSubmit.Submission == nil || afterSubmit.Submission.ResponseText == nil {
		t.Fatal("expected full submission content on detail GET")
	}
	if len(afterSubmit.Submission.Files) != 1 {
		t.Fatalf("expected 1 file on detail submission, got %d", len(afterSubmit.Submission.Files))
	}

	updatedText := "<p>Updated answer</p>"
	resubmitted, err := svc.SubmitAssignment(tenantID, studentID, "test-course", assignmentID, &updatedText, nil)
	if err != nil {
		t.Fatalf("resubmit assignment: %v", err)
	}
	if resubmitted.ResponseText == nil || *resubmitted.ResponseText != updatedText {
		t.Fatalf("expected updated response text, got %+v", resubmitted.ResponseText)
	}
	if len(resubmitted.Files) != 1 {
		t.Fatalf("expected existing file retained after text-only resubmit, got %d", len(resubmitted.Files))
	}

	afterGradeBlock, err := svc.GetStudentAssignment(tenantID, studentID, "test-course", assignmentID)
	if err != nil {
		t.Fatalf("GetStudentAssignment after resubmit: %v", err)
	}
	if !afterGradeBlock.CanSubmit || !afterGradeBlock.CanEdit {
		t.Fatal("expected can still edit before grading")
	}

	list, err := svc.ListCourseSubmissions(tenantID, courseID)
	if err != nil {
		t.Fatalf("ListCourseSubmissions: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 submission in admin list, got %d", len(list))
	}
	if list[0].AssignmentTitle != "Homework 1" {
		t.Fatalf("unexpected assignment title %q", list[0].AssignmentTitle)
	}
	if list[0].FileCount != 1 {
		t.Fatalf("expected file_count 1, got %d", list[0].FileCount)
	}

	feedback := "Well done"
	graded, err := svc.GradeSubmission(tenantID, courseID, submitted.ID, GradeAssignmentInput{
		Score:    8,
		Feedback: &feedback,
	})
	if err != nil {
		t.Fatalf("GradeSubmission: %v", err)
	}
	if graded.Status != models.AssignmentSubmissionStatusGraded {
		t.Fatalf("expected graded status, got %s", graded.Status)
	}
	if graded.Score != 8 || graded.Percentage != 80 {
		t.Fatalf("unexpected score/percentage: %.2f / %.2f", graded.Score, graded.Percentage)
	}
	if !graded.Passed {
		t.Fatal("expected passed true for score 8 with min 6")
	}
	if graded.InstructorFeedback == nil || *graded.InstructorFeedback != feedback {
		t.Fatal("expected instructor feedback")
	}

	studentList, err := svc.ListStudentSubmissions(tenantID, studentID, &courseID)
	if err != nil {
		t.Fatalf("ListStudentSubmissions: %v", err)
	}
	if len(studentList) != 1 || studentList[0].Status != models.AssignmentSubmissionStatusGraded {
		t.Fatalf("unexpected student submission list: %+v", studentList)
	}
	if studentList[0].ResponseText == nil || len(studentList[0].Files) != 1 {
		t.Fatal("expected full submission detail in student list")
	}

	detail, err := svc.GetStudentSubmission(tenantID, studentID, submitted.ID)
	if err != nil {
		t.Fatalf("GetStudentSubmission: %v", err)
	}
	if detail.Status != models.AssignmentSubmissionStatusGraded {
		t.Fatalf("expected graded detail, got %s", detail.Status)
	}
}

func TestSubmitAssignment_RequiresContent(t *testing.T) {
	db := setupAssignmentTestDB(t)
	tenantID, _, _, assignmentID, studentID := seedAssignmentFixture(t, db)
	svc := NewAssignmentService(db, mockUpload)

	_, err := svc.SubmitAssignment(tenantID, studentID, "test-course", assignmentID, nil, nil)
	if err == nil || err.Error() != "response text or at least one file is required" {
		t.Fatalf("expected empty submission error, got %v", err)
	}
}

func TestSubmitAssignment_RejectsOversizedFile(t *testing.T) {
	db := setupAssignmentTestDB(t)
	tenantID, _, _, assignmentID, studentID := seedAssignmentFixture(t, db)
	svc := NewAssignmentService(db, mockUpload)

	body := make([]byte, MaxAssignmentFileSizeBytes+1)
	file := makeFileHeader(t, "big.pdf", "application/pdf", body)
	_, err := svc.SubmitAssignment(tenantID, studentID, "test-course", assignmentID, nil, []*multipart.FileHeader{file})
	if err == nil || !containsErr(err, "exceeds maximum size") {
		t.Fatalf("expected file size error, got %v", err)
	}
}

func TestSubmitAssignment_RejectsDisallowedMime(t *testing.T) {
	db := setupAssignmentTestDB(t)
	tenantID, _, _, assignmentID, studentID := seedAssignmentFixture(t, db)
	svc := NewAssignmentService(db, mockUpload)

	file := makeFileHeader(t, "evil.html", "text/html", []byte("<html></html>"))
	_, err := svc.SubmitAssignment(tenantID, studentID, "test-course", assignmentID, nil, []*multipart.FileHeader{file})
	if err == nil || !containsErr(err, "not allowed") {
		t.Fatalf("expected mime error, got %v", err)
	}
}

func TestSubmitAssignment_RejectsAfterTimeLimit(t *testing.T) {
	db := setupAssignmentTestDB(t)
	tenantID, _, _, assignmentID, studentID := seedAssignmentFixture(t, db)
	svc := NewAssignmentService(db, mockUpload)

	past := time.Now().Add(-2 * time.Hour)
	session := models.AssignmentAttemptSession{
		TenantID:     tenantID,
		StudentID:    studentID,
		AssignmentID: assignmentID,
		StartedAt:    past.Add(-1 * time.Hour),
		ExpiresAt:    &past,
	}
	if err := db.Create(&session).Error; err != nil {
		t.Fatalf("create expired session: %v", err)
	}

	text := "<p>late</p>"
	_, err := svc.SubmitAssignment(tenantID, studentID, "test-course", assignmentID, &text, nil)
	if err == nil || err.Error() != "assignment time limit exceeded" {
		t.Fatalf("expected time limit error, got %v", err)
	}
}

func containsErr(err error, sub string) bool {
	return err != nil && strings.Contains(err.Error(), sub)
}

func TestGetStudentAssignment_NotEnrolled(t *testing.T) {
	db := setupAssignmentTestDB(t)
	_, _, _, assignmentID, _ := seedAssignmentFixture(t, db)
	svc := NewAssignmentService(db, mockUpload)

	_, err := svc.GetStudentAssignment(1, 999, "test-course", assignmentID)
	if err == nil || err.Error() != "enrollment required" {
		t.Fatalf("expected enrollment required, got %v", err)
	}
}
