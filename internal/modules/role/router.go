package role

import (
	"dashlearn/internal/middleware"
	"dashlearn/internal/utils"

	"github.com/gin-gonic/gin"
)

func RegsiterRoleRoutes(rg *gin.RouterGroup) {

	handler := NewRoleHandler(utils.DB)

	roleGroup := rg.Group("/role", middleware.AuthMiddleware())
	{
		roleGroup.GET("/collection", handler.GetRoles)
		roleGroup.GET("/user", handler.GetRoleByUserID)
		roleGroup.POST("/create", handler.CreateRole)
		roleGroup.DELETE("/delete/:role_id", handler.DeleteRoleHandler)
	}

}
