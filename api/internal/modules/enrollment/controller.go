package enrollment

import (
	"dashlearn/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EnrollmentHandler struct {
	service EnrollmentService
}

func NewEnrollmentHandler(db *gorm.DB) *EnrollmentHandler {
	return &EnrollmentHandler{
		service: NewEnrollmentService(db),
	}
}

func (h *EnrollmentHandler) CreateEnrollment(c *gin.Context) {
	var input models.Enrollment
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Create(input, c.GetUint("tenant_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Enrollment created successfully"})
}

func (h *EnrollmentHandler) GetEnrollments(c *gin.Context) {

	enrollments, err := h.service.GetAll(c.GetUint("tenant_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": enrollments})
}

func (h *EnrollmentHandler) GetEnrolledCourses(c *gin.Context) {
	enrollments, err := h.service.GetEnrolledCourses(c.GetUint("tenant_id"), c.GetUint("user_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": enrollments})
}

func (h *EnrollmentHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid enrollment ID"})
		return
	}

	if err := h.service.Delete(uint(id), c.GetUint("tenant_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "Enrollment deleted successfully"})
}
