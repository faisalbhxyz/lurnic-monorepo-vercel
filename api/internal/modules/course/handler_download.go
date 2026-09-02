package course

import (
	"net/http"
	"strconv"

	"dashlearn/internal/apiresponse"

	"github.com/gin-gonic/gin"
)

func (h *CourseHandler) DownloadLesson(c *gin.Context) {
	slug := c.Param("slug")
	tenantID := c.GetUint("tenant_id")
	studentID := c.GetUint("user_id")

	lessonID, err := strconv.ParseUint(c.Param("lessonId"), 10, 64)
	if err != nil || lessonID == 0 {
		apiresponse.Error(c, http.StatusNotFound, "Lesson not found", "LESSON_NOT_FOUND")
		return
	}

	result, err := h.downloadService.ResolveLessonDownload(tenantID, studentID, slug, uint(lessonID))
	if err != nil {
		switch err {
		case ErrDownloadCourseNotFound, ErrDownloadLessonNotFound:
			apiresponse.Error(c, http.StatusNotFound, "Lesson not found", "LESSON_NOT_FOUND")
		case ErrDownloadNotEnrolled:
			apiresponse.Forbidden(c, "You must be enrolled in this course", "NOT_ENROLLED")
		case ErrDownloadNotDownloadable:
			apiresponse.Error(c, http.StatusUnprocessableEntity, "This lesson cannot be downloaded offline", "NOT_DOWNLOADABLE")
		default:
			apiresponse.Error(c, http.StatusInternalServerError, "Failed to prepare download", "INTERNAL_ERROR")
		}
		return
	}

	if c.Query("format") == "json" {
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"download_url": result.RedirectURL,
				"file_name":    result.FileName,
				"content_type": result.ContentType,
			},
		})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=\""+result.FileName+"\"")
	c.Redirect(http.StatusFound, result.RedirectURL)
}
