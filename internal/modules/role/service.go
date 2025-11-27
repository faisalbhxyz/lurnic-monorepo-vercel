package role

import (
	"dashlearn/internal/models"
	"errors"

	"gorm.io/gorm"
)

type RoleService interface {
	GetRoles(tenantID uint) (*[]RoleRes, error)
	GetRoleByUserID(auth_user_id uint, tenant_id uint) (*UserRoleRes, error)
	CreateRole(input CreateRoleInput, tenantID uint) error
	DeleteRole(roleID uint, tenantID uint) error
}

type roleService struct {
	db *gorm.DB
}

func NewRoleService(db *gorm.DB) RoleService {
	return &roleService{db: db}
}

func (s *roleService) GetRoles(tenantID uint) (*[]RoleRes, error) {
	var roles []models.Role
	var roleRes []RoleRes

	if err := s.db.Where("tenant_id = ? AND name != ?", tenantID, "superadmin").Find(&roles).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("role not found")
		}
		return nil, errors.New("something went wrong, please try again")
	}

	for _, role := range roles {
		roleRes = append(roleRes, RoleRes{
			ID:          role.ID,
			Name:        role.Name,
			Permissions: role.Permissions,
			CreatedAt:   role.CreatedAt,
			UpdatedAt:   role.UpdatedAt,
		})
	}

	return &roleRes, nil
}

func (s *roleService) GetRoleByUserID(auth_user_id uint, tenant_id uint) (*UserRoleRes, error) {
	var user models.User
	var roleRes UserRoleRes

	if err := s.db.Preload("Role").Where("tenant_id = ? AND id = ?", tenant_id, auth_user_id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("something went wrong, please try again")
	}

	roleRes = UserRoleRes{
		UserID: user.UserID,
		Name:   user.Name,
		Role:   user.Role.Name,
		// Permissions: user.Role.Permissions,
	}

	return &roleRes, nil
}

func (s *roleService) CreateRole(input CreateRoleInput, tenantID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Check if role already exists for this tenant
		var existingRole models.Role
		if err := tx.Where("tenant_id = ? AND name = ?", tenantID, input.Title).First(&existingRole).Error; err == nil {
			return errors.New("role already exists")
		} else if err != gorm.ErrRecordNotFound {
			return errors.New("something went wrong, please try again")
		}

		// Create new role
		newRole := models.Role{
			Name:        input.Title,
			Permissions: input.Permissions,
			TenantID:    &tenantID,
		}

		if err := tx.Create(&newRole).Error; err != nil {
			return errors.New("role not created")
		}

		return nil
	})
}

func (s *roleService) DeleteRole(roleID uint, tenantID uint) error {
	var role models.Role
	if err := s.db.Where("id = ? AND tenant_id = ?", roleID, tenantID).First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("role not found")
		}
		return errors.New("something went wrong, please try again")
	}

	if err := s.db.Delete(&role).Error; err != nil {
		return errors.New("role not deleted")
	}
	return nil
}
