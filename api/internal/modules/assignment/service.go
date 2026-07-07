package assignment

import (
	"errors"
	"fmt"
	"math"
	"mime/multipart"
	"strings"
	"time"

	"dashlearn/internal/models"
	"dashlearn/internal/modules/certificate"
	"dashlearn/internal/response"

	"gorm.io/gorm"
)

type AssignmentService interface {
	GetStudentAssignment(tenantID, studentID uint, slug string, assignmentID uint) (*StudentAssignmentResponse, error)
	SubmitAssignment(tenantID, studentID uint, slug string, assignmentID uint, responseText *string, files []*multipart.FileHeader) (*AssignmentSubmissionDetail, error)
	ListCourseSubmissions(tenantID, courseID uint) ([]AssignmentSubmissionListItem, error)
	GetCourseSubmission(tenantID, courseID, submissionID uint) (*AssignmentSubmissionDetail, error)
	GradeSubmission(tenantID, courseID, submissionID uint, input GradeAssignmentInput) (*AssignmentSubmissionDetail, error)
	ListStudentSubmissions(tenantID, studentID uint, courseID *uint) ([]AssignmentSubmissionDetail, error)
	GetStudentSubmission(tenantID, studentID, submissionID uint) (*AssignmentSubmissionDetail, error)
}

type assignmentService struct {
	db     *gorm.DB
	upload func(multipart.File, *multipart.FileHeader) (string, error)
}

func NewAssignmentService(db *gorm.DB, upload func(multipart.File, *multipart.FileHeader) (string, error)) AssignmentService {
	return &assignmentService{db: db, upload: upload}
}

