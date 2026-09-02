package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// OptionalStudentAuthMiddleware parses Bearer token when present and sets user_id/tenant_id.
// Missing or invalid tokens do not abort the request — used for public endpoints with optional auth.
func OptionalStudentAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		_, _ = authenticateStudent(c, parts[1])
		c.Next()
	}
}
