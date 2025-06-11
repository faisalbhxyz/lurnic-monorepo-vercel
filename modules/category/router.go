package category

import (
	"dashlearn/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterCategoryRoutes(rg *gin.RouterGroup) {
	routerGroup := rg.Group("/category")
	{
		routerGroup.GET("/", middleware.AuthMiddleware(), GetAllCategory)
		routerGroup.GET("/:id", middleware.AuthMiddleware(), GetCategoryByID)
		routerGroup.POST("/create", middleware.AuthMiddleware(), CreateCategory)
		routerGroup.PUT("/update/:id", middleware.AuthMiddleware(), UpdateCategory)
		// routerGroup.DELETE("/delete/:id", middleware.AuthMiddleware(), DeleteCategory)
	}
}
