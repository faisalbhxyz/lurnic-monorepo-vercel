package category

import (
	"dashlearn/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterCategoryRoutes(rg *gin.RouterGroup) {
	authgroup := rg.Group("/private/category", middleware.AuthMiddleware())
	{
		authgroup.GET("/", GetAllCategory)
		authgroup.GET("/:id", GetCategoryByID)
		authgroup.POST("/create", CreateCategory)
		authgroup.PUT("/update/:id", UpdateCategory)
		// routerGroup.DELETE("/delete/:id", middleware.AuthMiddleware(), DeleteCategory)
	}

	publicGroup := rg.Group("/category", middleware.GetTenantID())
	{
		publicGroup.GET("/", GetAllCategoryPublic)
		// publicGroup.GET("/:id", GetCategoryByIDPublic)
	}
}
