package apiresponse

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Error writes a storefront-compatible error body: { "message": "...", "error": "CODE" }.
func Error(c *gin.Context, status int, message string, code string) {
	body := gin.H{"message": message}
	if code != "" {
		body["error"] = code
	}
	c.JSON(status, body)
}

// Validation writes a 422 validation error.
func Validation(c *gin.Context, message string) {
	Error(c, http.StatusUnprocessableEntity, message, "VALIDATION_ERROR")
}

// Unauthorized writes a 401 with optional machine code.
func Unauthorized(c *gin.Context, message, code string) {
	Error(c, http.StatusUnauthorized, message, code)
}

// Forbidden writes a 403 with optional machine code.
func Forbidden(c *gin.Context, message, code string) {
	Error(c, http.StatusForbidden, message, code)
}

// NotFound writes a 404.
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message, "NOT_FOUND")
}

// Conflict writes a 409.
func Conflict(c *gin.Context, message, code string) {
	Error(c, http.StatusConflict, message, code)
}
