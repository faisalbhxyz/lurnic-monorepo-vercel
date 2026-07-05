package certificate

import (
	"dashlearn/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterCertificateRoutes(rg *gin.RouterGroup) {
	handler := NewHandler()

	studentGroup := rg.Group("/student", middleware.GetTenantID(), middleware.StudentAuthMiddleware())
	{
		studentGroup.GET("/certificates", handler.ListStudentCertificates)
		studentGroup.GET("/certificates/:id", handler.GetStudentCertificate)
	}

	studentCourseGroup := rg.Group("/course", middleware.GetTenantID(), middleware.StudentAuthMiddleware())
	{
		studentCourseGroup.GET("/:slug/certificate", handler.GetCourseCertificate)
	}
}
