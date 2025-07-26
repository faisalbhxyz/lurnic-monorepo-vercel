package course

import (
	"dashlearn/internal/middleware"
	"dashlearn/internal/utils"

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
		publicGroup.GET("/:slug", handler.GetBySlugPublic)
		publicGroup.GET("/search", handler.GetSearchCourses)
		publicGroup.GET("/category/:category", handler.GetAllPublicByCategory)
		publicGroup.GET("/menu/:subcategory", handler.GetAllPublicBySubCategory)
	}
}
