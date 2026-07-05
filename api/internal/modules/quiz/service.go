package quiz

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"dashlearn/internal/models"
	"dashlearn/internal/modules/certificate"
	"dashlearn/internal/response"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type QuizService interface {
	GetStudentQuiz(tenantID, studentID uint, slug string, quizID uint) (*StudentQuizResponse, error)
	SubmitQuiz(tenantID, studentID uint, slug string, quizID uint, input SubmitQuizInput) (*QuizSubmissionDetail, error)
	ListCourseSubmissions(tenantID, courseID uint) ([]QuizSubmissionListItem, error)
	GetCourseSubmission(tenantID, courseID, submissionID uint) (*QuizSubmissionDetail, error)
	ListStudentSubmissions(tenantID, studentID uint, courseID *uint) ([]QuizSubmissionListItem, error)
}

type quizService struct {
	db *gorm.DB
}

func NewQuizService(db *gorm.DB) QuizService {
	return &quizService{db: db}
}

func (s *quizService) GetStudentQuiz(tenantID, studentID uint, slug string, quizID uint) (*StudentQuizResponse, error) {
	course, quiz, err := s.loadPublishedQuizForStudent(tenantID, studentID, slug, quizID)
	if err != nil {
		return nil, err
	}

	attemptsUsed, err := s.countAttempts(tenantID, studentID, quizID)
	if err != nil {
		return nil, err
	}
	if !quiz.EnableRetry && attemptsUsed > 0 {
		return nil, errors.New("quiz retry is disabled")
	}
	if quiz.EnableRetry && quiz.RetryAttempts > 0 && attemptsUsed >= quiz.RetryAttempts {
		return nil, errors.New("maximum quiz attempts reached")
	}

	canRetry := quiz.EnableRetry && (quiz.RetryAttempts == 0 || attemptsUsed < quiz.RetryAttempts)

	questions := prepareQuestionsForAttempt(quiz.Questions, quiz)
	resp := &StudentQuizResponse{
		CourseQuizResponse: buildQuizResponse(*quiz, questions, false),
		AttemptsUsed:       attemptsUsed,
		CanRetry:           canRetry,
	}
	_ = course
	return resp, nil
}

