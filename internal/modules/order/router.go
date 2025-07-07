package order

import (
	"dashlearn/middleware"
	"dashlearn/utils"

	"github.com/gin-gonic/gin"
)

func RegisterCourseRoutes(rg *gin.RouterGroup) {

	handler := NewOrderHandler(utils.DB)

	authGroup := rg.Group("/private/order", middleware.AuthMiddleware())
	{
		authGroup.GET("/", handler.GetAll)
		authGroup.POST("/create", handler.Create)
		authGroup.PUT("/mark-as-paid/:id", handler.MarkAsPaid)
		authGroup.DELETE("/delete/:id", handler.Delete)
	}

	publicGroup := rg.Group("/order", middleware.StudentAuthMiddleware())
	{
		publicGroup.POST("/create", handler.Create)
		// publicGroup.GET("/:id", handler.GetByID)
	}
}
