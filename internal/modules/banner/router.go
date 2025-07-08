package banner

import (
	"dashlearn/internal/middleware"
	"dashlearn/internal/utils"

	"github.com/gin-gonic/gin"
)

func RegisterBannerRoutes(rg *gin.RouterGroup) {

	handler := NewBannerHandler(utils.DB)

	authgroup := rg.Group("/private/banner", middleware.AuthMiddleware())
	{
		authgroup.GET("/", handler.GetAll)
		authgroup.GET("/:id", handler.GetByID)
		authgroup.POST("/create", handler.Create)
		authgroup.PUT("/update/:id", handler.Update)
		authgroup.DELETE("/delete/:id", handler.Delete)
	}

	publicGroup := rg.Group("/banners", middleware.GetTenantID())
	{
		publicGroup.GET("/", handler.GetAll)
	}
}
