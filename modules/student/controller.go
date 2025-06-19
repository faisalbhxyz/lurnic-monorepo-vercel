package student

import (
	"dashlearn/models"
	"dashlearn/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lucsky/cuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetStudents(ctx *gin.Context) {
	var users []models.Student
	utils.DB.Preload("Enrollments").Where("tenant_id = ?", ctx.GetUint("tenant_id")).Find(&users)
	ctx.JSON(http.StatusOK, gin.H{"data": users})
}

func GetStudentLite(ctx *gin.Context) {
	var users []struct {
		ID        uint    `json:"id"`
		FirstName string  `json:"first_name"`
		LastName  *string `json:"last_name"`
	}
	if err := utils.DB.Table("students").Where("tenant_id = ?", ctx.GetUint("tenant_id")).
		Select("id", "first_name", "last_name", "tenant_id").
		Find(&users).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": users})
}

func CreateStudent(ctx *gin.Context) {
	var input CreateStudentInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			errorsMap := make(map[string]string)
			for _, fieldErr := range validationErrors {
				field := fieldErr.Field()
				tag := fieldErr.Tag()

				switch field {
				case "FirstName":
					switch tag {
					case "required":
						errorsMap["firstname"] = "First Name is required"
					}
				case "Email":
					switch tag {
					case "required":
						errorsMap["email"] = "Email is required"
					case "email":
						errorsMap["email"] = "Invalid email format"
					}
				case "Password":
					switch tag {
					case "required":
						errorsMap["password"] = "Password is required"
					case "min":
						errorsMap["password"] = "Password must be at least 6 characters long"
					}
				}
			}
		}

		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong. Please try again."})
		return
	}

	var existingUser models.Student
	if err := utils.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	} else if err != gorm.ErrRecordNotFound {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong. Please try again."})
		return
	}

	newUser := models.Student{
		UserID:    cuid.New(),
		FirstName: input.FirstName,
		LastName:  utils.ZeroToNil(input.LastName),
		Phone:     utils.ZeroToNil(input.Phone),
		Email:     input.Email,
		Password:  string(hashedPassword),
		Status:    true,
		TenantID:  ctx.GetUint("tenant_id"),
	}

	if err := utils.DB.Create(&newUser).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create student"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Student created successfully"})
}

func LoginStudent(ctx *gin.Context) {
	var input LoginStudentInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
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
			ctx.JSON(http.StatusBadRequest, gin.H{"error": errorsMap})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.Student
	err := utils.DB.Where("email = ?", input.Email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email not found"})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong. Please try again."})
		return
	}

	// Compare the provided password with the stored hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"user_id": user.UserID,
			"name": user.FirstName + func() string {
				if user.LastName != nil {
					return " " + *user.LastName
				}
				return ""
			}(),
			"phone": user.Phone,
			"email": user.Email,
		},
	})

}
