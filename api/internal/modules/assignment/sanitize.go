package assignment

import (
	"fmt"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

var assignmentHTMLPolicy = bluemonday.UGCPolicy()

func sanitizeResponseText(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", nil
	}
	sanitized := strings.TrimSpace(assignmentHTMLPolicy.Sanitize(trimmed))
	if sanitized == "" {
		return "", nil
	}
	if len(sanitized) > MaxAssignmentResponseTextLength {
		return "", fmt.Errorf("response text exceeds maximum length of %d characters", MaxAssignmentResponseTextLength)
	}
	return sanitized, nil
}
