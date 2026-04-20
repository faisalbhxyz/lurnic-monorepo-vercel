package subcategory

import (
	"dashlearn/internal/models"
	"dashlearn/internal/response"
	"dashlearn/internal/utils"
	"errors"

	"gorm.io/gorm"
)

type SubCategoryService interface {
	GetAll(tenantID uint) ([]response.SubCategoryResponse, error)
	GetByID(tenantID uint, id uint64) (*response.SubCategoryResponse, error)
	Create(input CreateSubCategoryInput, tenantID uint) error
	Update(id uint64, input CreateSubCategoryInput, tenantID uint) error
	Delete(id uint64, tenantID uint) error
}

type subcategoryService struct {
	db *gorm.DB
}

func NewSubCategoryService(db *gorm.DB) SubCategoryService {
	return &subcategoryService{
		db: db,
	}
}

func (s *subcategoryService) GetAll(tenantID uint) ([]response.SubCategoryResponse, error) {
	var subcategories []models.SubCategory
	err := s.db.Where("tenant_id = ?", tenantID).Preload("Category").Find(&subcategories).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	result := make([]response.SubCategoryResponse, len(subcategories))

	for i, subcategory := range subcategories {
		result[i] = response.SubCategoryResponse{
			ID:          subcategory.ID,
			Name:        subcategory.Name,
			Slug:        subcategory.Slug,
			Description: subcategory.Description,
			Category: &response.CategoryResponse{
				ID:   subcategory.Category.ID,
				Name: subcategory.Category.Name,
				Slug: subcategory.Category.Slug,
			},
		}
	}

	return result, nil
}

func (s *subcategoryService) GetByID(tenantID uint, id uint64) (*response.SubCategoryResponse, error) {
	var subcategory models.SubCategory
	result := s.db.Where("id = ? AND tenant_id = ?", id, tenantID).Preload("Category").First(&subcategory)

	if result.Error != nil {
		return nil, result.Error
	}

	response := response.SubCategoryResponse{
		ID:          subcategory.ID,
		Name:        subcategory.Name,
		Slug:        subcategory.Slug,
		Description: subcategory.Description,
		Category: &response.CategoryResponse{
			ID:   subcategory.Category.ID,
			Name: subcategory.Category.Name,
			Slug: subcategory.Category.Slug,
		},
	}

	return &response, nil
}

func (s *subcategoryService) Create(input CreateSubCategoryInput, tenantID uint) error {
	if s.db.Where("slug = ? AND tenant_id = ?", input.Slug, tenantID).First(&models.SubCategory{}).RowsAffected > 0 {
		return errors.New("sub category with this slug already exists")
	}

	subcategory := models.SubCategory{
		CategoryID:  input.CategoryID,
		Name:        input.Name,
		Slug:        input.Slug,
		Description: utils.EmptyStringToNil(input.Description),
		TenantID:    tenantID,
	}

	return s.db.Create(&subcategory).Error
}

func (s *subcategoryService) Update(id uint64, input CreateSubCategoryInput, tenantID uint) error {
	var category models.SubCategory

	if err := s.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&category).Error; err != nil {
		return err
	}

	if s.db.Where("slug = ? AND tenant_id = ? AND id != ?", input.Slug, tenantID, id).First(&models.SubCategory{}).RowsAffected > 0 {
		return errors.New("sub category with this slug already exists")
	}

	category.CategoryID = input.CategoryID
	category.Name = input.Name
	category.Slug = input.Slug
	category.Description = utils.EmptyStringToNil(input.Description)

	return s.db.Save(&category).Error
}

func (s *subcategoryService) Delete(id uint64, tenantID uint) error {
	return s.db.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&models.SubCategory{}).Error
}
