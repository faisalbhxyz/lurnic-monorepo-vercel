package instructor

import (
	"dashlearn/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterInstructorRoutes(rg *gin.RouterGroup) {
	authGroup := rg.Group("/private/instructor", middleware.AuthMiddleware())
	{
		authGroup.GET("/", GetInstructors)
		authGroup.GET("/lite", GetInstructorsLite)
		authGroup.GET("/details/:id", GetInstructorDetails)
		authGroup.POST("/register", CreateInstructor)
		authGroup.PUT("/update/:id", UpdateInstructor)
		authGroup.DELETE("/delete/:id", DeleteInstructor)
		// userGroup.POST("/login", LoginUser)
		// userGroup.POST("/upload", UploadUser)
	}
}
