package instructor

import (
	"dashlearn/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterInstructorRoutes(rg *gin.RouterGroup) {
	routesGroup := rg.Group("/instructor")
	{
		// userGroup.GET("/", middleware.AuthMiddleware(), GetUsers)
		routesGroup.POST("/register", middleware.AuthMiddleware(), CreateInstructor)
		// userGroup.POST("/login", LoginUser)
		// userGroup.POST("/upload", UploadUser)
	}
}
