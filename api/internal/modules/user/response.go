package user

import (
	"dashlearn/internal/modules/role"
	"time"
)

type UserRes struct {
	ID        uint      `json:"id,omitempty"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Phone     *string   `json:"phone"`
	Email     string    `json:"email"`
	Status    bool      `json:"status,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	RoleID    *uint     `json:"role_id,omitempty"`
	// Role      Role      `gorm:"foreignKey:RoleID;references:ID" json:"role,omitempty"`
}

type UserResLite struct {
	ID     uint    `json:"id,omitempty"`
	UserID string  `json:"user_id"`
	Name   string  `json:"name"`
	Phone  *string `json:"phone"`
	Email  string  `json:"email"`
}

type LoginUserRes struct {
	Token string      `json:"token"`
	User  UserResLite `json:"user"`
}

type TeamMemberRes struct {
	ID        uint         `json:"id,omitempty"`
	UserID    string       `json:"user_id"`
	Name      string       `json:"name"`
	Phone     *string      `json:"phone"`
	Email     string       `json:"email"`
	Status    bool         `json:"status,omitempty"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	RoleID    uint         `json:"role_id"`
	Role      role.RoleRes `json:"role"`
}
