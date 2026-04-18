package category

import (
	"dashlearn/internal/middleware"
	"dashlearn/internal/utils"

	"github.com/gin-gonic/gin"
)

func RegisterCategoryRoutes(rg *gin.RouterGroup) {

	handler := NewCategoryHandler(utils.DB)

	authgroup := rg.Group("/private/category", middleware.AuthMiddleware())
	{
		authgroup.GET("", handler.GetAll)
		authgroup.GET("/:id", handler.GetByID)
		authgroup.POST("/create", handler.Create)
		authgroup.PUT("/update/:id", handler.Update)
		authgroup.DELETE("/delete/:id", handler.Delete)
	}

	publicGroup := rg.Group("/category", middleware.GetTenantID())
	{
		publicGroup.GET("", handler.GetAll)
		// publicGroup.GET("/:id", GetCategoryByIDPublic)
	}
}
