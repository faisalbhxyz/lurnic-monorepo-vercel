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
		studentCourseGroup.GET("/:slug/quizzes/:quizId/questions/:questionIndex", handler.GetStudentQuizQuestion)
		studentCourseGroup.POST("/:slug/quizzes/:quizId/submit", handler.SubmitQuiz)
		studentCourseGroup.POST("/:slug/quizzes/:quizId/skip", handler.SkipQuiz)
	}

	studentGroup := rg.Group("/student", middleware.GetTenantID(), middleware.StudentAuthMiddleware())
	{
		studentGroup.GET("/quiz-submissions", handler.ListStudentSubmissions)
		studentGroup.GET("/quiz-submissions/:submissionId", handler.GetStudentSubmission)
	}

	adminGroup := rg.Group("/private/course", middleware.AuthMiddleware())
	{
		adminGroup.GET("/:id/quiz-submissions", handler.ListCourseSubmissions)
		adminGroup.GET("/:id/quiz-submissions/:submissionId", handler.GetCourseSubmission)
		adminGroup.POST("/:id/quiz-submissions/:submissionId/feedback", handler.UpdateSubmissionFeedback)
	}
}
