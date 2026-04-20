package role

import "gorm.io/datatypes"

type CreateRoleInput struct {
	Title       string         `json:"title" binding:"required"`
	Permissions datatypes.JSON `json:"permissions" binding:"required"`
}
