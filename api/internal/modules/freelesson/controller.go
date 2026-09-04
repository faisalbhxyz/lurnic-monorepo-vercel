package freelesson

import (
	"dashlearn/internal/utils"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	service *Service
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{service: NewService(db)}
}

func (h *Handler) ListCatalog(c *gin.Context) {
	tenantID := c.GetUint("tenant_id")
	classSlug := strings.TrimSpace(c.Query("class_slug"))

	limit := 50
	if raw := c.Query("limit"); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil {
			limit = n
		}
	}
	offset := 0
	if raw := c.Query("offset"); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil {
			offset = n
		}
	}

	data, meta, err := h.service.ListCatalog(tenantID, classSlug, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data, "meta": meta})
}

func (h *Handler) ListLibrary(c *gin.Context) {
	tenantID := c.GetUint("tenant_id")
	studentID := c.GetUint("user_id")

	data, err := h.service.ListLibrary(tenantID, studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *Handler) AddToLibrary(c *gin.Context) {
	tenantID := c.GetUint("tenant_id")
	studentID := c.GetUint("user_id")

	var req AddFreeLessonsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := h.service.AddToLibrary(tenantID, studentID, req.LessonIDs)
	if err != nil {
		status := http.StatusUnprocessableEntity
		msg := err.Error()
		if strings.Contains(msg, "required") {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data, "message": "Free lessons saved"})
}

func (h *Handler) RemoveFromLibrary(c *gin.Context) {
	tenantID := c.GetUint("tenant_id")
	studentID := c.GetUint("user_id")

	lessonID, err := strconv.ParseUint(c.Param("lessonId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lesson ID"})
		return
	}

	if err := h.service.RemoveFromLibrary(tenantID, studentID, uint(lessonID)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Free lesson not found in library"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Free lesson removed"})
}

// NewHandlerDefault uses the shared DB connection.
func NewHandlerDefault() *Handler {
	return NewHandler(utils.DB)
}
