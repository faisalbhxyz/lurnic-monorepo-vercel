package banner

import (
	"context"
	"dashlearn/internal/models"
	"dashlearn/internal/response"
	"dashlearn/internal/utils"
	"fmt"

	"gorm.io/gorm"
)

type BannerService interface {
	GetAll(tenantID uint) ([]response.BannerResponse, error)
	GetByID(tenantID uint, id uint64) (*response.BannerResponse, error)
	Create(input CreateBannerInput, tenantID uint) error
	Update(id uint64, input UpdateBannerInput, tenantID uint) error
	Delete(id uint64, tenantID uint) error
}

type bannerService struct {
	db *gorm.DB
}

func NewBannerService(db *gorm.DB) BannerService {
	return &bannerService{
		db: db,
	}
}

func (s *bannerService) GetAll(tenantID uint) ([]response.BannerResponse, error) {
	var banners []response.BannerResponse
	err := s.db.Where("tenant_id = ?", tenantID).Find(&banners).Error
	return banners, err
}

func (s *bannerService) GetByID(tenantID uint, id uint64) (*response.BannerResponse, error) {
	var banner response.BannerResponse
	if err := s.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&banner).Error; err != nil {
		return nil, err
	}
	return &banner, nil
}

func (s *bannerService) Create(input CreateBannerInput, tenantID uint) error {
	banner := models.Banner{
		Title:    utils.ZeroToNil(input.Title),
		Url:      utils.ZeroToNil(input.Url),
		Image:    input.Image,
		TenantID: tenantID,
	}

	return s.db.Create(&banner).Error
}

func (s *bannerService) Update(id uint64, input UpdateBannerInput, tenantID uint) error {
	var banner response.BannerResponse

	if err := s.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&banner).Error; err != nil {
		return err
	}

	if input.Image != nil {
		err := utils.DeleteCDNFile(context.Background(), banner.Image)
		if err != nil {
			fmt.Println("Error deleting old image:", err)
		}
	}

	updatedbanner := models.Banner{
		Title: utils.ZeroToNil(input.Title),
		Url:   utils.ZeroToNil(input.Url),
		Image: func() string {
			if input.Image != nil && *input.Image != "" {
				return *input.Image
			}
			return banner.Image
		}(),
	}

	return s.db.Where(models.Banner{ID: uint(id), TenantID: tenantID}).Updates(&updatedbanner).Error
}

func (s *bannerService) Delete(id uint64, tenantID uint) error {
	var banner models.Banner

	if err := s.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&banner).Error; err != nil {
		return err
	}

	if err := utils.DeleteCDNFile(context.Background(), banner.Image); err != nil {
		fmt.Println("Error deleting image:", err)
		return err
	}

	return s.db.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&models.Banner{}).Error
}
