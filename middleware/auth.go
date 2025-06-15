package middleware

import (
	"dashlearn/models"
	"dashlearn/utils"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var SecretKey = []byte(os.Getenv("JWT_SECRET"))

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			ctx.Abort()
			return
		}

		// Expected: Bearer <token>
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			ctx.Abort()
			return
		}

		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			// Validate signing method
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return SecretKey, nil
		})

		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			ctx.Abort()
			return
		}

		// You can store claims/user info in context if needed
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if exp, ok := claims["exp"].(float64); ok && int64(exp) < time.Now().Unix() {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
				ctx.Abort()
				return
			}

			// get tenantID
			var user models.User
			utils.DB.Where("user_id = ?", claims["user_id"]).Select("id", "user_id", "tenant_id").First(&user)

			// Set user info in context
			ctx.Set("user_id", user.ID)
			ctx.Set("tenant_id", user.TenantID)
		}

		ctx.Next()
	}
}
