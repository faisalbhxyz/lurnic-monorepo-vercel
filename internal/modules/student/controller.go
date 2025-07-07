package student

import (
	"dashlearn/models"
	"dashlearn/utils"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lucsky/cuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetStudents(c *gin.Context) {
	var users []models.Student
	utils.DB.Preload("Enrollments").Where("tenant_id = ?", c.GetUint("tenant_id")).Find(&users)
	c.JSON(http.StatusOK, gin.H{"data": users})
}

func GetStudentLite(c *gin.Context) {
	var users []struct {
		ID        uint    `json:"id"`
		FirstName string  `json:"first_name"`
		LastName  *string `json:"last_name"`
	}
	if err := utils.DB.Table("students").Where("tenant_id = ?", c.GetUint("tenant_id")).
		Select("id", "first_name", "last_name", "tenant_id").
		Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

func CreateStudent(c *gin.Context) {
	var input CreateStudentInput
	if err := c.ShouldBindJSON(&input); err != nil {
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

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong. Please try again."})
		return
	}

	var existingUser models.Student
	if err := utils.DB.Where(
		"(email = ? OR phone = ?) AND tenant_id = ?",
		input.Email, input.Phone, c.GetUint("tenant_id"),
	).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email or phone already exists"})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong. Please try again."})
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
		TenantID:  c.GetUint("tenant_id"),
	}

	if err := utils.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create student"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Student created successfully"})
}

func CreateStudentPublic(c *gin.Context) {
	var input CreateStudentInput
	if err := c.ShouldBindJSON(&input); err != nil {
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

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong. Please try again."})
		return
	}

	var existingUser models.Student
	if err := utils.DB.Where(
		"(email = ? OR phone = ?) AND tenant_id = ?",
		input.Email, input.Phone, c.GetUint("tenant_id"),
	).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email or phone already exists"})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong. Please try again."})
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
		TenantID:  c.GetUint("tenant_id"),
	}

	if err := utils.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create student"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Student created successfully"})
}

func LoginStudent(c *gin.Context) {
	var input LoginStudentInput

	if err := c.ShouldBindJSON(&input); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			errorsMap := make(map[string]string)
			for _, fieldErr := range validationErrors {
				field := fieldErr.Field()
				tag := fieldErr.Tag()
				switch field {
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
			c.JSON(http.StatusBadRequest, gin.H{"error": errorsMap})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.Student
	err := utils.DB.Where("email = ? AND tenant_id = ?", input.Email, c.GetUint("tenant_id")).First(&user).Error
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

func GetStudentDetailsByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}
	var users models.StudentDetailsRes
	utils.DB.
		Where("tenant_id = ? AND id = ?", c.GetUint("tenant_id"), id).
		First(&users)

	c.JSON(http.StatusOK, gin.H{"data": users})
}

func GetStudentDetails(c *gin.Context) {
	var users models.StudentDetailsRes
	utils.DB.
		Preload("Enrollments").
		Where("tenant_id = ? AND id = ?", c.GetUint("tenant_id"), c.GetUint("user_id")).
		First(&users)

	c.JSON(http.StatusOK, gin.H{"data": users})
}

func UpdateStudent(c *gin.Context) {
	// Parse the student ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	// Bind request body
	var input UpdateStudentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID := c.GetUint("tenant_id")

	// Optional: Check if student exists first (skip if you're fine with silent fail)
	var student models.Student
	if err := utils.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&student).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	// Update fields
	if err := utils.DB.
		Model(&student).
		Updates(models.Student{
			FirstName: input.FirstName,
			LastName:  input.LastName,
			Phone:     input.Phone,
		}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Student updated successfully"})
}

func DeleteStudent(c *gin.Context) {
	// Parse the student ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	// Delete the student
	if err := utils.DB.Where("id = ? AND tenant_id = ?", id, c.GetUint("tenant_id")).Delete(&models.Student{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Student deleted successfully"})
}
