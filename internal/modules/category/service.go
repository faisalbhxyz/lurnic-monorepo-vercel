package category

import (
	"dashlearn/internal/models"
	"dashlearn/internal/utils"
	"errors"

	"gorm.io/gorm"
)

type CategoryService interface {
	GetAll(tenantID uint) ([]models.Category, error)
	GetByID(tenantID uint, id uint64) (*models.Category, error)
	Create(input CreateCategoryInput, tenantID uint) error
	Update(id uint64, input CreateCategoryInput, tenantID uint) error
	Delete(id uint64, tenantID uint) error
}

type categoryService struct {
	db *gorm.DB
}

func NewCategoryService(db *gorm.DB) CategoryService {
	return &categoryService{
		db: db,
	}
}

func (s *categoryService) GetAll(tenantID uint) ([]models.Category, error) {
	var categories []models.Category
	err := s.db.Where("tenant_id = ?", tenantID).Find(&categories).Error
	return categories, err
}

func (s *categoryService) GetByID(tenantID uint, id uint64) (*models.Category, error) {
	var category models.Category
	result := s.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&category)
	if result.Error != nil {
		return nil, result.Error
	}
	return &category, nil
}

func (s *categoryService) Create(input CreateCategoryInput, tenantID uint) error {
	if s.db.Where("slug = ? AND tenant_id = ?", input.Slug, tenantID).First(&models.Category{}).RowsAffected > 0 {
		return errors.New("category with this slug already exists")
	}

	category := models.Category{
		Name:        input.Name,
		Slug:        input.Slug,
		Description: utils.EmptyStringToNil(input.Description),
		TenantID:    tenantID,
	}

	return s.db.Create(&category).Error
}

func (s *categoryService) Update(id uint64, input CreateCategoryInput, tenantID uint) error {
	var category models.Category

	if err := s.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&category).Error; err != nil {
		return err
	}

	if s.db.Where("slug = ? AND tenant_id = ? AND id != ?", input.Slug, tenantID, id).First(&models.Category{}).RowsAffected > 0 {
		return errors.New("category with this slug already exists")
	}

	category.Name = input.Name
	category.Slug = input.Slug
	category.Description = utils.EmptyStringToNil(input.Description)

	return s.db.Save(&category).Error
}

func (s *categoryService) Delete(id uint64, tenantID uint) error {
	return s.db.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&models.Category{}).Error
}
