package assignment

import (
	"testing"
	"time"

	"dashlearn/internal/models"
)

func TestAssignmentTimeLimitDuration(t *testing.T) {
	if got := assignmentTimeLimitDuration(0, models.CourseAssignmentTimeLimitOptionWeek); got != 0 {
		t.Fatalf("expected 0 duration, got %v", got)
	}
	if got := assignmentTimeLimitDuration(2, models.CourseAssignmentTimeLimitOptionHour); got != 2*time.Hour {
		t.Fatalf("expected 2h, got %v", got)
	}
}

func TestIsAllowedMimeType(t *testing.T) {
	cases := []struct {
		mime    string
		allowed bool
	}{
		{"application/pdf", true},
		{"image/jpeg", true},
		{"application/octet-stream", false},
		{"text/html", false},
	}
	for _, tc := range cases {
		if got := isAllowedMimeType(tc.mime); got != tc.allowed {
			t.Fatalf("mime %q: expected %v, got %v", tc.mime, tc.allowed, got)
		}
	}
}

func TestSanitizeResponseText_StripsScript(t *testing.T) {
	out, err := sanitizeResponseText(`<p>Hello</p><script>alert(1)</script>`)
	if err != nil {
		t.Fatalf("sanitize: %v", err)
	}
	if out == "" || contains(out, "script") {
		t.Fatalf("unexpected sanitized output: %q", out)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 || indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

func TestSubmissionEditable(t *testing.T) {
	if !submissionEditable(models.AssignmentSubmissionStatusPendingReview, false) {
		t.Fatal("expected editable while pending review")
	}
	if submissionEditable(models.AssignmentSubmissionStatusGraded, false) {
		t.Fatal("expected not editable when graded")
	}
	if submissionEditable(models.AssignmentSubmissionStatusPendingReview, true) {
		t.Fatal("expected not editable when session expired")
	}
}
