package main

import (
	"dashlearn/internal/models"
	"dashlearn/internal/utils"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/lucsky/cuid"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	_ = godotenv.Load()

	if err := utils.ConnectDatabase(); err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	appKey := envOrDefault("SEED_APP_KEY", "local-dev")
	email := envOrDefault("SEED_EMAIL", "admin@local.dev")
	password := envOrDefault("SEED_PASSWORD", "password123")
	name := envOrDefault("SEED_NAME", "Local Admin")

	tenantID, err := ensureTenant(appKey)
	if err != nil {
		log.Fatalf("failed to ensure tenant: %v", err)
	}

	// Prefer global superadmin role if present; else create a tenant-scoped admin role.
	roleID, err := ensureRole(tenantID)
	if err != nil {
		log.Fatalf("failed to ensure role: %v", err)
	}

	if err := upsertUser(tenantID, roleID, name, email, password); err != nil {
		log.Fatalf("failed to seed user: %v", err)
	}

	fmt.Println("Seed complete.")
	fmt.Printf("Tenant app-key: %s\n", appKey)
	fmt.Printf("Login email: %s\n", email)
	fmt.Printf("Login password: %s\n", password)
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func ensureTenant(appKey string) (uint, error) {
	var tenant models.Tenant
	if err := utils.DB.Where("app_key = ?", appKey).First(&tenant).Error; err == nil {
		return tenant.ID, nil
	}

	tenant = models.Tenant{AppKey: appKey}
	if err := utils.DB.Create(&tenant).Error; err != nil {
		return 0, err
	}
	return tenant.ID, nil
}

func ensureRole(tenantID uint) (uint, error) {
	// global superadmin
	var superadmin models.Role
	if err := utils.DB.Where("name = ? AND tenant_id IS NULL", "superadmin").First(&superadmin).Error; err == nil {
		return superadmin.ID, nil
	}

	tenantRoleName := "admin"
	var role models.Role
	if err := utils.DB.Where("name = ? AND tenant_id = ?", tenantRoleName, tenantID).First(&role).Error; err == nil {
		return role.ID, nil
	}

	role = models.Role{
		Name:     tenantRoleName,
		TenantID: &tenantID,
	}
	if err := utils.DB.Create(&role).Error; err != nil {
		return 0, err
	}
	return role.ID, nil
}

func upsertUser(tenantID, roleID uint, name, email, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	var user models.User
	if err := utils.DB.Where("email = ? AND tenant_id = ?", email, tenantID).First(&user).Error; err == nil {
		user.Name = name
		user.Password = string(hashed)
		user.Status = true
		rid := roleID
		user.RoleID = &rid
		return utils.DB.Save(&user).Error
	}

	rid := roleID
	user = models.User{
		UserID:   cuid.New(),
		Name:     name,
		Email:    email,
		Password: string(hashed),
		Status:   true,
		TenantID: tenantID,
		RoleID:   &rid,
	}
	return utils.DB.Create(&user).Error
}

