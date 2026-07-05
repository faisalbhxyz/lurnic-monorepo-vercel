package student

import (
	"dashlearn/internal/models"
	"dashlearn/internal/utils"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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
	if err := c.ShouldBindWith(&input, binding.FormMultipart); err != nil {
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

	profileImageURL, uploadErr := parseProfileImageUpload(c)
	if uploadErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": uploadErr.Error()})
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
		UserID:       cuid.New(),
		FirstName:    input.FirstName,
		LastName:     utils.ZeroToNil(input.LastName),
		Phone:        utils.ZeroToNil(input.Phone),
		Email:        input.Email,
		Password:     string(hashedPassword),
		ProfileImage: profileImageURL,
		Status:       true,
		TenantID:     c.GetUint("tenant_id"),
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

func ForgotPasswordStudent(c *gin.Context) {
	var input ForgotPasswordInput
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
				case "ResetURL":
					switch tag {
					case "required":
						errorsMap["reset_url"] = "Reset URL is required"
					case "url":
						errorsMap["reset_url"] = "Reset URL must be a valid URL"
					}
				}
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": errorsMap})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID := c.GetUint("tenant_id")
	genericMessage := gin.H{
		"message": "If an account exists for this email, a password reset link has been sent.",
	}

	var user models.Student
	err := utils.DB.Where("email = ? AND tenant_id = ?", input.Email, tenantID).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, genericMessage)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong. Please try again."})
		return
	}

	resetToken, err := utils.GeneratePasswordResetJWT(user.UserID, user.Email, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong. Please try again."})
		return
	}

	resetLink, err := buildPasswordResetLink(input.ResetURL, resetToken, user.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"reset_url": "Reset URL must include scheme and host"}})
		return
	}

	smtpCfg := utils.LoadSMTPConfig()
	if smtpCfg.Enabled() {
		subject := "Reset your password"
		body := fmt.Sprintf(
			"Hello %s,\n\nWe received a request to reset your password. Open the link below to choose a new password. This link expires in 1 hour.\n\n%s\n\nIf you did not request this, you can ignore this email.\n",
			user.FirstName,
			resetLink,
		)
		if err := smtpCfg.Send(user.Email, subject, body); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send reset email. Please try again later."})
			return
		}
	} else if strings.EqualFold(os.Getenv("GIN_MODE"), "debug") {
		response := gin.H{
			"message":         genericMessage["message"],
			"dev_reset_link":  resetLink,
			"dev_reset_token": resetToken,
		}
		c.JSON(http.StatusOK, response)
		return
	} else {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Password reset email is not configured"})
		return
	}

	c.JSON(http.StatusOK, genericMessage)
}

func ResetPasswordStudent(c *gin.Context) {
	var input ResetPasswordInput
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
				case "Token":
					if tag == "required" {
						errorsMap["token"] = "Reset token is required"
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

	claims, err := utils.ParsePasswordResetJWT(input.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired reset token"})
		return
	}

	tenantID := c.GetUint("tenant_id")
	if claims.TenantID != tenantID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired reset token"})
		return
	}

	if !strings.EqualFold(claims.Email, input.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired reset token"})
		return
	}

	var user models.Student
	err = utils.DB.Where("user_id = ? AND email = ? AND tenant_id = ?", claims.UserID, input.Email, tenantID).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired reset token"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong. Please try again."})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong. Please try again."})
		return
	}

	if err := utils.DB.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong. Please try again."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

func buildPasswordResetLink(resetURL, token, email string) (string, error) {
	parsed, err := url.Parse(strings.TrimSpace(resetURL))
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return "", errors.New("reset URL must include scheme and host")
	}

	query := parsed.Query()
	query.Set("token", token)
	query.Set("email", email)
	parsed.RawQuery = query.Encode()
	return parsed.String(), nil
}

func GetStudentDetailsByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	details, err := buildStudentAdminDetails(utils.DB, c.GetUint("tenant_id"), uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong. Please try again."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": details})
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
	if err := c.ShouldBindWith(&input, binding.FormMultipart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profileImageURL, uploadErr := parseProfileImageUpload(c)
	if uploadErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": uploadErr.Error()})
		return
	}
	input.ProfileImageURL = profileImageURL

	tenantID := c.GetUint("tenant_id")

	// Optional: Check if student exists first (skip if you're fine with silent fail)
	var student models.Student
	if err := utils.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&student).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	updates := map[string]interface{}{
		"first_name": input.FirstName,
		"last_name":  input.LastName,
		"phone":      input.Phone,
	}

	if input.ProfileImageURL != nil && *input.ProfileImageURL != "" {
		updates["profile_image"] = input.ProfileImageURL
		if student.ProfileImage != nil && *student.ProfileImage != "" {
			if delErr := utils.DeleteFromBunny(*student.ProfileImage); delErr != nil {
				fmt.Println("Failed to delete old profile image:", delErr)
			}
		}
	}

	if err := utils.DB.Model(&student).Updates(updates).Error; err != nil {
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
