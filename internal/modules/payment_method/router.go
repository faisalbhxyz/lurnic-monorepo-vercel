package paymentmethod

import (
	"dashlearn/internal/middleware"
	"dashlearn/internal/utils"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup) {

	handler := NewPaymentMethodHandler(utils.DB)

	authgroup := rg.Group("/private/payment-method", middleware.AuthMiddleware())
	{
		authgroup.GET("/", handler.GetAll)
		authgroup.GET("/:id", handler.GetByID)
		authgroup.POST("/create", handler.Create)
		authgroup.PUT("/update/:id", handler.Update)
		authgroup.DELETE("/delete/:id", handler.Delete)
	}

	publicGroup := rg.Group("/payment-methods", middleware.GetTenantID())
	{
		publicGroup.GET("/", handler.GetAll)
		publicGroup.GET("/:id", handler.GetByID)
	}
}
