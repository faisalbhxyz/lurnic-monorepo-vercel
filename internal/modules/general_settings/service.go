package generalsettings

import (
	"context"
	"dashlearn/internal/models"
	"dashlearn/internal/response"
	"dashlearn/internal/utils"
	"fmt"

	"gorm.io/gorm"
)

type GeneralSettingsService interface {
	GetGeneralSettings(tenantID uint) (response.GeneralSettingsRes, error)
	UpdateGeneralSettings(input *CreateOrUpdateGeneralSettingsInput, tenantID uint) error
}

type generalsettingsService struct {
	db *gorm.DB
}

func NewGeneralSettingsService(db *gorm.DB) GeneralSettingsService {
	return &generalsettingsService{
		db: db,
	}
}

func (s *generalsettingsService) GetGeneralSettings(tenantID uint) (response.GeneralSettingsRes, error) {
	var res response.GeneralSettingsRes
	if err := s.db.Model(&models.GeneralSettings{}).
		Where("tenant_id = ?", tenantID).
		First(&res).Error; err != nil {
		return response.GeneralSettingsRes{}, err
	}
	return res, nil
}

func (s *generalsettingsService) UpdateGeneralSettings(input *CreateOrUpdateGeneralSettingsInput, tenantID uint) error {
	var existingRes response.GeneralSettingsRes

	if err := s.db.Model(&models.GeneralSettings{}).
		Where("tenant_id = ?", tenantID).
		First(&existingRes).Error; err != nil {
		//create new record
		newSettings := models.GeneralSettings{
			OrgName:       input.OrgName,
			Logo:          utils.ZeroToNil(input.Logo),
			Favicon:       utils.ZeroToNil(input.Favicon),
			StudentPrefix: input.StudentPrefix,
			TeacherPrefix: input.TeacherPrefix,
			TenantID:      tenantID,
		}

		return s.db.Create(&newSettings).Error
	} else {
		//update record

		if input.Logo != nil && *input.Logo != "" && existingRes.Logo != nil && *existingRes.Logo != "" {
			err := utils.DeleteCDNFile(context.Background(), *existingRes.Logo)
			if err != nil {
				fmt.Println("Error deleting old image:", err)
			}
		} else if input.Logo != nil && *input.Logo != "" && existingRes.Logo == nil {
			input.Logo = utils.ZeroToNil(input.Logo)
		} else {
			input.Logo = existingRes.Logo
		}

		if input.Favicon != nil && *input.Favicon != "" && existingRes.Favicon != nil && *existingRes.Favicon != "" {
			err := utils.DeleteCDNFile(context.Background(), *existingRes.Favicon)
			if err != nil {
				fmt.Println("Error deleting old image:", err)
			}
		} else if input.Favicon != nil && *input.Favicon != "" && existingRes.Favicon == nil {
			input.Favicon = utils.ZeroToNil(input.Favicon)
		} else {
			input.Favicon = existingRes.Favicon
		}

		updatedSettings := models.GeneralSettings{
			OrgName:       input.OrgName,
			Logo:          input.Logo,
			Favicon:       input.Favicon,
			StudentPrefix: input.StudentPrefix,
			TeacherPrefix: input.TeacherPrefix,
		}

		return s.db.Where("tenant_id = ?", tenantID).Updates(&updatedSettings).Error
	}

}
