package instructor

import (
	"dashlearn/internal/models"
	"dashlearn/internal/response"
	"dashlearn/internal/utils"
	"errors"
	"fmt"

	"github.com/lucsky/cuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type InstructorService interface {
	GetInstructors(tenantID uint) ([]models.Instructor, error)
	GetInstructorsLite(tenantID uint) ([]response.InstructorResponse, error)
	GetInstructorDetails(tenantID uint, id uint) (*response.InstructorDetailsResponse, error)
	CreateInstructor(input CreateInstructorInput, tenantID uint) error
	UpdateInstructor(input UpdateInstructorInput, tenantID uint, id uint) error
	DeleteInstructor(tenantID uint, id uint) error
}

type instructorService struct {
	db *gorm.DB
}

func NewInstructorService(db *gorm.DB) InstructorService {
	return &instructorService{
		db: db,
	}
}

func (s *instructorService) GetInstructors(tenantID uint) ([]models.Instructor, error) {
	var users []models.Instructor
	if err := s.db.Where("tenant_id = ?", tenantID).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (s *instructorService) GetInstructorsLite(tenantID uint) ([]response.InstructorResponse, error) {
	var instructors []models.Instructor
	var instructorResponses []response.InstructorResponse

	if err := s.db.
		Where("tenant_id = ?", tenantID).
		Select("id", "first_name", "last_name", "email", "image", "role", "designation").
		Find(&instructors).Error; err != nil {
		return nil, err
	}

	for _, instructor := range instructors {
		res := response.InstructorResponse{
			ID:          instructor.ID,
			FirstName:   instructor.FirstName,
			LastName:    instructor.LastName,
			Email:       instructor.Email,
			Image:       instructor.Image,
			Role:        instructor.Role,
			Designation: instructor.Designation,
		}
		instructorResponses = append(instructorResponses, res)
	}

	return instructorResponses, nil
}

func (s *instructorService) GetInstructorDetails(tenantID uint, id uint) (*response.InstructorDetailsResponse, error) {
	var instructor models.Instructor

	if err := s.db.Where("tenant_id = ? AND id = ?", tenantID, id).First(&instructor).Error; err != nil {
		return nil, err
	}

	instructorRes := response.InstructorDetailsResponse{
		ID:          instructor.ID,
		FirstName:   instructor.FirstName,
		LastName:    utils.ZeroToNil(instructor.LastName),
		Phone:       utils.ZeroToNil(instructor.Phone),
		Email:       instructor.Email,
		Role:        utils.ZeroToNil(instructor.Role),
		Designation: utils.ZeroToNil(instructor.Designation),
		Image:       utils.ZeroToNil(instructor.Image),
	}

	return &instructorRes, nil
}

func (s *instructorService) CreateInstructor(input CreateInstructorInput, tenantID uint) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("something went wrong. Please try again")
	}

	var existingUser models.Instructor
	if err := utils.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		return errors.New("email already exists")
	} else if err != gorm.ErrRecordNotFound {
		return errors.New("something went wrong. Please try again")
	}

	newUser := models.Instructor{
		UserID:      cuid.New(),
		FirstName:   input.FirstName,
		LastName:    utils.EmptyStringToNil(input.LastName),
		Phone:       utils.EmptyStringToNil(input.Phone),
		Role:        utils.EmptyStringToNil(input.Role),
		Designation: utils.EmptyStringToNil(input.Designation),
		Image:       input.ImageURL,
		Email:       input.Email,
		Password:    string(hashedPassword),
		Status:      true,
		TenantID:    tenantID,
	}

	if err := utils.DB.Create(&newUser).Error; err != nil {
		return errors.New("failed to create instructor")
	}

	return nil
}

func (s *instructorService) UpdateInstructor(input UpdateInstructorInput, tenantID uint, id uint) error {
	var instructor models.Instructor
	if err := utils.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&instructor).Error; err != nil {
		return errors.New("instructor not found")
	}

	updates := map[string]interface{}{
		"first_name":  input.FirstName,
		"last_name":   utils.EmptyStringToNil(input.LastName),
		"phone":       utils.EmptyStringToNil(input.Phone),
		"role":        utils.EmptyStringToNil(input.Role),
		"designation": utils.EmptyStringToNil(input.Designation), // can be nil → will set DB NULL
	}

	// Handle image update
	if input.ImageURL != nil && *input.ImageURL != "" {
		updates["image"] = input.ImageURL

		if instructor.Image != nil {
			if delErr := utils.DeleteFromBunny(*instructor.Image); delErr != nil {
				fmt.Println("Failed to delete old file:", delErr)
			}
		}
	}

	return utils.DB.Model(&instructor).Updates(updates).Error
}

func (s *instructorService) DeleteInstructor(tenantID uint, id uint) error {
	var instructor models.Instructor

	if err := utils.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&instructor).Error; err != nil {
		return errors.New("instructor not found")
	}

	if instructor.Image != nil && *instructor.Image != "" {
		if delErr := utils.DeleteFromBunny(*instructor.Image); delErr != nil {
			// You can log or ignore deletion errors as per your need
			fmt.Println("Failed to delete old file:", delErr)
		}
	}

	if err := utils.DB.
		Where("id = ? AND tenant_id = ?", id, tenantID).
		Delete(&models.Instructor{}).Error; err != nil {
		return errors.New("something went wrong. Please try again")
	}

	return nil
}
