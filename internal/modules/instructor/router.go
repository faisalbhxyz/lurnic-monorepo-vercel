package instructor

import (
	"dashlearn/internal/middleware"
	"dashlearn/internal/utils"

	"github.com/gin-gonic/gin"
)

func RegisterInstructorRoutes(rg *gin.RouterGroup) {

	handler := NewInstructorHandler(utils.DB)

	authGroup := rg.Group("/private/instructor", middleware.AuthMiddleware())
	{
		authGroup.GET("", handler.GetInstructors)
		authGroup.GET("/lite", handler.GetInstructorsLite)
		authGroup.GET("/details/:id", handler.GetInstructorDetails)
		authGroup.POST("/register", handler.CreateInstructor)
		authGroup.PUT("/update/:id", handler.UpdateInstructor)
		authGroup.DELETE("/delete/:id", handler.DeleteInstructor)
		// userGroup.POST("/login", LoginUser)
		// userGroup.POST("/upload", UploadUser)
	}

	publicRoutes := rg.Group("/instructor", middleware.GetTenantID())
	{
		publicRoutes.GET("/all", handler.GetInstructorsLite)
	}
}
