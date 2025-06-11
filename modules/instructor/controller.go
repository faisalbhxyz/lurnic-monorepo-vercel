package instructor

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

// func GetUsers(ctx *gin.Context) {
// 	var users []User
// 	utils.DB.Preload("Tenant").Select("id", "user_id", "name", "phone", "email", "status", "created_at", "updated_at", "tenant_id").Find(&users)
// 	ctx.JSON(http.StatusOK, gin.H{"data": users})
// }

func CreateInstructor(ctx *gin.Context) {
	var input CreateInstructorInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			errorsMap := make(map[string]string)
			for _, fieldErr := range validationErrors {
				field := fieldErr.Field()
				tag := fieldErr.Tag()

				switch field {
				case "FirstName":
					if tag == "required" {
						errorsMap["firstname"] = "First Name is required"
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
		}

		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong. Please try again."})
		return
	}

	var existingUser models.Instructor
	if err := utils.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	} else if err != gorm.ErrRecordNotFound {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong. Please try again."})
		return
	}

	newUser := models.Instructor{
		UserID:    cuid.New(),
		FirstName: input.FirstName,
		LastName:  utils.EmptyStringToNil(input.LastName),
		Phone:     utils.EmptyStringToNil(input.Phone),
		Email:     input.Email,
		Password:  string(hashedPassword),
		Status:    true,
		TenantID:  ctx.GetUint("tenant_id"),
	}

	if err := utils.DB.Create(&newUser).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create instructor"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Instructor created successfully"})
}

// func LoginUser(ctx *gin.Context) {
// 	var input LoginUserInput

// 	if err := ctx.ShouldBindJSON(&input); err != nil {
// 		var validationErrors validator.ValidationErrors
// 		if errors.As(err, &validationErrors) {
// 			errorsMap := make(map[string]string)
// 			for _, fieldErr := range validationErrors {
// 				field := fieldErr.Field()
// 				tag := fieldErr.Tag()
// 				switch field {
// 				case "Email":
// 					if tag == "required" {
// 						errorsMap["email"] = "Email is required"
// 					} else if tag == "email" {
// 						errorsMap["email"] = "Invalid email format"
// 					}
// 				case "Password":
// 					if tag == "required" {
// 						errorsMap["password"] = "Password is required"
// 					} else if tag == "min" {
// 						errorsMap["password"] = "Password must be at least 6 characters long"
// 					}
// 				}
// 			}
// 			ctx.JSON(http.StatusBadRequest, gin.H{"error": errorsMap})
// 			return
// 		}
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	var user User
// 	err := utils.DB.Where("email = ?", input.Email).First(&user).Error
// 	if errors.Is(err, gorm.ErrRecordNotFound) {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email not found"})
// 		return
// 	} else if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong. Please try again."})
// 		return
// 	}

// 	// Compare the provided password with the stored hashed password
// 	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
// 		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
// 		return
// 	}

// 	// Generate JWT token
// 	token, err := utils.GenerateJWT(user.UserID)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{
// 		"token": token,
// 		"user": gin.H{
// 			"user_id": user.UserID,
// 			"name":    user.Name,
// 			"phone":   user.Phone,
// 			"email":   user.Email,
// 		},
// 	})

// }

// func UploadUser(c *gin.Context) {
// 	fileHeader, err := c.FormFile("file")
// 	fmt.Println("FILE", fileHeader.Filename)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
// 		return
// 	}

// 	// Call your UploadFile util function
// 	url, err := utils.UploadFile(context.Background(), fileHeader)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Return the uploaded file URL
// 	c.JSON(http.StatusOK, gin.H{
// 		"url": url,
// 	})
// }
