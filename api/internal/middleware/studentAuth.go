package middleware

import (
	"dashlearn/internal/models"
	"dashlearn/internal/utils"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func StudentAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		tokenStr := parts[1]
		userID, sessionID, err := utils.ParseStudentSessionID(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		var user models.Student
		if err := utils.DB.Where("user_id = ?", userID).Select("id", "user_id", "tenant_id", "status").First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		if !user.Status {
			c.JSON(http.StatusForbidden, gin.H{"error": "Account is inactive"})
			c.Abort()
			return
		}

		var session models.StudentSession
		err = utils.DB.Where("student_id = ? AND session_id = ?", user.ID, sessionID).First(&session).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Session expired or logged in on another device",
				"code":    "SESSION_REPLACED",
				"message": "Your account was logged in on another device. Please sign in again.",
			})
			c.Abort()
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong. Please try again."})
			c.Abort()
			return
		}

		if time.Since(session.LastSeenAt) > time.Minute {
			_ = utils.DB.Model(&session).Update("last_seen_at", time.Now()).Error
		}

		c.Set("user_id", user.ID)
		c.Set("tenant_id", user.TenantID)
		c.Set("session_id", sessionID)
		c.Next()
	}
}
