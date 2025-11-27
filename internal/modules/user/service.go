package user

import (
	"dashlearn/internal/models"
	"dashlearn/internal/modules/role"
	"dashlearn/internal/utils"
	"errors"

	"github.com/lucsky/cuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	GetUsers() []UserRes
	CreateUser(input CreateUserInput) error
	LoginUser(input LoginUserInput) (*LoginUserRes, error)
	CreateTeamMember(input CreateTeamMemberInput, tenantID uint) error
	GetTeamMembers(tenant_id uint) []TeamMemberRes
	GetTeamMemberByUID(tenant_id uint, user_id uint) TeamMemberRes
	UpdateTeamMember(tenant_id uint, user_id uint, input UpdateTeamMemberInput) error
	DeleteTeamMember(user_id uint, tenant_id uint) error
}

type userService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) UserService {
	return &userService{db: db}
}

func (s *userService) GetUsers() []UserRes {
	var users []models.User
	var usersRes []UserRes
	s.db.
		Select("user_id", "name", "phone", "email", "status", "created_at", "updated_at").
		Find(&users)

	for _, user := range users {
		usersRes = append(usersRes, UserRes{
			UserID:    user.UserID,
			Name:      user.Name,
			Phone:     user.Phone,
			Email:     user.Email,
			Status:    user.Status,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}
	return usersRes
}

func (s *userService) CreateUser(input CreateUserInput) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return errors.New("something went wrong, please try again")
		}

		var existingUser models.User
		if err := utils.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
			return errors.New("email already exists")
		} else if err != gorm.ErrRecordNotFound {
			return errors.New("something went wrong. Please try again")
		}

		// newTenant := models.Tenant{
		// 	AppKey: cuid.New(),
		// }

		// if err := tx.Create(&newTenant).Error; err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong. Please try again."})
		// 	return
		// }

		// var superadmin models.Role
		// if err := tx.Where("name = ?", "superadmin").First(&superadmin).Error; err != nil {
		// 	return errors.New("something went wrong, please try again")
		// }

		// var tenant models.Tenant
		// if err := tx.Where("email = ?", "john@gmail.com").First(&tenant).Error; err != nil {
		// 	return errors.New("something went wrong, please try again")
		// }

		newUser := models.User{
			UserID:   cuid.New(),
			Name:     input.Name,
			Phone:    input.Phone,
			Email:    input.Email,
			Password: string(hashedPassword),
			Status:   true,
			// TenantID: tenant.ID,
		}

		if err := tx.Create(&newUser).Error; err != nil {
			return errors.New("failed to create user")
		}

		return nil
	})
}

func (s *userService) LoginUser(input LoginUserInput) (*LoginUserRes, error) {
	var user models.User
	err := s.db.Where("email = ?", input.Email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, errors.New("something went wrong, please try again")
	}

	// Compare the provided password with the stored hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, errors.New("invalid password")
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.UserID)
	if err != nil {
		return nil, errors.New("something went wrong, please try again")
	}

	return &LoginUserRes{
		Token: token,
		User: UserResLite{
			UserID: user.UserID,
			Name:   user.Name,
			Phone:  user.Phone,
			Email:  user.Email,
		},
	}, nil

}

func (s *userService) CreateTeamMember(input CreateTeamMemberInput, tenantID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return errors.New("something went wrong. Please try again")
		}

		var existingUser models.User
		if err := tx.Where("user_id = ? OR email = ?", input.UserID, input.Email).First(&existingUser).Error; err == nil {
			if existingUser.UserID == input.UserID {
				return errors.New("user with this user id already exists")
			}
			if existingUser.Email == input.Email {
				return errors.New("email already exists")
			}
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("something went wrong, please try again")
		}

		roleID := uint(input.Role)
		newUser := models.User{
			UserID:   input.UserID,
			Name:     input.Name,
			Phone:    utils.ZeroToNil(input.Phone),
			RoleID:   &roleID,
			Email:    input.Email,
			Password: string(hashedPassword),
			Status:   true,
			TenantID: tenantID,
		}

		if err := utils.DB.Create(&newUser).Error; err != nil {
			return errors.New("failed to create user")
		}

		return nil
	})
}

func (s *userService) UpdateTeamMember(tenant_id uint, user_id uint, input UpdateTeamMemberInput) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var user models.User
		err := tx.Where("id = ? AND tenant_id = ?", user_id, tenant_id).First(&user).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("user not found")
			}
			return err
		}

		updates := map[string]any{}

		if input.Name != nil {
			updates["name"] = input.Name
		}
		if input.Phone != nil {
			updates["phone"] = input.Phone
		}
		if input.Role != nil {
			updates["role_id"] = input.Role
		}

		if len(updates) == 0 {
			return nil
		}

		if err := tx.Model(&user).Updates(updates).Error; err != nil {
			return errors.New("failed to update user")
		}

		return nil
	})
}

func (s *userService) GetTeamMemberByUID(tenant_id uint, user_id uint) TeamMemberRes {
	var user models.User
	var usersRes TeamMemberRes
	s.db.
		Preload("Role").
		Where("role_id != ? AND id = ? AND tenant_id = ?", 1, user_id, tenant_id).
		First(&user)

	usersRes = TeamMemberRes{
		UserID: user.UserID,
		Name:   user.Name,
		Phone:  user.Phone,
		Email:  user.Email,
		Role: role.RoleRes{
			ID:          user.Role.ID,
			Name:        user.Role.Name,
			Permissions: user.Role.Permissions,
		},
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return usersRes
}

func (s *userService) GetTeamMembers(tenant_id uint) []TeamMemberRes {
	var users []models.User
	var usersRes []TeamMemberRes
	s.db.
		Preload("Role").
		Where("role_id != ?", 1).
		Find(&users)

	for _, user := range users {
		usersRes = append(usersRes, TeamMemberRes{
			ID:     user.ID,
			UserID: user.UserID,
			Name:   user.Name,
			Phone:  user.Phone,
			Email:  user.Email,
			Role: role.RoleRes{
				ID:          user.Role.ID,
				Name:        user.Role.Name,
				Permissions: user.Role.Permissions,
			},
			Status:    user.Status,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}
	return usersRes
}

func (s *userService) DeleteTeamMember(user_id uint, tenant_id uint) error {
	var user models.User
	err := s.db.Where("id = ? AND tenant_id = ?", user_id, tenant_id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("user not found")
	} else if err != nil {
		return errors.New("something went wrong, please try again")
	}

	if err := s.db.Delete(&user).Error; err != nil {
		return errors.New("user not deleted")
	}
	return nil
}
