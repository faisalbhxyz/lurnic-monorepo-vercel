package coursereview

import (
	"log"
	"net/http"
	"strconv"

	"dashlearn/internal/apiresponse"
	"dashlearn/internal/utils"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler() *Handler {
	return &Handler{service: NewService(utils.DB)}
}

func (h *Handler) ListCourseReviews(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	studentID := c.GetUint("user_id")
	authed := studentID > 0

	data, err := h.service.ListReviews(ListReviewsParams{
		TenantID:  c.GetUint("tenant_id"),
		Slug:      c.Param("slug"),
		StudentID: studentID,
		Authed:    authed,
		Page:      page,
		PerPage:   perPage,
	})
	if err != nil {
		if err == errCourseNotFound {
			apiresponse.NotFound(c, "Course not found")
			return
		}
		log.Printf("[coursereview] ListCourseReviews failed slug=%s tenant=%d: %v", c.Param("slug"), c.GetUint("tenant_id"), err)
		apiresponse.Error(c, http.StatusInternalServerError, "Failed to load reviews", "INTERNAL_ERROR")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *Handler) CreateCourseReview(c *gin.Context) {
	var input CreateReviewInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apiresponse.Validation(c, "Invalid review payload")
		return
	}

	data, err := h.service.CreateReview(
		c.GetUint("tenant_id"),
		c.GetUint("user_id"),
		c.Param("slug"),
		input,
	)
	if err != nil {
		switch err {
		case errCourseNotFound:
			apiresponse.NotFound(c, "Course not found")
		case errNotEnrolled:
			apiresponse.Forbidden(c, "You must be enrolled to review this course", "NOT_ENROLLED")
		case errCourseNotCompleted:
			apiresponse.Forbidden(c, "Complete the course before leaving a review", "COURSE_NOT_COMPLETED")
		case errReviewExists:
			apiresponse.Conflict(c, "You have already reviewed this course", "REVIEW_ALREADY_EXISTS")
		case errValidation:
			apiresponse.Validation(c, "Invalid review data")
		default:
			apiresponse.Error(c, http.StatusInternalServerError, "Failed to submit review", "INTERNAL_ERROR")
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Review submitted successfully",
		"data":    data,
	})
}
