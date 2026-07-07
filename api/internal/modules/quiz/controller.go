package quiz

import (
	"net/http"
	"strconv"

	"dashlearn/internal/utils"

	"github.com/gin-gonic/gin"
)

type QuizHandler struct {
	service QuizService
}

func NewQuizHandler() *QuizHandler {
	return &QuizHandler{service: NewQuizService(utils.DB)}
}

func (h *QuizHandler) GetStudentQuiz(c *gin.Context) {
	quizID, err := strconv.ParseUint(c.Param("quizId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quiz ID"})
		return
	}

	data, err := h.service.GetStudentQuiz(
		c.GetUint("tenant_id"),
		c.GetUint("user_id"),
		c.Param("slug"),
		uint(quizID),
	)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "course not found" || err.Error() == "quiz not found" {
			status = http.StatusNotFound
		} else if err.Error() == "enrollment required" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *QuizHandler) GetStudentQuizQuestion(c *gin.Context) {
	quizID, err := strconv.ParseUint(c.Param("quizId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quiz ID"})
		return
	}
	questionIndex, err := strconv.Atoi(c.Param("questionIndex"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question index"})
		return
	}

	data, err := h.service.GetStudentQuizQuestion(
		c.GetUint("tenant_id"),
		c.GetUint("user_id"),
		c.Param("slug"),
		uint(quizID),
		questionIndex,
	)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "course not found" || err.Error() == "quiz not found" {
			status = http.StatusNotFound
		} else if err.Error() == "enrollment required" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *QuizHandler) SubmitQuiz(c *gin.Context) {
	quizID, err := strconv.ParseUint(c.Param("quizId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quiz ID"})
		return
	}

	var input SubmitQuizInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := h.service.SubmitQuiz(
		c.GetUint("tenant_id"),
		c.GetUint("user_id"),
		c.Param("slug"),
		uint(quizID),
		input,
	)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "course not found" || err.Error() == "quiz not found" {
			status = http.StatusNotFound
		} else if err.Error() == "enrollment required" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Quiz submitted successfully",
		"data":    data,
	})
}

func (h *QuizHandler) SkipQuiz(c *gin.Context) {
	quizID, err := strconv.ParseUint(c.Param("quizId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quiz ID"})
		return
	}

	data, err := h.service.SkipQuiz(
		c.GetUint("tenant_id"),
		c.GetUint("user_id"),
		c.Param("slug"),
		uint(quizID),
	)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "course not found" || err.Error() == "quiz not found" {
			status = http.StatusNotFound
		} else if err.Error() == "enrollment required" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Quiz skipped successfully",
		"data":    data,
	})
}

func (h *QuizHandler) ListStudentSubmissions(c *gin.Context) {
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

func (h *QuizHandler) GetStudentSubmission(c *gin.Context) {
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

func (h *QuizHandler) ListCourseSubmissions(c *gin.Context) {
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

func (h *QuizHandler) GetCourseSubmission(c *gin.Context) {
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

func (h *QuizHandler) UpdateSubmissionFeedback(c *gin.Context) {
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

	var input UpdateQuizSubmissionFeedbackInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := h.service.UpdateSubmissionFeedback(
		c.GetUint("tenant_id"),
		uint(courseID),
		uint(submissionID),
		input,
	)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "submission not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Instructor feedback updated successfully",
		"data":    data,
	})
}
