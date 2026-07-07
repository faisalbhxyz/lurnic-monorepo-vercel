package assignment

import (
	"dashlearn/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAssignmentRoutes(rg *gin.RouterGroup) {
	handler := NewAssignmentHandler()

	studentCourseGroup := rg.Group("/course", middleware.GetTenantID(), middleware.StudentAuthMiddleware())
	{
		studentCourseGroup.GET("/:slug/assignments/:assignmentId", handler.GetStudentAssignment)
		studentCourseGroup.POST("/:slug/assignments/:assignmentId/submit", handler.SubmitAssignment)
	}

	studentGroup := rg.Group("/student", middleware.GetTenantID(), middleware.StudentAuthMiddleware())
	{
		studentGroup.GET("/assignment-submissions", handler.ListStudentSubmissions)
		studentGroup.GET("/assignment-submissions/:submissionId", handler.GetStudentSubmission)
	}

	adminGroup := rg.Group("/private/course", middleware.AuthMiddleware())
	{
		adminGroup.GET("/:id/assignment-submissions", handler.ListCourseSubmissions)
		adminGroup.GET("/:id/assignment-submissions/:submissionId", handler.GetCourseSubmission)
		adminGroup.POST("/:id/assignment-submissions/:submissionId/grade", handler.GradeSubmission)
	}
}
