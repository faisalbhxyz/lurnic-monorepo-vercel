package middleware

import (
	"dashlearn/internal/apiresponse"
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
			apiresponse.Unauthorized(c, "Authorization header is missing", "")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			apiresponse.Unauthorized(c, "Invalid token format", "")
			c.Abort()
			return
		}

		if ok, sessionReplaced := authenticateStudent(c, parts[1]); !ok {
			if sessionReplaced {
				apiresponse.Unauthorized(c, "Your account was logged in on another device. Please sign in again.", "SESSION_REPLACED")
			} else if c.Writer.Status() == 0 {
				apiresponse.Unauthorized(c, "Invalid or expired token", "")
			}
			c.Abort()
			return
		}

		c.Next()
	}
}

// authenticateStudent validates a student JWT and sets user_id, tenant_id, session_id on context.
// Returns (ok, sessionReplaced).
func authenticateStudent(c *gin.Context, tokenStr string) (bool, bool) {
	userID, sessionID, err := utils.ParseStudentSessionID(tokenStr)
	if err != nil {
		return false, false
	}

	var user models.Student
	if err := utils.DB.Where("user_id = ?", userID).Select("id", "user_id", "tenant_id", "status").First(&user).Error; err != nil {
		return false, false
	}

	if !user.Status {
		apiresponse.Forbidden(c, "Account is inactive", "")
		return false, false
	}

	var session models.StudentSession
	err = utils.DB.Where("student_id = ? AND session_id = ?", user.ID, sessionID).First(&session).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, true
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong. Please try again.", "error": "INTERNAL_ERROR"})
		return false, false
	}

	if time.Since(session.LastSeenAt) > time.Minute {
		_ = utils.DB.Model(&session).Update("last_seen_at", time.Now()).Error
	}

	c.Set("user_id", user.ID)
	c.Set("tenant_id", user.TenantID)
	c.Set("session_id", sessionID)
	return true, false
}
