package subcategory

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type SubCategoryHandler struct {
	service SubCategoryService
}

func NewCategoryHandler(db *gorm.DB) *SubCategoryHandler {
	return &SubCategoryHandler{
		service: NewSubCategoryService(db),
	}
}

func (h *SubCategoryHandler) GetAll(c *gin.Context) {
	tenantID := c.GetUint("tenant_id")
	categories, err := h.service.GetAll(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": categories})
}

func (h *SubCategoryHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}
	tenantID := c.GetUint("tenant_id")

	category, err := h.service.GetByID(tenantID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": category})
}

func (h *SubCategoryHandler) Create(c *gin.Context) {
	var input CreateSubCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		handleValidationError(c, err)
		return
	}

	if err := h.service.Create(input, c.GetUint("tenant_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Sub category created successfully"})
}

func (h *SubCategoryHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var input CreateSubCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		handleValidationError(c, err)
		return
	}

	if err := h.service.Update(id, input, c.GetUint("tenant_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Sub category updated successfully"})
}

func handleValidationError(c *gin.Context, err error) {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		errorsMap := make(map[string]string)
		for _, fieldErr := range validationErrors {
			switch fieldErr.Field() {
			case "Name":
				if fieldErr.Tag() == "required" {
					errorsMap["name"] = "Name is required"
				}
			case "Slug":
				if fieldErr.Tag() == "required" {
					errorsMap["slug"] = "Slug is required"
				}
			}
		}
		c.JSON(http.StatusBadRequest, gin.H{"errors": errorsMap})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}

func (h *SubCategoryHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}
	if err := h.service.Delete(id, c.GetUint("tenant_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
