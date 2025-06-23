package enrollment

import (
	"dashlearn/middleware"
	"dashlearn/utils"

	"github.com/gin-gonic/gin"
)

func RegisterEnrollmentRoutes(rg *gin.RouterGroup) {

	handler := NewEnrollmentHandler(utils.DB)

	authGroup := rg.Group("/private/enrollment", middleware.AuthMiddleware())
	{
		authGroup.GET("/", handler.GetEnrollments)
		authGroup.POST("/create", handler.CreateEnrollment)
		authGroup.DELETE("/delete/:id", handler.Delete)
	}

	publicGroup := rg.Group("/enrolled", middleware.GetTenantID())
	{
		publicGroup.GET("/courses", middleware.StudentAuthMiddleware(), handler.GetEnrolledCourses)
	}
}
