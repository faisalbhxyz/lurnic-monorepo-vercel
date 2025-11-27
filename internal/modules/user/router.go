package user

import (
	"dashlearn/internal/middleware"
	"dashlearn/internal/utils"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(rg *gin.RouterGroup) {
	// userGroup := rg.Group("/user")
	// {
	// 	userGroup.GET("/", middleware.AuthMiddleware(), GetUsers)
	// 	userGroup.GET("/cehck", CheckUser)
	// 	userGroup.POST("/register", CreateUser)
	// 	userGroup.POST("/login", LoginUser)
	// 	userGroup.POST("/upload", UploadUser)
	// }

	handler := NewUserHandler(utils.DB)

	userGroup := rg.Group("/user")
	{
		// userGroup.GET("/", middleware.AuthMiddleware(), GetUsers)
		userGroup.POST("/register", handler.CreateUser)
		userGroup.POST("/login", handler.LoginUser)
	}
	teamMemberGroup := rg.Group("/team-member", middleware.AuthMiddleware())
	{
		teamMemberGroup.GET("/collection", handler.GetTeamMembers)
		teamMemberGroup.POST("/create", handler.CreateTeamMember)
		teamMemberGroup.GET("/details/:id", handler.GetTeamMemberByUID)
		teamMemberGroup.PUT("/update/:id", handler.UpdateTeamMember)
		teamMemberGroup.DELETE("/delete/:id", handler.DeleteTeamMember)
	}

}
