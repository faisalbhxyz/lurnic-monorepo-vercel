package course

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestDownloadLesson_InvalidLessonID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := &CourseHandler{
		downloadService: NewDownloadService(nil),
	}

	router := gin.New()
	router.GET("/course/:slug/lessons/:lessonId/download", func(c *gin.Context) {
		c.Set("tenant_id", uint(1))
		c.Set("user_id", uint(10))
		handler.DownloadLesson(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/course/test-course/lessons/not-a-number/download", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for invalid lesson id, got %d body=%s", rec.Code, rec.Body.String())
	}
}
