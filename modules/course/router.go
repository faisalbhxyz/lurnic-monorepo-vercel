package course

import (
	"dashlearn/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterCourseRoutes(rg *gin.RouterGroup) {
	routerGroup := rg.Group("/course")
	{
		routerGroup.GET("/", middleware.AuthMiddleware(), GetCourses)
		routerGroup.POST("/create", middleware.AuthMiddleware(), CreateCourse)
	}
}
