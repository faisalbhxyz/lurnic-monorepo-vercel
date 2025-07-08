package user

import (
	"context"
	"dashlearn/internal/models"
	"dashlearn/internal/utils"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lucsky/cuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CheckUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "WORKING"})
}
func GetUsers(c *gin.Context) {
	var users []models.User
	utils.DB.Preload("Tenant").Select("id", "user_id", "name", "phone", "email", "status", "created_at", "updated_at", "tenant_id").Find(&users)
	c.JSON(http.StatusOK, gin.H{"data": users})
}

func CreateUser(c *gin.Context) {
	var input CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {

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
				case "Email":
					if tag == "required" {
						errorsMap["email"] = "Email is required"
					} else if tag == "email" {
						errorsMap["email"] = "Invalid email format"
					}
				case "Password":
					if tag == "required" {
						errorsMap["password"] = "Password is required"
					} else if tag == "min" {
						errorsMap["password"] = "Password must be at least 6 characters long"
					}
				}
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": errorsMap})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong. Please try again."})
		return
	}

	var existingUser models.User
	if err := utils.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong. Please try again."})
		return
	}

	newTenant := models.Tenant{
		AppKey: cuid.New(),
	}

	if err := utils.DB.Create(&newTenant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong. Please try again."})
		return
	}

	var superadmin models.Role
	if err := utils.DB.Where("name = ?", "superadmin").First(&superadmin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong. Please try again."})
		return
	}

	newUser := models.User{
		UserID:   cuid.New(),
		Name:     input.Name,
		Phone:    input.Phone,
		Email:    input.Email,
		Password: string(hashedPassword),
		Status:   true,
		TenantID: newTenant.ID,
		RoleID:   &superadmin.ID,
	}

	if err := utils.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func LoginUser(c *gin.Context) {
	var input LoginUserInput

	if err := c.ShouldBindJSON(&input); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			errorsMap := make(map[string]string)
			for _, fieldErr := range validationErrors {
				field := fieldErr.Field()
				tag := fieldErr.Tag()
				switch field {
				case "Email":
					if tag == "required" {
						errorsMap["email"] = "Email is required"
					} else if tag == "email" {
						errorsMap["email"] = "Invalid email format"
					}
				case "Password":
					if tag == "required" {
						errorsMap["password"] = "Password is required"
					} else if tag == "min" {
						errorsMap["password"] = "Password must be at least 6 characters long"
					}
				}
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": errorsMap})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	err := utils.DB.Where("email = ?", input.Email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong. Please try again."})
		return
	}

	// Compare the provided password with the stored hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"user_id": user.UserID,
			"name":    user.Name,
			"phone":   user.Phone,
			"email":   user.Email,
		},
	})

}

func UploadUser(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	fmt.Println("FILE", fileHeader.Filename)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	// Call your UploadFile util function
	url, err := utils.UploadFile(context.Background(), fileHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the uploaded file URL
	c.JSON(http.StatusOK, gin.H{
		"url": url,
	})
}
