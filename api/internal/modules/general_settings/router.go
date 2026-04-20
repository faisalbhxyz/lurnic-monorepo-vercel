package generalsettings

import (
	"dashlearn/internal/middleware"
	"dashlearn/internal/utils"

	"github.com/gin-gonic/gin"
)

func RegisterGeneralSettingsRoutes(rg *gin.RouterGroup) {

	handler := NewGeneralSettingsHandler(utils.DB)

	authGroup := rg.Group("/private/general-settings", middleware.AuthMiddleware())
	{
		authGroup.GET("", handler.Get)
		authGroup.PUT("/update", handler.UpdateOrCreate)
	}
}
