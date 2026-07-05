package utils

import (
	"os"
	"testing"
)

func TestPasswordResetJWTRoundTrip(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-for-password-reset")
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))

	token, err := GeneratePasswordResetJWT("user_abc", "student@example.com", 7)
	if err != nil {
		t.Fatalf("GeneratePasswordResetJWT: %v", err)
	}

	claims, err := ParsePasswordResetJWT(token)
	if err != nil {
		t.Fatalf("ParsePasswordResetJWT: %v", err)
	}

	if claims.UserID != "user_abc" || claims.Email != "student@example.com" || claims.TenantID != 7 {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

func TestParsePasswordResetJWTRejectsLoginToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-for-password-reset")
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))

	loginToken, err := GenerateJWT("user_abc")
	if err != nil {
		t.Fatalf("GenerateJWT: %v", err)
	}

	if _, err := ParsePasswordResetJWT(loginToken); err == nil {
		t.Fatal("expected login token to be rejected")
	}
}
