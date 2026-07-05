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
	}

	adminGroup := rg.Group("/private/course", middleware.AuthMiddleware())
	{
		adminGroup.GET("/:courseId/assignment-submissions", handler.ListCourseSubmissions)
		adminGroup.GET("/:courseId/assignment-submissions/:submissionId", handler.GetCourseSubmission)
		adminGroup.POST("/:courseId/assignment-submissions/:submissionId/grade", handler.GradeSubmission)
	}
}
