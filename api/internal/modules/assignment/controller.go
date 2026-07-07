package assignment

import (
	"mime/multipart"
	"net/http"
	"strconv"

	"dashlearn/internal/utils"

	"github.com/gin-gonic/gin"
)

type AssignmentHandler struct {
	service AssignmentService
}

func NewAssignmentHandler() *AssignmentHandler {
	return &AssignmentHandler{
		service: NewAssignmentService(utils.DB, utils.UploadToBunny),
	}
}

func (h *AssignmentHandler) GetStudentAssignment(c *gin.Context) {
	assignmentID, err := strconv.ParseUint(c.Param("assignmentId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
		return
	}

	data, err := h.service.GetStudentAssignment(
		c.GetUint("tenant_id"),
		c.GetUint("user_id"),
		c.Param("slug"),
		uint(assignmentID),
	)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "course not found" || err.Error() == "assignment not found" {
			status = http.StatusNotFound
		} else if err.Error() == "enrollment required" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *AssignmentHandler) SubmitAssignment(c *gin.Context) {
	assignmentID, err := strconv.ParseUint(c.Param("assignmentId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
		return
	}

	var responseText *string
	if raw := c.PostForm("response_text"); raw != "" {
		responseText = &raw
	}

	form, err := c.MultipartForm()
	if err != nil && err != http.ErrNotMultipart {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var files []*multipart.FileHeader
	if form != nil {
		files = form.File["files"]
	}

	data, err := h.service.SubmitAssignment(
		c.GetUint("tenant_id"),
		c.GetUint("user_id"),
		c.Param("slug"),
		uint(assignmentID),
		responseText,
		files,
	)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "course not found" || err.Error() == "assignment not found" {
			status = http.StatusNotFound
		} else if err.Error() == "enrollment required" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Assignment submitted successfully",
		"data":    data,
	})
}

func (h *AssignmentHandler) ListStudentSubmissions(c *gin.Context) {
	var courseID *uint
	if raw := c.Query("course_id"); raw != "" {
		parsed, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course_id"})
			return
		}
		v := uint(parsed)
		courseID = &v
	}

	data, err := h.service.ListStudentSubmissions(c.GetUint("tenant_id"), c.GetUint("user_id"), courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *AssignmentHandler) GetStudentSubmission(c *gin.Context) {
	submissionID, err := strconv.ParseUint(c.Param("submissionId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid submission ID"})
		return
	}

	data, err := h.service.GetStudentSubmission(c.GetUint("tenant_id"), c.GetUint("user_id"), uint(submissionID))
	if err != nil {
		status := http.StatusNotFound
		if err.Error() != "submission not found" {
			status = http.StatusInternalServerError
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *AssignmentHandler) ListCourseSubmissions(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	data, err := h.service.ListCourseSubmissions(c.GetUint("tenant_id"), uint(courseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *AssignmentHandler) GetCourseSubmission(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}
	submissionID, err := strconv.ParseUint(c.Param("submissionId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid submission ID"})
		return
	}

	data, err := h.service.GetCourseSubmission(c.GetUint("tenant_id"), uint(courseID), uint(submissionID))
	if err != nil {
		status := http.StatusNotFound
		if err.Error() != "submission not found" {
			status = http.StatusInternalServerError
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *AssignmentHandler) GradeSubmission(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}
	submissionID, err := strconv.ParseUint(c.Param("submissionId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid submission ID"})
		return
	}

	var input GradeAssignmentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := h.service.GradeSubmission(c.GetUint("tenant_id"), uint(courseID), uint(submissionID), input)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "submission not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Assignment graded successfully",
		"data":    data,
	})
}