func (s *assignmentService) GetStudentAssignment(tenantID, studentID uint, slug string, assignmentID uint) (*StudentAssignmentResponse, error) {
	_, assignment, err := s.loadPublishedAssignmentForStudent(tenantID, studentID, slug, assignmentID)
	if err != nil {
		return nil, err
	}

	session, err := s.resolveAssignmentSession(tenantID, studentID, assignment)
	if err != nil {
		return nil, err
	}

	var existing models.AssignmentSubmission
	hasSubmitted := false
	if err := s.db.Preload("Files").
		Where("tenant_id = ? AND assignment_id = ? AND student_id = ?", tenantID, assignmentID, studentID).
		First(&existing).Error; err == nil {
		hasSubmitted = true
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	sessionExpired := session.IsExpired()
	editable := hasSubmitted && submissionEditable(existing.Status, sessionExpired)

	resp := &StudentAssignmentResponse{
		CourseAssignmentResponse: buildAssignmentResponse(*assignment),
		HasSubmitted:             hasSubmitted,
		CanSubmit:                (!hasSubmitted && !sessionExpired) || editable,
		CanEdit:                  editable,
		StartedAt:                session.StartedAt.Format(time.RFC3339),
		MaxFileSizeBytes:         MaxAssignmentFileSizeBytes,
		AllowedMimeTypes:         AllowedAssignmentMimeTypes,
		MaxResponseTextLength:    MaxAssignmentResponseTextLength,
	}
	if session.ExpiresAt != nil {
		resp.DeadlineAt = session.ExpiresAt.Format(time.RFC3339)
		resp.SecondsRemaining = secondsRemaining(session.ExpiresAt)
	}

	if hasSubmitted {
		resp.Submission = s.buildSubmissionSummary(existing)
	}

	return resp, nil
}

func (s *assignmentService) SubmitAssignment(tenantID, studentID uint, slug string, assignmentID uint, responseText *string, files []*multipart.FileHeader) (*AssignmentSubmissionDetail, error) {
	_, assignment, err := s.loadPublishedAssignmentForStudent(tenantID, studentID, slug, assignmentID)
	if err != nil {
		return nil, err
	}

	session, err := s.resolveAssignmentSession(tenantID, studentID, assignment)
	if err != nil {
		return nil, err
	}
	if session.IsExpired() {
		return nil, errors.New("assignment time limit exceeded")
	}

	var existing models.AssignmentSubmission
	hasExisting := false
	if err := s.db.Preload("Files").
		Where("tenant_id = ? AND assignment_id = ? AND student_id = ?", tenantID, assignmentID, studentID).
		First(&existing).Error; err == nil {
		hasExisting = true
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if hasExisting {
		if existing.Status == models.AssignmentSubmissionStatusGraded {
			return nil, errors.New("assignment already graded")
		}
		if existing.Status != models.AssignmentSubmissionStatusPendingReview {
			return nil, errors.New("assignment cannot be updated")
		}
		return s.updateSubmission(&existing, assignment, responseText, files)
	}

	return s.createSubmission(tenantID, studentID, assignment, responseText, files)
}

func (s *assignmentService) createSubmission(tenantID, studentID uint, assignment *models.CourseAssignment, responseText *string, files []*multipart.FileHeader) (*AssignmentSubmissionDetail, error) {
	text, uploadedFiles, err := s.prepareSubmissionContent(responseText, files, assignment, nil, false)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	submission := models.AssignmentSubmission{
		TenantID:     tenantID,
		CourseID:     assignment.CourseID,
		ChapterID:    assignment.ChapterID,
		AssignmentID: assignment.ID,
		StudentID:    studentID,
		ResponseText: text,
		MaxScore:     assignment.TotalMarks,
		Status:       models.AssignmentSubmissionStatusPendingReview,
		SubmittedAt:  now,
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&submission).Error; err != nil {
			return err
		}
		for i := range uploadedFiles {
			uploadedFiles[i].SubmissionID = submission.ID
			if err := tx.Create(&uploadedFiles[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	_, _ = certificate.NewService(s.db).TryIssueCertificate(tenantID, studentID, assignment.CourseID)

	return s.buildSubmissionDetail(submission.ID)
}

func (s *assignmentService) updateSubmission(existing *models.AssignmentSubmission, assignment *models.CourseAssignment, responseText *string, files []*multipart.FileHeader) (*AssignmentSubmissionDetail, error) {
	text, uploadedFiles, err := s.prepareSubmissionContent(responseText, files, assignment, existing, true)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	updates := map[string]interface{}{
		"submitted_at": now,
		"status":       models.AssignmentSubmissionStatusPendingReview,
	}
	if text != nil {
		updates["response_text"] = text
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&existing).Updates(updates).Error; err != nil {
			return err
		}
		if len(uploadedFiles) > 0 {
			if err := tx.Where("submission_id = ?", existing.ID).Delete(&models.AssignmentSubmissionFile{}).Error; err != nil {
				return err
			}
			for i := range uploadedFiles {
				uploadedFiles[i].SubmissionID = existing.ID
				if err := tx.Create(&uploadedFiles[i]).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.buildSubmissionDetail(existing.ID)
}

func (s *assignmentService) prepareSubmissionContent(
	responseText *string,
	files []*multipart.FileHeader,
	assignment *models.CourseAssignment,
	existing *models.AssignmentSubmission,
	isUpdate bool,
) (*string, []models.AssignmentSubmissionFile, error) {
	rawText := strings.TrimSpace(stringValue(responseText))
	var sanitizedText string
	var textProvided bool
	if rawText != "" {
		textProvided = true
		var err error
		sanitizedText, err = sanitizeResponseText(rawText)
		if err != nil {
			return nil, nil, err
		}
	}

	limit := assignment.FileUploadLimit
	if limit <= 0 {
		limit = 1
	}
	if len(files) > limit {
		return nil, nil, fmt.Errorf("maximum %d file(s) allowed", limit)
	}

	uploadedFiles := make([]models.AssignmentSubmissionFile, 0, len(files))
	for _, fileHeader := range files {
		if err := validateSubmissionFile(fileHeader); err != nil {
			return nil, nil, err
		}

		file, err := fileHeader.Open()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to open uploaded file: %w", err)
		}

		url, err := s.upload(file, fileHeader)
		file.Close()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to upload file: %w", err)
		}

		mimeType := fileHeader.Header.Get("Content-Type")
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}

		uploadedFiles = append(uploadedFiles, models.AssignmentSubmissionFile{
			URL:      url,
			FileName: fileHeader.Filename,
			MimeType: mimeType,
			Size:     fileHeader.Size,
		})
	}

	existingFiles := 0
	hasExistingText := false
	if existing != nil {
		existingFiles = len(existing.Files)
		hasExistingText = existing.ResponseText != nil && strings.TrimSpace(*existing.ResponseText) != ""
	}

	finalFileCount := existingFiles
	if len(uploadedFiles) > 0 {
		finalFileCount = len(uploadedFiles)
	}

	hasText := sanitizedText != "" || (!textProvided && isUpdate && hasExistingText)
	if !hasText && finalFileCount == 0 {
		return nil, nil, errors.New("response text or at least one file is required")
	}

	var textPtr *string
	if sanitizedText != "" {
		textPtr = &sanitizedText
	}

	return textPtr, uploadedFiles, nil
}

func (s *assignmentService) ListCourseSubmissions(tenantID, courseID uint) ([]AssignmentSubmissionListItem, error) {
	var rows []models.AssignmentSubmission
	err := s.db.
		Preload("Assignment").
		Preload("Student").
		Preload("Files").
		Where("tenant_id = ? AND course_id = ?", tenantID, courseID).
		Order("submitted_at DESC").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	return s.mapSubmissionList(rows)
}

func (s *assignmentService) GetCourseSubmission(tenantID, courseID, submissionID uint) (*AssignmentSubmissionDetail, error) {
	var submission models.AssignmentSubmission
	err := s.db.Where("id = ? AND tenant_id = ? AND course_id = ?", submissionID, tenantID, courseID).First(&submission).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("submission not found")
		}
		return nil, err
	}
	return s.buildSubmissionDetail(submission.ID)
}

func (s *assignmentService) GetStudentSubmission(tenantID, studentID, submissionID uint) (*AssignmentSubmissionDetail, error) {
	var submission models.AssignmentSubmission
	err := s.db.Where("id = ? AND tenant_id = ? AND student_id = ?", submissionID, tenantID, studentID).First(&submission).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("submission not found")
		}
		return nil, err
	}
	return s.buildSubmissionDetail(submission.ID)
}

