package user

import (
	"dashlearn/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(rg *gin.RouterGroup) {
	userGroup := rg.Group("/user")
	{
		userGroup.GET("/", middleware.AuthMiddleware(), GetUsers)
		userGroup.GET("/cehck", CheckUser)
		userGroup.POST("/register", CreateUser)
		userGroup.POST("/login", LoginUser)
		userGroup.POST("/upload", UploadUser)
	}
}
