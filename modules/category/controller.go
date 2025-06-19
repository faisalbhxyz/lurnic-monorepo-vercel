package category

import (
	"dashlearn/models"
	"dashlearn/utils"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func GetAllCategory(ctx *gin.Context) {
	var categories []models.Category
	utils.DB.Where("tenant_id = ?", ctx.GetUint("tenant_id")).Find(&categories)
	ctx.JSON(http.StatusOK, gin.H{"data": categories})
}

func GetAllCategoryPublic(ctx *gin.Context) {
	var categories []models.Category
	utils.DB.Where("tenant_id = ?", ctx.GetUint("tenant_id")).Find(&categories)
	ctx.JSON(http.StatusOK, gin.H{"data": categories})
}

func GetCategoryByID(ctx *gin.Context) {
	categoryID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var category models.Category
	utils.DB.Where("id = ? AND tenant_id = ?", categoryID, ctx.GetUint("tenant_id")).First(&category)
	ctx.JSON(http.StatusOK, gin.H{"data": category})
}

func CreateCategory(ctx *gin.Context) {
	var input CreateCategoryInput

	if err := ctx.ShouldBindJSON(&input); err != nil {

		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			errorsMap := make(map[string]string)
			for _, fieldErr := range validationErrors {
				field := fieldErr.Field()
				tag := fieldErr.Tag()

				switch field {
				case "Name":
					if tag == "required" {
						errorsMap["name"] = "Name is required"
					}
				case "Slug":
					if tag == "required" {
						errorsMap["slug"] = "Slug is required"
					}
				}
			}
			ctx.JSON(http.StatusBadRequest, gin.H{"errors": errorsMap})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if utils.DB.Where("slug = ? AND tenant_id = ?", input.Slug, ctx.GetUint("tenant_id")).First(&models.Category{}).RowsAffected > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Category with this slug already exists"})
		return
	}

	newCategory := models.Category{
		Name:        input.Name,
		Slug:        input.Slug,
		Description: utils.EmptyStringToNil(input.Description),
		TenantID:    ctx.GetUint("tenant_id"),
	}

	if err := utils.DB.Create(&newCategory).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Category created successfully"})
}

func UpdateCategory(ctx *gin.Context) {
	var category models.Category
	categoryID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	if err := utils.DB.Where("id = ? AND tenant_id = ?", categoryID, ctx.GetUint("tenant_id")).First(&category).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	var input CreateCategoryInput

	if err := ctx.ShouldBindJSON(&input); err != nil {

		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			errorsMap := make(map[string]string)
			for _, fieldErr := range validationErrors {
				field := fieldErr.Field()
				tag := fieldErr.Tag()

				switch field {
				case "Name":
					if tag == "required" {
						errorsMap["name"] = "Name is required"
					}
				case "Slug":
					if tag == "required" {
						errorsMap["slug"] = "Slug is required"
					}
				}
			}
			ctx.JSON(http.StatusBadRequest, gin.H{"errors": errorsMap})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if utils.DB.Where("slug = ? AND tenant_id = ? AND id != ?", input.Slug, ctx.GetUint("tenant_id"), categoryID).First(&models.Category{}).RowsAffected > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Category with this slug already exists"})
		return
	}

	//update category
	category.Name = input.Name
	category.Slug = input.Slug
	category.Description = utils.EmptyStringToNil(input.Description)

	if err := utils.DB.Save(&category).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Category updated successfully"})
}
