package quiz

import (
	"dashlearn/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterQuizRoutes(rg *gin.RouterGroup) {
	handler := NewQuizHandler()

	studentCourseGroup := rg.Group("/course", middleware.GetTenantID(), middleware.StudentAuthMiddleware())
	{
		studentCourseGroup.GET("/:slug/quizzes/:quizId", handler.GetStudentQuiz)
		studentCourseGroup.POST("/:slug/quizzes/:quizId/submit", handler.SubmitQuiz)
	}

	studentGroup := rg.Group("/student", middleware.GetTenantID(), middleware.StudentAuthMiddleware())
	{
		studentGroup.GET("/quiz-submissions", handler.ListStudentSubmissions)
	}

	adminGroup := rg.Group("/private/course", middleware.AuthMiddleware())
	{
		adminGroup.GET("/:courseId/quiz-submissions", handler.ListCourseSubmissions)
		adminGroup.GET("/:courseId/quiz-submissions/:submissionId", handler.GetCourseSubmission)
	}
}