func (s *quizService) SubmitQuiz(tenantID, studentID uint, slug string, quizID uint, input SubmitQuizInput) (*QuizSubmissionDetail, error) {
	_, quiz, err := s.loadPublishedQuizForStudent(tenantID, studentID, slug, quizID)
	if err != nil {
		return nil, err
	}

	attemptsUsed, err := s.countAttempts(tenantID, studentID, quizID)
	if err != nil {
		return nil, err
	}
	if !quiz.EnableRetry && attemptsUsed > 0 {
		return nil, errors.New("quiz retry is disabled")
	}
	if quiz.EnableRetry && quiz.RetryAttempts > 0 && attemptsUsed >= quiz.RetryAttempts {
		return nil, errors.New("maximum quiz attempts reached")
	}

	questionMap := make(map[uint]models.QuizQuestion, len(quiz.Questions))
	for _, q := range quiz.Questions {
		questionMap[q.ID] = q
	}

	answerByQuestion := make(map[uint]SubmitQuizAnswerInput, len(input.Answers))
	for _, ans := range input.Answers {
		if _, exists := questionMap[ans.QuestionID]; !exists {
			return nil, fmt.Errorf("invalid question_id: %d", ans.QuestionID)
		}
		answerByQuestion[ans.QuestionID] = ans
	}

	for _, q := range quiz.Questions {
		if q.AnswerRequired {
			if _, ok := answerByQuestion[q.ID]; !ok {
				return nil, fmt.Errorf("answer required for question %d", q.ID)
			}
		}
	}

	attemptNumber := attemptsUsed + 1
	now := time.Now()
	submission := models.QuizSubmission{
		TenantID:      tenantID,
		CourseID:      quiz.CourseID,
		ChapterID:     quiz.ChapterID,
		QuizID:        quiz.ID,
		StudentID:     studentID,
		AttemptNumber: attemptNumber,
		SubmittedAt:   now,
		Status:        models.QuizSubmissionStatusSubmitted,
	}

	var score float32
	var maxScore float32
	pendingReview := false
	answerRows := make([]models.QuizSubmissionAnswer, 0, len(quiz.Questions))

	for _, q := range quiz.Questions {
		maxScore += q.Marks
		submitted, hasAnswer := answerByQuestion[q.ID]

		var answerJSON []byte
		var isCorrect *bool
		marksAwarded := float32(0)

		if hasAnswer {
			answerJSON, _ = json.Marshal(submitted.Value)
			if q.CorrectAnswer != nil {
				correct, graded, gradeErr := gradeAnswer(q, submitted.Value)
				if gradeErr != nil {
					pendingReview = true
				} else if graded {
					isCorrect = &correct
					if correct {
						marksAwarded = q.Marks
						score += q.Marks
					}
				} else {
					pendingReview = true
				}
			} else {
				pendingReview = true
			}
		} else if q.AnswerRequired {
			pendingReview = true
		}

		answerRows = append(answerRows, models.QuizSubmissionAnswer{
			QuestionID:   q.ID,
			Answer:       answerJSON,
			IsCorrect:    isCorrect,
			MarksAwarded: marksAwarded,
		})
	}

	percentage := float32(0)
	if maxScore > 0 {
		percentage = float32(math.Round(float64(score/maxScore*100*100)) / 100)
	}
	passed := percentage >= quiz.MinimumPassPercentage

	submission.Score = score
	submission.MaxScore = maxScore
	submission.Percentage = percentage
	submission.Passed = passed
	if pendingReview {
		submission.Status = models.QuizSubmissionStatusPendingReview
	} else {
		submission.Status = models.QuizSubmissionStatusGraded
		gradedAt := now
		submission.GradedAt = &gradedAt
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&submission).Error; err != nil {
			return err
		}
		for i := range answerRows {
			answerRows[i].SubmissionID = submission.ID
			if err := tx.Create(&answerRows[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	_, _ = certificate.NewService(s.db).TryIssueCertificate(tenantID, studentID, quiz.CourseID)

	return s.buildSubmissionDetail(submission.ID, quiz.RevealAnswers)
}

func (s *quizService) ListCourseSubmissions(tenantID, courseID uint) ([]QuizSubmissionListItem, error) {
	var rows []models.QuizSubmission
	err := s.db.
		Preload("Quiz").
		Preload("Student").
		Where("tenant_id = ? AND course_id = ?", tenantID, courseID).
		Order("submitted_at DESC").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	return s.mapSubmissionList(rows)
}

func (s *quizService) GetCourseSubmission(tenantID, courseID, submissionID uint) (*QuizSubmissionDetail, error) {
	var submission models.QuizSubmission
	err := s.db.Where("id = ? AND tenant_id = ? AND course_id = ?", submissionID, tenantID, courseID).First(&submission).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("submission not found")
		}
		return nil, err
	}
	return s.buildSubmissionDetail(submission.ID, true)
}

func (s *quizService) ListStudentSubmissions(tenantID, studentID uint, courseID *uint) ([]QuizSubmissionListItem, error) {
	q := s.db.Preload("Quiz").Preload("Student").Where("tenant_id = ? AND student_id = ?", tenantID, studentID)
	if courseID != nil {
		q = q.Where("course_id = ?", *courseID)
	}
	var rows []models.QuizSubmission
	if err := q.Order("submitted_at DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return s.mapSubmissionList(rows)
}

func (s *quizService) loadPublishedQuizForStudent(tenantID, studentID uint, slug string, quizID uint) (*models.CourseDetails, *models.CourseQuiz, error) {
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

	var quiz models.CourseQuiz
	if err := s.db.
		Preload("Questions").
		Where("id = ? AND course_id = ? AND is_published = ?", quizID, course.ID, true).
		First(&quiz).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("quiz not found")
		}
		return nil, nil, err
	}
	if len(quiz.Questions) == 0 {
		return nil, nil, errors.New("quiz has no questions")
	}

	return &course, &quiz, nil
}

func (s *quizService) isEnrolled(tenantID, studentID, courseID uint) bool {
	var count int64
	s.db.Model(&models.Enrollment{}).
		Where("tenant_id = ? AND student_id = ? AND course_id = ?", tenantID, studentID, courseID).
		Count(&count)
	return count > 0
}

