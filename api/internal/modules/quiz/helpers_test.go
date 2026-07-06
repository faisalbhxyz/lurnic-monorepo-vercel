package quiz

import (
	"testing"

	"dashlearn/internal/models"
)

func TestBuildQuestionOrderRandomizeAndLimit(t *testing.T) {
	visible := 2
	quiz := &models.CourseQuiz{
		ID:                    9,
		RandomizeQuestions:    true,
		TotalVisibleQuestions: &visible,
	}
	questions := []models.QuizQuestion{
		{ID: 1, Title: "Q1"},
		{ID: 2, Title: "Q2"},
		{ID: 3, Title: "Q3"},
	}

	order := buildQuestionOrder(questions, quiz)
	if len(order) != 2 {
		t.Fatalf("expected 2 visible questions, got %d", len(order))
	}

	seen := map[uint]bool{}
	for _, id := range order {
		if id < 1 || id > 3 {
			t.Fatalf("unexpected question id %d", id)
		}
		seen[id] = true
	}
	if len(seen) != 2 {
		t.Fatalf("expected unique question ids, got %#v", order)
	}
}

func TestQuizTimeLimitDurationZeroMeansUnlimited(t *testing.T) {
	if got := quizTimeLimitDuration(0, models.CourseQuizTimeLimitOptionWeek); got != 0 {
		t.Fatalf("expected 0 duration, got %v", got)
	}
}

func TestQuestionsFromOrderPreservesSessionOrder(t *testing.T) {
	all := []models.QuizQuestion{
		{ID: 1, Title: "Q1"},
		{ID: 2, Title: "Q2"},
		{ID: 3, Title: "Q3"},
	}
	ordered := questionsFromOrder(all, []uint{3, 1})
	if len(ordered) != 2 || ordered[0].ID != 3 || ordered[1].ID != 1 {
		t.Fatalf("unexpected order: %#v", ordered)
	}
}