func (s *assignmentService) GradeSubmission(tenantID, courseID, submissionID uint, input GradeAssignmentInput) (*AssignmentSubmissionDetail, error) {
	var submission models.AssignmentSubmission
	err := s.db.
		Preload("Assignment").
		Where("id = ? AND tenant_id = ? AND course_id = ?", submissionID, tenantID, courseID).
		First(&submission).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("submission not found")
		}
		return nil, err
	}

	if input.Score > submission.Assignment.TotalMarks {
		return nil, fmt.Errorf("score cannot exceed total marks (%.2f)", submission.Assignment.TotalMarks)
	}

	percentage := float32(0)
	if submission.Assignment.TotalMarks > 0 {
		percentage = float32(math.Round(float64(input.Score/submission.Assignment.TotalMarks*100*100)) / 100)
	}
	passed := input.Score >= submission.Assignment.MinimumPassMarks
	now := time.Now()

	submission.Score = input.Score
	submission.MaxScore = submission.Assignment.TotalMarks
	submission.Percentage = percentage
	submission.Passed = passed
	submission.Status = models.AssignmentSubmissionStatusGraded
	submission.InstructorFeedback = input.Feedback
	submission.GradedAt = &now

	if err := s.db.Save(&submission).Error; err != nil {
		return nil, err
	}

	return s.buildSubmissionDetail(submission.ID)
}

func (s *assignmentService) ListStudentSubmissions(tenantID, studentID uint, courseID *uint) ([]AssignmentSubmissionDetail, error) {
	q := s.db.Preload("Assignment").Preload("Student").Preload("Files").
		Where("tenant_id = ? AND student_id = ?", tenantID, studentID)
	if courseID != nil {
		q = q.Where("course_id = ?", *courseID)
	}
	var rows []models.AssignmentSubmission
	if err := q.Order("submitted_at DESC").Find(&rows).Error; err != nil {
		return nil, err
	}

	items := make([]AssignmentSubmissionDetail, 0, len(rows))
	for _, row := range rows {
		detail, err := s.buildSubmissionDetail(row.ID)
		if err != nil {
			return nil, err
		}
		items = append(items, *detail)
	}
	return items, nil
}

func (s *assignmentService) resolveAssignmentSession(tenantID, studentID uint, assignment *models.CourseAssignment) (*models.AssignmentAttemptSession, error) {
	var session models.AssignmentAttemptSession
	err := s.db.Where(
		"tenant_id = ? AND student_id = ? AND assignment_id = ?",
		tenantID, studentID, assignment.ID,
	).First(&session).Error
	if err == nil {
		return &session, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	startedAt := time.Now()
	var expiresAt *time.Time
	if duration := assignmentTimeLimitDuration(assignment.TimeLimit, assignment.TimeLimitOption); duration > 0 {
		expiry := startedAt.Add(duration)
		expiresAt = &expiry
	}

	session = models.AssignmentAttemptSession{
		TenantID:     tenantID,
		StudentID:    studentID,
		AssignmentID: assignment.ID,
		StartedAt:    startedAt,
		ExpiresAt:    expiresAt,
	}
	if err := s.db.Create(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *assignmentService) loadPublishedAssignmentForStudent(tenantID, studentID uint, slug string, assignmentID uint) (*models.CourseDetails, *models.CourseAssignment, error) {
	var course models.CourseDetails
	if err := s.db.Where("tenant_id = ? AND slug = ?", tenantID, slug).First(&course).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("course not found")
		}
		return nil, nil, err
	}

	if !s.isEnrolled(tenantID, studentID, course.ID) {
		return nil, nil, errors.New("enrollment required")
	}

	var assignment models.CourseAssignment
	if err := s.db.
		Where("id = ? AND course_id = ? AND is_published = ?", assignmentID, course.ID, true).
		First(&assignment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("assignment not found")
		}
		return nil, nil, err
	}

	return &course, &assignment, nil
}

