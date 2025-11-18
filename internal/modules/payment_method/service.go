package paymentmethod

import (
	"dashlearn/internal/models"
	"dashlearn/internal/utils"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type PaymentMethodService interface {
	GetAll(tenantID uint) ([]PaymentMethodResponse, error)
	GetByID(tenantID uint, id uint64) (*PaymentMethodResponse, error)
	Create(input CreatePaymentMethodInput, tenantID uint) error
	Update(id uint64, input CreatePaymentMethodInput, tenantID uint) error
	Delete(id uint64, tenantID uint) error
}

type paymentMethodService struct {
	db *gorm.DB
}

func NewPaymentMethodService(db *gorm.DB) PaymentMethodService {
	return &paymentMethodService{
		db: db,
	}
}

func (s *paymentMethodService) GetAll(tenantID uint) ([]PaymentMethodResponse, error) {
	var payments []models.PaymentMethod
	var res []PaymentMethodResponse

	err := s.db.Where("tenant_id = ?", tenantID).Find(&payments).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	for _, payment := range payments {
		res = append(res, PaymentMethodResponse{
			ID:          payment.ID,
			Title:       payment.Title,
			Image:       payment.Image,
			Instruction: payment.Instruction,
			CreatedAt:   payment.CreatedAt,
			UpdatedAt:   payment.UpdatedAt,
		})
	}

	return res, nil
}

func (s *paymentMethodService) GetByID(tenantID uint, id uint64) (*PaymentMethodResponse, error) {
	var payment models.PaymentMethod
	var response PaymentMethodResponse
	result := s.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&payment)
	if result.Error != nil {
		return nil, result.Error
	}

	response = PaymentMethodResponse{
		ID:          payment.ID,
		Title:       payment.Title,
		Image:       payment.Image,
		Instruction: payment.Instruction,
		CreatedAt:   payment.CreatedAt,
		UpdatedAt:   payment.UpdatedAt,
	}

	return &response, nil
}

func (s *paymentMethodService) Create(input CreatePaymentMethodInput, tenantID uint) error {
	if s.db.Where("title = ? AND tenant_id = ?", input.Title, tenantID).First(&models.PaymentMethod{}).RowsAffected > 0 {
		return errors.New("payment method with this title already exists")
	}

	method := models.PaymentMethod{
		Title:       input.Title,
		Image:       input.Image,
		Instruction: input.Instruction,
		TenantID:    tenantID,
	}

	return s.db.Create(&method).Error
}

func (s *paymentMethodService) Update(id uint64, input CreatePaymentMethodInput, tenantID uint) error {
	var category models.Category

	if err := s.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&category).Error; err != nil {
		return err
	}

	// if s.db.Where("slug = ? AND tenant_id = ? AND id != ?", input.Slug, tenantID, id).First(&models.Category{}).RowsAffected > 0 {
	// 	return errors.New("category with this slug already exists")
	// }

	// category.Name = input.Name
	// category.Slug = input.Slug
	// category.Description = utils.EmptyStringToNil(input.Description)

	return s.db.Save(&category).Error
}

func (s *paymentMethodService) Delete(id uint64, tenantID uint) error {
	var payment models.PaymentMethod

	if err := s.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&payment).Error; err != nil {
		return errors.New("payment method not found")
	}

	if payment.Image != nil && *payment.Image != "" {
		if delErr := utils.DeleteFromBunny(*payment.Image); delErr != nil {
			// You can log or ignore deletion errors as per your need
			fmt.Println("Failed to delete old file:", delErr)
		}
	}

	if err := s.db.
		Where("id = ? AND tenant_id = ?", id, tenantID).
		Delete(&models.PaymentMethod{}).Error; err != nil {
		return errors.New("something went wrong. Please try again")
	}

	return nil
}
