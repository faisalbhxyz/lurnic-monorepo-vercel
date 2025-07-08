package student

import (
	"dashlearn/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterStudentRoutes(rg *gin.RouterGroup) {
	authGroup := rg.Group("/private/student", middleware.AuthMiddleware())
	{
		authGroup.GET("/", GetStudents)
		authGroup.GET("/lite", GetStudentLite)
		authGroup.GET("/details/:id", GetStudentDetailsByID)
		authGroup.POST("/register", CreateStudent)
		authGroup.PUT("/update/:id", UpdateStudent)
		authGroup.DELETE("/delete/:id", DeleteStudent)
	}

	publicGroup := rg.Group("/student", middleware.GetTenantID())
	{
		publicGroup.POST("/login", LoginStudent)
		publicGroup.POST("/register", CreateStudentPublic)
		publicGroup.GET("/details", middleware.StudentAuthMiddleware(), GetStudentDetails)
	}
}
