package paymentmethod

import (
	"dashlearn/internal/models"
	"dashlearn/internal/utils"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type PaymentMethodService interface {
	GetAll(tenantID uint, payment_status *bool) ([]PaymentMethodResponse, error)
	GetByID(tenantID uint, id uint64) (*PaymentMethodResponse, error)
	Create(input CreatePaymentMethodInput, tenantID uint) error
	Update(id uint, input UpdatePaymentMethodInput, tenantID uint) error
	Delete(id uint, tenantID uint) error
}

type paymentMethodService struct {
	db *gorm.DB
}

func NewPaymentMethodService(db *gorm.DB) PaymentMethodService {
	return &paymentMethodService{
		db: db,
	}
}

func (s *paymentMethodService) GetAll(tenantID uint, payment_status *bool) ([]PaymentMethodResponse, error) {
	var payments []models.PaymentMethod
	var res []PaymentMethodResponse
	var err error

	if payment_status != nil {
		err = s.db.Where("tenant_id = ? AND status = ?", tenantID, *payment_status).Find(&payments).Error
	} else {
		err = s.db.Where("tenant_id = ?", tenantID).Find(&payments).Error
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	for _, payment := range payments {
		res = append(res, PaymentMethodResponse{
			ID:          payment.ID,
			Title:       payment.Title,
			Image:       payment.Image,
			Instruction: payment.Instruction,
			Status:      payment.Status,
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
		Status:      payment.Status,
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

func (s *paymentMethodService) Update(id uint, input UpdatePaymentMethodInput, tenantID uint) error {
	var payment models.PaymentMethod

	// 1. Find the payment method under the correct tenant
	if err := s.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&payment).Error; err != nil {
		return err
	}

	// 2. Check uniqueness of title (except current)
	var existing models.PaymentMethod
	if err := s.db.
		Where("title = ? AND tenant_id = ? AND id != ?", input.Title, tenantID, id).
		First(&existing).Error; err == nil {
		return errors.New("payment method with this title already exists")
	}

	// 3. If new image provided, delete old one
	if input.Image != nil && *input.Image != "" {
		if payment.Image != nil && *payment.Image != "" {
			if delErr := utils.DeleteFromBunny(*payment.Image); delErr != nil {
				fmt.Println("Failed to delete old file:", delErr)
			}
		}
	}

	// 4. Update fields
	payment.Title = input.Title
	payment.Image = input.Image
	payment.Status = input.Status == "true"
	payment.Instruction = input.Instruction

	// 5. Save changes
	return s.db.Save(&payment).Error
}

func (s *paymentMethodService) Delete(id uint, tenantID uint) error {
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
