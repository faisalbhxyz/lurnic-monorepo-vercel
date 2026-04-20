package role

import (
	"time"

	"gorm.io/datatypes"
)

type RoleRes struct {
	ID          uint           `json:"id"`
	Name        string         `json:"name"`
	Permissions datatypes.JSON `json:"permissions"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type UserRoleRes struct {
	UserID      string         `json:"user_id"`
	Name        string         `json:"name"`
	Role        string         `json:"role"`
	Permissions datatypes.JSON `json:"permissions"`
}
