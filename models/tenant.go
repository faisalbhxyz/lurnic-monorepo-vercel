package models

import (
	"time"
)

type Tenant struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	AppKey    string    `gorm:"uniqueIndex" json:"app_key"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
