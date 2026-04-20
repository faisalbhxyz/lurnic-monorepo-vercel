package paymentmethod

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func handleValidationError(c *gin.Context, err error) {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		errorsMap := make(map[string]string)
		for _, fieldErr := range validationErrors {
			switch fieldErr.Field() {
			case "Title":
				if fieldErr.Tag() == "required" {
					errorsMap["title"] = "Title is required"
				}
			case "Instruction":
				if fieldErr.Tag() == "required" {
					errorsMap["instruction"] = "Instruction is required"
				}
			}
		}
		c.JSON(http.StatusBadRequest, gin.H{"errors": errorsMap})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}
