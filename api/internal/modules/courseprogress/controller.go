package courseprogress

import (
	"dashlearn/internal/utils"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	service *Service
}

func NewHandler() *Handler {
	return &Handler{service: NewService(utils.DB)}
}

func (h *Handler) GetCourseProgress(c *gin.Context) {
	slug := c.Param("slug")
	tenantID := c.GetUint("tenant_id")
	studentID := c.GetUint("student_id")

	data, err := h.service.GetCourseProgress(tenantID, studentID, slug)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *Handler) MarkLessonComplete(c *gin.Context) {
	slug := c.Param("slug")
	tenantID := c.GetUint("tenant_id")
	studentID := c.GetUint("student_id")

	lessonID, err := strconv.ParseUint(c.Param("lessonId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lesson ID"})
		return
	}

	data, err := h.service.MarkLessonComplete(tenantID, studentID, slug, uint(lessonID))
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data, "message": "Lesson marked complete"})
}

func (h *Handler) GetLessonVideoProgress(c *gin.Context) {
	slug := c.Param("slug")
	tenantID := c.GetUint("tenant_id")
	studentID := c.GetUint("student_id")

	lessonID, err := strconv.ParseUint(c.Param("lessonId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lesson ID"})
		return
	}

	data, err := h.service.GetLessonVideoProgress(tenantID, studentID, slug, uint(lessonID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "progress not found"})
			return
		}
		status := http.StatusBadRequest
		if err.Error() == "course not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *Handler) UpdateLessonVideoProgress(c *gin.Context) {
	slug := c.Param("slug")
	tenantID := c.GetUint("tenant_id")
	studentID := c.GetUint("student_id")

	lessonID, err := strconv.ParseUint(c.Param("lessonId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lesson ID"})
		return
	}

	var req UpdateLessonVideoProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := h.service.UpdateLessonVideoProgress(tenantID, studentID, slug, uint(lessonID), req)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "course not found" || err.Error() == "lesson not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}
