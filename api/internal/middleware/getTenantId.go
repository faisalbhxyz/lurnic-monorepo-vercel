package middleware

import (
	"dashlearn/internal/models"
	"dashlearn/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTenantID() gin.HandlerFunc {
	return func(c *gin.Context) {
		appKeyHeader := c.GetHeader("app-key")

		if appKeyHeader == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "App key (app-key) header is missing"})
			c.Abort()
			return
		}

		// get tenantID
		var tenant models.Tenant
		utils.DB.Where("app_key = ?", appKeyHeader).Select("id", "app_key").First(&tenant)

		if tenant.ID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tenant not found"})
			c.Abort()
			return
		}

		c.Set("tenant_id", tenant.ID)

		c.Next()
	}
}