func (s *quizService) countAttempts(tenantID, studentID, quizID uint) (int, error) {
	var count int64
	err := s.db.Model(&models.QuizSubmission{}).
		Where("tenant_id = ? AND student_id = ? AND quiz_id = ?", tenantID, studentID, quizID).
		Count(&count).Error
	return int(count), err
}

func (s *quizService) mapSubmissionList(rows []models.QuizSubmission) ([]QuizSubmissionListItem, error) {
	chapterTitles := map[uint]string{}
	items := make([]QuizSubmissionListItem, 0, len(rows))
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

		items = append(items, QuizSubmissionListItem{
			ID:            row.ID,
			QuizID:        row.QuizID,
			QuizTitle:     row.Quiz.Title,
			ChapterID:     row.ChapterID,
			ChapterTitle:  chapterTitle,
			StudentID:     row.StudentID,
			StudentName:   strings.TrimSpace(studentName),
			StudentEmail:  row.Student.Email,
			AttemptNumber: row.AttemptNumber,
			Score:         row.Score,
			MaxScore:      row.MaxScore,
			Percentage:    row.Percentage,
			Passed:        row.Passed,
			Status:        row.Status,
			SubmittedAt:   row.SubmittedAt.Format(time.RFC3339),
		})
	}
	return items, nil
}

func (s *quizService) buildSubmissionDetail(submissionID uint, revealAnswers bool) (*QuizSubmissionDetail, error) {
	var submission models.QuizSubmission
	if err := s.db.
		Preload("Quiz").
		Preload("Student").
		Preload("Answers").
		Preload("Answers.Question").
		First(&submission, submissionID).Error; err != nil {
		return nil, err
	}

	list, err := s.mapSubmissionList([]models.QuizSubmission{submission})
	if err != nil || len(list) == 0 {
		return nil, errors.New("submission not found")
	}

	answers := make([]QuizSubmissionAnswerResponse, 0, len(submission.Answers))
	sort.Slice(submission.Answers, func(i, j int) bool {
		return submission.Answers[i].QuestionID < submission.Answers[j].QuestionID
	})

	for _, ans := range submission.Answers {
		var submitted interface{}
		_ = json.Unmarshal(ans.Answer, &submitted)

		item := QuizSubmissionAnswerResponse{
			QuestionID:      ans.QuestionID,
			QuestionTitle:   ans.Question.Title,
			QuestionType:    ans.Question.Type,
			SubmittedAnswer: submitted,
			IsCorrect:       ans.IsCorrect,
			MarksAwarded:    ans.MarksAwarded,
		}
		if revealAnswers || submission.Quiz.RevealAnswers {
			item.AnswerExplanation = ans.Question.AnswerExplanation
			if ans.Question.CorrectAnswer != nil {
				_ = json.Unmarshal(*ans.Question.CorrectAnswer, &item.CorrectAnswer)
			}
		}
		answers = append(answers, item)
	}

	return &QuizSubmissionDetail{
		QuizSubmissionListItem: list[0],
		RevealAnswers:          revealAnswers || submission.Quiz.RevealAnswers,
		Answers:                answers,
	}, nil
}

func prepareQuestionsForAttempt(questions []models.QuizQuestion, quiz *models.CourseQuiz) []models.QuizQuestion {
	items := make([]models.QuizQuestion, len(questions))
	copy(items, questions)
	if quiz.RandomizeQuestions {
		// simple shuffle by id xor for deterministic tests; good enough for MVP
		sort.Slice(items, func(i, j int) bool {
			return (items[i].ID ^ uint(quiz.ID))%7 < (items[j].ID ^ uint(quiz.ID))%7
		})
	}
	if quiz.TotalVisibleQuestions != nil && *quiz.TotalVisibleQuestions > 0 && *quiz.TotalVisibleQuestions < len(items) {
		items = items[:*quiz.TotalVisibleQuestions]
	}
	return items
}

