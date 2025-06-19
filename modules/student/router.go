package student

import (
	"dashlearn/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterStudentRoutes(rg *gin.RouterGroup) {
	authGroup := rg.Group("/private/student")
	{
		authGroup.GET("/", middleware.AuthMiddleware(), GetStudents)
		authGroup.GET("/lite", middleware.AuthMiddleware(), GetStudentLite)
		authGroup.POST("/register", middleware.AuthMiddleware(), CreateStudent)
		// userGroup.POST("/upload", UploadUser)
	}

	publicGroup := rg.Group("/student")
	{
		publicGroup.POST("/login", LoginStudent)
	}
}
