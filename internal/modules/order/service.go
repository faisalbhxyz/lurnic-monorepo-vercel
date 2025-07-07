package order

import (
	"dashlearn/models"
	"dashlearn/response"
	"errors"

	"gorm.io/gorm"
)

type OrderService interface {
	GetAll(tenantID uint) ([]response.GetAllOrderResponse, error)
	Create(input CreateOrderInput, tenantID uint, studentID uint) (response.CreateOrderResponse, error)
	Delete(tenantID uint, orderID uint) error
	MarkAsPaid(tenantID uint, orderID uint) error
}

type orderService struct {
	db *gorm.DB
}

func NewOrderService(db *gorm.DB) OrderService {
	return &orderService{
		db: db,
	}
}

func (s *orderService) GetAll(tenantID uint) ([]response.GetAllOrderResponse, error) {
	var orders []response.GetAllOrderResponse
	err := s.db.Where("tenant_id = ?", tenantID).Preload("Course").Preload("Student").Find(&orders).Error
	return orders, err
}

func (s *orderService) Create(input CreateOrderInput, tenantID uint, studentID uint) (response.CreateOrderResponse, error) {
	// check if the course exists
	var course models.CourseDetails
	if err := s.db.Where(models.CourseDetails{
		ID:       input.CourseID,
		TenantID: tenantID,
	}).First(&course).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.CreateOrderResponse{}, errors.New("course not found")
		}
		return response.CreateOrderResponse{}, err
	}

	// Check if an order already exists for the given course and student
	var existingOrder models.Order
	if err := s.db.Where(models.Order{
		CourseID:  input.CourseID,
		StudentID: studentID,
		TenantID:  tenantID,
	}).First(&existingOrder).Error; err == nil {
		return response.CreateOrderResponse{}, errors.New("an order for this course already exists for this student")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// some other DB error
		return response.CreateOrderResponse{}, err
	}

	// next invoice ID
	var nextInvoiceID int64
	if err := s.db.Model(&models.Order{}).
		Where("tenant_id = ?", tenantID).
		Order("invoice_id DESC").
		Limit(1).
		Select("invoice_id").
		Scan(&nextInvoiceID).Error; err != nil {
		return response.CreateOrderResponse{}, err
	}
	nextInvoiceID++

	var coursePrice float64
	if course.PricingModel == models.CoursePricingModelPaid {
		if course.SalePrice == nil || *course.SalePrice == 0 {
			coursePrice = float64(*course.RegularPrice)
		} else {
			coursePrice = float64(*course.SalePrice)
		}
	} else {
		coursePrice = 0
	}

	newOrder := models.Order{
		InvoiceID:     nextInvoiceID,
		StudentID:     studentID,
		CourseID:      input.CourseID,
		TenantID:      tenantID,
		DiscountType:  "none",
		Discount:      0,
		Total:         coursePrice,
		PaymentType:   "manual",
		PaymentStatus: models.OrderPaymentStatusUnpaid,
	}

	if err := s.db.Create(&newOrder).Error; err != nil {
		return response.CreateOrderResponse{}, err
	}

	// build response
	orderResponse := response.CreateOrderResponse{
		ID:           newOrder.ID,
		InvoiceID:    newOrder.InvoiceID,
		CourseID:     newOrder.CourseID,
		Total:        newOrder.Total,
		CustomerNote: newOrder.CustomerNote,
		CreatedAt:    newOrder.CreatedAt,
		UpdatedAt:    newOrder.UpdatedAt,
	}

	return orderResponse, nil
}

func (s *orderService) MarkAsPaid(tenantID uint, orderID uint) error {
	//check if order exists
	var order models.Order
	if err := s.db.Where(&models.Order{
		ID:       orderID,
		TenantID: tenantID,
	}).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("order not found")
		}
		return err
	}

	updatedOrder := models.Order{
		PaymentStatus: models.OrderPaymentStatusPaid,
	}

	if err := s.db.Where(&models.Order{
		ID:       orderID,
		TenantID: tenantID,
	}).Updates(updatedOrder).Error; err != nil {
		// return errors.New("failed to mark order as paid")
		return err
	}

	//check if student is already enrolled in this course
	if s.db.Where(&models.Enrollment{
		CourseID:  order.CourseID,
		StudentID: order.StudentID,
		TenantID:  tenantID,
	}).First(&models.Enrollment{}).RowsAffected > 0 {
		return nil
	}

	//new enrollment
	newEnrollment := models.Enrollment{
		StudentID: order.StudentID,
		CourseID:  order.CourseID,
		TenantID:  tenantID,
	}

	return s.db.Create(&newEnrollment).Error

}

func (s *orderService) Delete(tenantID uint, orderID uint) error {
	return s.db.Where(&models.Order{
		ID:       orderID,
		TenantID: tenantID,
	}).Delete(&models.Order{}).Error
}
