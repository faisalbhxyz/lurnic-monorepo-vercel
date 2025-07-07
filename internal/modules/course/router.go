package course

import (
	"dashlearn/middleware"
	"dashlearn/utils"

	"github.com/gin-gonic/gin"
)

func RegisterCourseRoutes(rg *gin.RouterGroup) {

	handler := NewCourseHandler(utils.DB)

	authGroup := rg.Group("/private/course", middleware.AuthMiddleware())
	{
		authGroup.GET("/", handler.GetAll)
		authGroup.GET("/lite", handler.GetAllLite)
		authGroup.POST("/create", handler.Create)
		authGroup.GET("/:id", handler.GetByID)
		authGroup.PUT("/update/:id", handler.Update)
		authGroup.DELETE("/delete/:id", handler.Delete)
	}

	publicGroup := rg.Group("/course", middleware.GetTenantID())
	{
		publicGroup.GET("/", handler.GetAllPublic)
		publicGroup.GET("/:id", handler.GetByID)
	}
}