func buildQuizResponse(quiz models.CourseQuiz, questions []models.QuizQuestion, revealSensitive bool) response.CourseQuizResponse {
	questionResponses := make([]response.CourseQuizQuestionsResponse, 0, len(questions))
	for _, question := range questions {
		questionResponses = append(questionResponses, sanitizeQuestionResponse(question, revealSensitive))
	}
	return response.CourseQuizResponse{
		ID:                    quiz.ID,
		CourseID:              quiz.CourseID,
		ChapterID:             quiz.ChapterID,
		Title:                 quiz.Title,
		Instructions:          quiz.Instructions,
		IsPublished:           quiz.IsPublished,
		RandomizeQuestions:    quiz.RandomizeQuestions,
		SingleQuizView:        quiz.SingleQuizView,
		TimeLimit:             quiz.TimeLimit,
		TimeLimitOption:       quiz.TimeLimitOption,
		TotalVisibleQuestions: quiz.TotalVisibleQuestions,
		RevealAnswers:         quiz.RevealAnswers,
		EnableRetry:           quiz.EnableRetry,
		RetryAttempts:         quiz.RetryAttempts,
		MinimumPassPercentage: quiz.MinimumPassPercentage,
		CreatedAt:             quiz.CreatedAt,
		UpdatedAt:             quiz.UpdatedAt,
		Questions:             questionResponses,
	}
}

func SanitizeQuestionResponse(question models.QuizQuestion, revealSensitive bool) response.CourseQuizQuestionsResponse {
	return sanitizeQuestionResponse(question, revealSensitive)
}

func sanitizeQuestionResponse(question models.QuizQuestion, revealSensitive bool) response.CourseQuizQuestionsResponse {
	res := response.CourseQuizQuestionsResponse{
		ID:             question.ID,
		QuizID:         question.QuizID,
		Title:          question.Title,
		Details:        question.Details,
		Media:          question.Media,
		Options:        question.Options,
		Type:           question.Type,
		Marks:          question.Marks,
		AnswerRequired: question.AnswerRequired,
		CreatedAt:      question.CreatedAt,
		UpdatedAt:      question.UpdatedAt,
	}
	if revealSensitive {
		res.AnswerExplanation = question.AnswerExplanation
		res.CorrectAnswer = question.CorrectAnswer
	}
	return res
}

func gradeAnswer(question models.QuizQuestion, submitted interface{}) (correct bool, graded bool, err error) {
	if question.CorrectAnswer == nil {
		return false, false, nil
	}

	var expected map[string]interface{}
	if err := json.Unmarshal(*question.CorrectAnswer, &expected); err != nil {
		return false, false, err
	}

	switch question.Type {
	case models.QuizQuestionTypeTrueFalse:
		expectedVal, ok := expected["value"]
		if !ok {
			return false, false, nil
		}
		return compareScalar(expectedVal, submitted), true, nil
	case models.QuizQuestionTypeSingleChoice:
		expectedVal, ok := expected["value"]
		if !ok {
			return false, false, nil
		}
		return compareScalar(expectedVal, submitted), true, nil
	case models.QuizQuestionTypeMultipleChoice:
		expectedVals, ok := expected["values"].([]interface{})
		if !ok {
			return false, false, nil
		}
		submittedVals, ok := toStringSlice(submitted)
		if !ok {
			return false, true, nil
		}
		return compareStringSets(expectedVals, submittedVals), true, nil
	default:
		return false, false, nil
	}
}

func compareScalar(expected interface{}, submitted interface{}) bool {
	return fmt.Sprint(expected) == fmt.Sprint(submitted)
}

func toStringSlice(value interface{}) ([]string, bool) {
	switch v := value.(type) {
	case []interface{}:
		out := make([]string, 0, len(v))
		for _, item := range v {
			out = append(out, fmt.Sprint(item))
		}
		return out, true
	case []string:
		return v, true
	default:
		return nil, false
	}
}

func compareStringSets(expected []interface{}, submitted []string) bool {
	if len(expected) != len(submitted) {
		return false
	}
	expectedSet := map[string]int{}
	for _, item := range expected {
		expectedSet[fmt.Sprint(item)]++
	}
	for _, item := range submitted {
		expectedSet[item]--
	}
	for _, count := range expectedSet {
		if count != 0 {
			return false
		}
	}
	return true
}

// Ensure JSON columns marshal cleanly when saving from admin later.
func toJSONColumn(v interface{}) *datatypes.JSON {
	if v == nil {
		return nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	j := datatypes.JSON(b)
	return &j
}
