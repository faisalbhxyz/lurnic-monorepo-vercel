package user

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type UserHandler struct {
	service UserService
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{
		service: NewUserService(db),
	}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	users := h.service.GetUsers()
	c.JSON(http.StatusOK, gin.H{"data": users})
}

func (h *UserHandler) CreateUser(c *gin.Context) {
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
			c.JSON(http.StatusBadRequest, gin.H{"errors": errorsMap})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.CreateUser(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (h *UserHandler) LoginUser(c *gin.Context) {
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

	user, err := h.service.LoginUser(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)

}

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

func (h *UserHandler) CreateTeamMember(c *gin.Context) {
	var input CreateTeamMemberInput
	if err := c.ShouldBindJSON(&input); err != nil {

		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			errorsMap := make(map[string]string)
			for _, fieldErr := range validationErrors {
				field := fieldErr.Field()
				tag := fieldErr.Tag()

				switch field {
				case "UserID":
					if tag == "required" {
						errorsMap["user_id"] = "User ID is required"
					}
				case "Name":
					if tag == "required" {
						errorsMap["name"] = "Name is required"
					}
				case "Role":
					if tag == "required" {
						errorsMap["role"] = "Role is required"
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
			c.JSON(http.StatusBadRequest, gin.H{"errors": errorsMap})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.CreateTeamMember(input, c.GetUint("tenant_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (h *UserHandler) GetTeamMembers(c *gin.Context) {
	users := h.service.GetTeamMembers(c.GetUint("tenant_id"))
	c.JSON(http.StatusOK, gin.H{"data": users})
}

func (h *UserHandler) GetTeamMemberByUID(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	user := h.service.GetTeamMemberByUID(c.GetUint("tenant_id"), uint(userID))

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (h *UserHandler) UpdateTeamMember(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	fmt.Println("USER", userID)

	var input UpdateTeamMemberInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.service.UpdateTeamMember(c.GetUint("tenant_id"), uint(userID), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (h *UserHandler) DeleteTeamMember(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.service.DeleteTeamMember(uint(userID), c.GetUint("tenant_id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
