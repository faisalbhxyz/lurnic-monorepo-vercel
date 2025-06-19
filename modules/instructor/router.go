package instructor

import (
	"dashlearn/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterInstructorRoutes(rg *gin.RouterGroup) {
	authGroup := rg.Group("/private/instructor")
	{
		authGroup.GET("/", middleware.AuthMiddleware(), GetInstructors)
		authGroup.GET("/lite", middleware.AuthMiddleware(), GetInstructorsLite)
		authGroup.POST("/register", middleware.AuthMiddleware(), CreateInstructor)
		// userGroup.POST("/login", LoginUser)
		// userGroup.POST("/upload", UploadUser)
	}
}
