package utils

import (
	"dashlearn/models"
	"strings"
)

func CreateSuperadminIfNotExists() {
	var r models.Role
	err := DB.FirstOrCreate(&r, models.Role{Name: "superadmin", TenantID: nil}).Error
	if err != nil {
		panic("Failed to create or find superadmin role: " + err.Error())
	}
}

func EmptyStringToNil(s *string) *string {
	if s != nil && strings.TrimSpace(*s) == "" {
		return nil
	}
	return s
}
