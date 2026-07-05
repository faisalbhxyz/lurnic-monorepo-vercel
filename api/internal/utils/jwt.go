package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

const passwordResetPurpose = "student_password_reset"

func GenerateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func GeneratePasswordResetJWT(userID, email string, tenantID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   userID,
		"email":     email,
		"tenant_id": tenantID,
		"purpose":   passwordResetPurpose,
		"exp":       time.Now().Add(1 * time.Hour).Unix(),
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

type PasswordResetClaims struct {
	UserID   string
	Email    string
	TenantID uint
}

func ParsePasswordResetJWT(tokenStr string) (*PasswordResetClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired reset token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid reset token claims")
	}

	if purpose, _ := claims["purpose"].(string); purpose != passwordResetPurpose {
		return nil, errors.New("invalid reset token purpose")
	}

	userID, _ := claims["user_id"].(string)
	email, _ := claims["email"].(string)
	tenantIDFloat, _ := claims["tenant_id"].(float64)
	if userID == "" || email == "" || tenantIDFloat == 0 {
		return nil, errors.New("invalid reset token claims")
	}

	if exp, ok := claims["exp"].(float64); ok && int64(exp) < time.Now().Unix() {
		return nil, errors.New("reset token expired")
	}

	return &PasswordResetClaims{
		UserID:   userID,
		Email:    email,
		TenantID: uint(tenantIDFloat),
	}, nil
}
