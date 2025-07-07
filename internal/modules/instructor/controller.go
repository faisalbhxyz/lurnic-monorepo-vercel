package instructor

import (
	"context"
	"dashlearn/models"
	"dashlearn/utils"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/lucsky/cuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetInstructors(c *gin.Context) {
	var users []models.Instructor
	utils.DB.Where("tenant_id = ?", c.GetUint("tenant_id")).Find(&users)
	c.JSON(http.StatusOK, gin.H{"data": users})
}

func GetInstructorsLite(c *gin.Context) {
	var users []models.InstructorResponseLite
	utils.DB.Where("tenant_id = ?", c.GetUint("tenant_id")).Find(&users)
	c.JSON(http.StatusOK, gin.H{"data": users})
}

func GetInstructorDetails(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}
	var user models.InstructorDetailsResponse
	utils.DB.Where("tenant_id = ? AND id = ?", c.GetUint("tenant_id"), id).First(&user)
	c.JSON(http.StatusOK, gin.H{"data": user})
}

func CreateInstructor(c *gin.Context) {
	var input CreateInstructorInput
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

	var imageURL *string
	file, err := c.FormFile("image")
	if err == nil {
		if file.Size > 2*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Max image size is 2MB"})
			return
		}

		// ✅ 2. MIME type check
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open image file"})
			return
		}
		defer src.Close()

		// Detect content type
		buffer := make([]byte, 512)
		if _, err := src.Read(buffer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read file content"})
			return
		}
		contentType := http.DetectContentType(buffer)

		allowedTypes := map[string]bool{
			"image/jpeg": true,
			"image/png":  true,
			"image/jpg":  true,
			// "image/webp": true,
			// "image/gif":  true,
		}
		if !allowedTypes[contentType] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Only PNG, JPG formats are supported"})
			return
		}

		// ✅ 3. (Optional) Image dimension check
		// if contentType == "image/jpeg" || contentType == "image/png" {
		// 	// Need to re-seek for reading again
		// 	if seeker, ok := src.(io.Seeker); ok {
		// 		seeker.Seek(0, io.SeekStart)
		// 	}

		// 	img, _, err := image.Decode(src)
		// 	if err != nil {
		// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode image"})
		// 		return
		// 	}
		// 	width := img.Bounds().Dx()
		// 	height := img.Bounds().Dy()

		// 	if width > 1920 || height > 1080 {
		// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Image must be 1920x1080 pixels or smaller"})
		// 		return
		// 	}
		// }

		// save file
		url, err := utils.UploadFile(context.Background(), file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		imageURL = &url
	} else {
		imageURL = nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong. Please try again."})
		return
	}

	var existingUser models.Instructor
	if err := utils.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong. Please try again."})
		return
	}

	newUser := models.Instructor{
		UserID:      cuid.New(),
		FirstName:   input.FirstName,
		LastName:    utils.EmptyStringToNil(input.LastName),
		Phone:       utils.EmptyStringToNil(input.Phone),
		Role:        utils.EmptyStringToNil(input.Role),
		Designation: utils.EmptyStringToNil(input.Designation),
		Image:       imageURL,
		Email:       input.Email,
		Password:    string(hashedPassword),
		Status:      true,
		TenantID:    c.GetUint("tenant_id"),
	}

	if err := utils.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create instructor"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Instructor created successfully"})
}

// func LoginUser(c *gin.Context) {
// 	var input LoginUserInput

// 	if err := c.ShouldBindJSON(&input); err != nil {
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
// 			c.JSON(http.StatusBadRequest, gin.H{"error": errorsMap})
// 			return
// 		}
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	var user User
// 	err := utils.DB.Where("email = ?", input.Email).First(&user).Error
// 	if errors.Is(err, gorm.ErrRecordNotFound) {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Email not found"})
// 		return
// 	} else if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong. Please try again."})
// 		return
// 	}

// 	// Compare the provided password with the stored hashed password
// 	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
// 		return
// 	}

// 	// Generate JWT token
// 	token, err := utils.GenerateJWT(user.UserID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"token": token,
// 		"user": gin.H{
// 			"user_id": user.UserID,
// 			"name":    user.Name,
// 			"phone":   user.Phone,
// 			"email":   user.Email,
// 		},
// 	})

// }

func UpdateInstructor(c *gin.Context) {

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	var input UpdateInstructorInput
	if err := c.ShouldBindWith(&input, binding.FormMultipart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID := c.GetUint("tenant_id")

	var instructor models.Instructor
	if err := utils.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&instructor).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Instructor not found"})
		return
	}

	var imageURL *string
	file, err := c.FormFile("image")

	if file != nil && instructor.Image != nil {
		if delErr := utils.DeleteCDNFile(context.Background(), *instructor.Image); delErr != nil {
			// You can log or ignore deletion errors as per your need
			fmt.Println("Failed to delete old file:", delErr)
		}
	}

	if err == nil {
		if file.Size > 2*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Max image size is 2MB"})
			return
		}

		// ✅ 2. MIME type check
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open image file"})
			return
		}
		defer src.Close()

		// Detect content type
		buffer := make([]byte, 512)
		if _, err := src.Read(buffer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read file content"})
			return
		}
		contentType := http.DetectContentType(buffer)

		allowedTypes := map[string]bool{
			"image/jpeg": true,
			"image/png":  true,
			"image/jpg":  true,
			// "image/webp": true,
			// "image/gif":  true,
		}
		if !allowedTypes[contentType] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Only PNG, JPG formats are supported"})
			return
		}

		// ✅ 3. (Optional) Image dimension check
		// if contentType == "image/jpeg" || contentType == "image/png" {
		// 	// Need to re-seek for reading again
		// 	if seeker, ok := src.(io.Seeker); ok {
		// 		seeker.Seek(0, io.SeekStart)
		// 	}

		// 	img, _, err := image.Decode(src)
		// 	if err != nil {
		// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode image"})
		// 		return
		// 	}
		// 	width := img.Bounds().Dx()
		// 	height := img.Bounds().Dy()

		// 	if width > 1920 || height > 1080 {
		// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Image must be 1920x1080 pixels or smaller"})
		// 		return
		// 	}
		// }

		// save file
		url, err := utils.UploadFile(context.Background(), file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		imageURL = &url
	} else {
		imageURL = nil
	}

	updateUser := models.Instructor{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Phone:     input.Phone,
		Image:     imageURL,
	}

	utils.DB.Model(&instructor).Updates(updateUser)

	c.JSON(http.StatusOK, gin.H{"message": "Instructor updated successfully"})

}

func DeleteInstructor(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	tenantID := c.GetUint("tenant_id")

	var instructor models.Instructor
	if err := utils.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&instructor).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Instructor not found"})
		return
	}

	if instructor.Image != nil {
		if delErr := utils.DeleteCDNFile(context.Background(), *instructor.Image); delErr != nil {
			// You can log or ignore deletion errors as per your need
			fmt.Println("Failed to delete old file:", delErr)
		}
	}

	if err := utils.DB.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&models.Instructor{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Instructor deleted successfully"})
}
