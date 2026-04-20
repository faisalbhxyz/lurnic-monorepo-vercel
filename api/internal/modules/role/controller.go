package role

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type RoleHandler struct {
	service RoleService
}

func NewRoleHandler(db *gorm.DB) *RoleHandler {
	return &RoleHandler{
		service: NewRoleService(db),
	}
}

func (h *RoleHandler) GetRoles(c *gin.Context) {
	roles, err := h.service.GetRoles(c.GetUint("tenant_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": roles})
}

func (h *RoleHandler) GetRoleByUserID(c *gin.Context) {

	user_id := c.GetUint("user_id")
	tenant_id := c.GetUint("tenant_id")

	role, err := h.service.GetRoleByUserID(user_id, tenant_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": role})
}

func (h *RoleHandler) CreateRole(c *gin.Context) {
	var input CreateRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {

		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			errorsMap := make(map[string]string)
			for _, fieldErr := range validationErrors {
				field := fieldErr.Field()
				tag := fieldErr.Tag()

				switch field {
				case "Title":
					if tag == "required" {
						errorsMap["title"] = "Title is required"
					}
				case "Permissions":
					if tag == "required" {
						errorsMap["permissions"] = "Permissions are required"
					}
				}
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": errorsMap})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.CreateRole(input, c.GetUint("tenant_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Role created successfully"})
}

func (h *RoleHandler) DeleteRoleHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("role_id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	if err := h.service.DeleteRole(uint(id), c.GetUint("tenant_id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
}
