package freelesson

import (
	"dashlearn/internal/middleware"
	"dashlearn/internal/utils"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	handler := NewHandler(utils.DB)

	publicGroup := rg.Group("/free-lessons", middleware.GetTenantID())
	{
		publicGroup.GET("", handler.ListCatalog)
	}

	studentGroup := rg.Group("/student/free-lessons", middleware.GetTenantID(), middleware.StudentAuthMiddleware())
	{
		studentGroup.GET("", handler.ListLibrary)
		studentGroup.POST("", handler.AddToLibrary)
		studentGroup.DELETE("/:lessonId", handler.RemoveFromLibrary)
	}
}
