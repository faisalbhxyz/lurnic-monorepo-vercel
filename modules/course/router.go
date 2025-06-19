package course

import (
	"dashlearn/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterCourseRoutes(rg *gin.RouterGroup) {

	authGroup := rg.Group("/private/course", middleware.AuthMiddleware())
	{
		authGroup.GET("/", GetCourses)
		authGroup.GET("/lite", GetCoursesLite)
		authGroup.POST("/create", CreateCourse)
	}

	publicGroup := rg.Group("/course", middleware.GetTenantID())
	{
		publicGroup.GET("/", GetPublicCourses)
		publicGroup.GET("/:id", GetCourseByID)
	}
}