func (s *assignmentService) isEnrolled(tenantID, studentID, courseID uint) bool {
	var count int64
	s.db.Model(&models.Enrollment{}).
		Where("tenant_id = ? AND student_id = ? AND course_id = ?", tenantID, studentID, courseID).
		Count(&count)
	return count > 0
}

func (s *assignmentService) mapSubmissionList(rows []models.AssignmentSubmission) ([]AssignmentSubmissionListItem, error) {
	chapterTitles := map[uint]string{}
	items := make([]AssignmentSubmissionListItem, 0, len(rows))
	for _, row := range rows {
		chapterTitle := ""
		if title, ok := chapterTitles[row.ChapterID]; ok {
			chapterTitle = title
		} else {
			var chapter models.CourseChapter
			if err := s.db.Select("title").Where("id = ?", row.ChapterID).First(&chapter).Error; err == nil {
				chapterTitle = chapter.Title
				chapterTitles[row.ChapterID] = chapterTitle
			}
		}

		studentName := row.Student.FirstName
		if row.Student.LastName != nil {
			studentName += " " + *row.Student.LastName
		}

		items = append(items, AssignmentSubmissionListItem{
			ID:              row.ID,
			AssignmentID:    row.AssignmentID,
			AssignmentTitle: row.Assignment.Title,
			ChapterID:       row.ChapterID,
			ChapterTitle:    chapterTitle,
			StudentID:       row.StudentID,
			StudentName:     strings.TrimSpace(studentName),
			StudentEmail:    row.Student.Email,
			Score:           row.Score,
			MaxScore:        row.MaxScore,
			Percentage:      row.Percentage,
			Passed:          row.Passed,
			Status:          row.Status,
			SubmittedAt:     row.SubmittedAt.Format(time.RFC3339),
			FileCount:       len(row.Files),
		})
	}
	return items, nil
}

func (s *assignmentService) buildSubmissionSummary(submission models.AssignmentSubmission) *AssignmentSubmissionSummary {
	files := make([]AssignmentSubmissionFileResponse, 0, len(submission.Files))
	for _, file := range submission.Files {
		files = append(files, AssignmentSubmissionFileResponse{
			ID:       file.ID,
			URL:      file.URL,
			FileName: file.FileName,
			MimeType: file.MimeType,
			Size:     file.Size,
		})
	}

	return &AssignmentSubmissionSummary{
		ID:           submission.ID,
		Score:        submission.Score,
		MaxScore:     submission.MaxScore,
		Percentage:   submission.Percentage,
		Passed:       submission.Passed,
		Status:       submission.Status,
		SubmittedAt:  submission.SubmittedAt.Format(time.RFC3339),
		ResponseText: submission.ResponseText,
		Files:        files,
	}
}

func (s *assignmentService) buildSubmissionDetail(submissionID uint) (*AssignmentSubmissionDetail, error) {
	var submission models.AssignmentSubmission
	if err := s.db.
		Preload("Assignment").
		Preload("Student").
		Preload("Files").
		First(&submission, submissionID).Error; err != nil {
		return nil, err
	}

	list, err := s.mapSubmissionList([]models.AssignmentSubmission{submission})
	if err != nil || len(list) == 0 {
		return nil, errors.New("submission not found")
	}

	summary := s.buildSubmissionSummary(submission)

	return &AssignmentSubmissionDetail{
		AssignmentSubmissionListItem: list[0],
		ResponseText:                 summary.ResponseText,
		InstructorFeedback:           submission.InstructorFeedback,
		Files:                        summary.Files,
	}, nil
}

func buildAssignmentResponse(assignment models.CourseAssignment) response.CourseAssignmentResponse {
	return response.CourseAssignmentResponse{
		ID:               assignment.ID,
		CourseID:         assignment.CourseID,
		ChapterID:        assignment.ChapterID,
		Title:            assignment.Title,
		Instructions:     assignment.Instructions,
		Attachments:      assignment.Attachments,
		IsPublished:      assignment.IsPublished,
		TimeLimit:        assignment.TimeLimit,
		TimeLimitOption:  assignment.TimeLimitOption,
		FileUploadLimit:  assignment.FileUploadLimit,
		TotalMarks:       assignment.TotalMarks,
		MinimumPassMarks: assignment.MinimumPassMarks,
		CreatedAt:        assignment.CreatedAt,
		UpdatedAt:        assignment.UpdatedAt,
	}
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
