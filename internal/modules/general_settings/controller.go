package generalsettings

import (
	"context"
	"dashlearn/internal/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type GeneralSettingsHandler struct {
	service GeneralSettingsService
}

func NewGeneralSettingsHandler(db *gorm.DB) *GeneralSettingsHandler {
	return &GeneralSettingsHandler{
		service: NewGeneralSettingsService(db),
	}
}

func (h *GeneralSettingsHandler) Get(c *gin.Context) {
	tenantID := c.GetUint("tenant_id")
	res, err := h.service.GetGeneralSettings(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h *GeneralSettingsHandler) UpdateOrCreate(c *gin.Context) {
	var input CreateOrUpdateGeneralSettingsInput
	if err := c.ShouldBindWith(&input, binding.FormMultipart); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			errorsMap := make(map[string]string)
			for _, fieldErr := range validationErrors {
				switch fieldErr.Field() {
				case "OrgName":
					switch fieldErr.Tag() {
					case "required":
						errorsMap["org_name"] = "Name is required"
					case "min":
						errorsMap["org_name"] = "Name must be at least 3 characters long"
					case "max":
						errorsMap["org_name"] = "Name must be at most 100 characters long"
					}
				case "StudentPrefix":
					switch fieldErr.Tag() {
					case "required":
						errorsMap["student_prefix"] = "Student Prefix is required"
					case "min":
						errorsMap["student_prefix"] = "Student prefix must be at least 1 characters long"
					case "max":
						errorsMap["student_prefix"] = "Student prefix must be at most 10 characters long"
					}
				case "TeacherPrefix":
					switch fieldErr.Tag() {
					case "required":
						errorsMap["teacher_prefix"] = "Teacher Prefix is required"
					case "min":
						errorsMap["teacher_prefix"] = "Teacher prefix must be at least 1 characters long"
					case "max":
						errorsMap["teacher_prefix"] = "Teacher prefix must be at most 10 characters long"
					}
				}
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": errorsMap})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logo, err := c.FormFile("logo")
	if err == nil {
		imageURL, err := utils.UploadFile(context.Background(), logo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
			return
		}
		input.Logo = &imageURL
	} else {
		input.Logo = nil
	}

	favicon, err := c.FormFile("favicon")
	if err == nil {
		imageURL, err := utils.UploadFile(context.Background(), favicon)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
			return
		}
		input.Favicon = &imageURL
	} else {
		input.Favicon = nil
	}

	if output, err := json.MarshalIndent(input, "", "  "); err == nil {
		fmt.Println("Parsed Input:\n", string(output))
	}

	if err := h.service.UpdateGeneralSettings(&input, c.GetUint("tenant_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "General settings updated successfully"})
}
