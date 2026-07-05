package academicnote

import (
	"dashlearn/internal/middleware"
	"dashlearn/internal/utils"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	handler := NewHandler(utils.DB)

	authGroup := rg.Group("/private/academic-notes", middleware.AuthMiddleware())
	{
		authGroup.GET("/classes", handler.GetAllClasses)
		authGroup.GET("/classes/:id", handler.GetClassByID)
		authGroup.POST("/classes/create", handler.CreateClass)
		authGroup.PUT("/classes/update/:id", handler.UpdateClass)
		authGroup.DELETE("/classes/delete/:id", handler.DeleteClass)

		authGroup.POST("/subjects/create", handler.CreateSubject)
		authGroup.PUT("/subjects/update/:id", handler.UpdateSubject)
		authGroup.DELETE("/subjects/delete/:id", handler.DeleteSubject)

		authGroup.POST("/papers/create", handler.CreatePaper)
		authGroup.PUT("/papers/update/:id", handler.UpdatePaper)
		authGroup.DELETE("/papers/delete/:id", handler.DeletePaper)

		authGroup.POST("/notes/create", handler.CreateNote)
		authGroup.PUT("/notes/update/:id", handler.UpdateNote)
		authGroup.DELETE("/notes/delete/:id", handler.DeleteNote)
	}

	publicGroup := rg.Group("/academic-notes", middleware.GetTenantID())
	{
		publicGroup.GET("", handler.GetPublicClasses)
		publicGroup.GET("/:classSlug", handler.GetPublicClassBySlug)
		publicGroup.GET("/:classSlug/:subjectSlug/:paperSlug", handler.GetPublicNotesByPaperSlug)
	}
}
