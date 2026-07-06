package quiz

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	mathrand "math/rand"
	"sort"
	"time"

	"dashlearn/internal/models"
)

func quizTimeLimitDuration(limit int, option models.CourseQuizTimeLimitOption) time.Duration {
	if limit <= 0 {
		return 0
	}
	switch option {
	case models.CourseQuizTimeLimitOptionMinute:
		return time.Duration(limit) * time.Minute
	case models.CourseQuizTimeLimitOptionHour:
		return time.Duration(limit) * time.Hour
	case models.CourseQuizTimeLimitOptionDay:
		return time.Duration(limit) * 24 * time.Hour
	case models.CourseQuizTimeLimitOptionWeek:
		return time.Duration(limit) * 7 * 24 * time.Hour
	case models.CourseQuizTimeLimitOptionMonth:
		return time.Duration(limit) * 30 * 24 * time.Hour
	default:
		return time.Duration(limit) * time.Minute
	}
}

func buildQuestionOrder(questions []models.QuizQuestion, quiz *models.CourseQuiz) []uint {
	items := make([]models.QuizQuestion, len(questions))
	copy(items, questions)

	if quiz.RandomizeQuestions {
		shuffleQuestions(items)
	} else {
		sort.Slice(items, func(i, j int) bool {
			return items[i].ID < items[j].ID
		})
	}

	if quiz.TotalVisibleQuestions != nil &&
		*quiz.TotalVisibleQuestions > 0 &&
		*quiz.TotalVisibleQuestions < len(items) {
		items = items[:*quiz.TotalVisibleQuestions]
	}

	order := make([]uint, len(items))
	for i, q := range items {
		order[i] = q.ID
	}
	return order
}

func shuffleQuestions(items []models.QuizQuestion) {
	var seed uint64
	if err := binary.Read(rand.Reader, binary.LittleEndian, &seed); err != nil {
		seed = uint64(time.Now().UnixNano())
	}
	r := mathrand.New(mathrand.NewSource(int64(seed)))
	r.Shuffle(len(items), func(i, j int) {
		items[i], items[j] = items[j], items[i]
	})
}

func questionsFromOrder(all []models.QuizQuestion, order []uint) []models.QuizQuestion {
	byID := make(map[uint]models.QuizQuestion, len(all))
	for _, q := range all {
		byID[q.ID] = q
	}
	out := make([]models.QuizQuestion, 0, len(order))
	for _, id := range order {
		if q, ok := byID[id]; ok {
			out = append(out, q)
		}
	}
	return out
}

func decodeQuestionOrder(raw []byte) ([]uint, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	var order []uint
	if err := json.Unmarshal(raw, &order); err != nil {
		return nil, err
	}
	return order, nil
}

func secondsRemaining(expiresAt *time.Time) *int {
	if expiresAt == nil {
		return nil
	}
	remaining := int(time.Until(*expiresAt).Seconds())
	if remaining < 0 {
		remaining = 0
	}
	return &remaining
}
