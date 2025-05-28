package utils

import (
	"testing"
)

const testSecretKey = "secret"

func TestGenerateToken(t *testing.T) {
	email := "test@example.com"
	userID := "user123"

	tokenString, err := GenerateToken(email, userID)

	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	if tokenString == "" {
		t.Error("Generated token string is empty")
	}

	gotUserID, err := VerifyToken(tokenString)
	if err != nil {
		t.Fatalf("VerifyToken failed to validate generated token: %v", err)
	}

	if gotUserID != userID {
		t.Errorf("Expected userID %q, got %q", userID, gotUserID)
	}
}
