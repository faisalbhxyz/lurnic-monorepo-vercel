package certificate

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

func (h *Handler) ListStudentCertificates(c *gin.Context) {
	tenantID := c.GetUint("tenant_id")
	studentID := c.GetUint("student_id")

	items, err := h.service.ListForStudent(tenantID, studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *Handler) GetStudentCertificate(c *gin.Context) {
	tenantID := c.GetUint("tenant_id")
	studentID := c.GetUint("student_id")

	certificateID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid certificate ID"})
		return
	}

	item, err := h.service.GetForStudent(tenantID, studentID, uint(certificateID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Certificate not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *Handler) GetCourseCertificate(c *gin.Context) {
	tenantID := c.GetUint("tenant_id")
	studentID := c.GetUint("student_id")
	slug := c.Param("slug")

	item, err := h.service.GetForCourseSlug(tenantID, studentID, slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Certificate not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}
