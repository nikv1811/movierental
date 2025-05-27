package utils

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "mySecretPassword123"
	hashedPassword, err := HashPassword(password)

	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if hashedPassword == "" {
		t.Error("Hashed password is empty")
	}

	// Verify that the hashed password is valid using CheckPasswordHash
	if !CheckPasswordHash(password, hashedPassword) {
		t.Error("CheckPasswordHash failed to verify a newly hashed password")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "testPassword"
	// Manually hash a password for consistent testing (or use the helper)
	hashedPassword, _ := HashPassword(password) // Assume hashing works for this test

	tests := []struct {
		name           string
		password       string
		hashedPassword string
		expected       bool
	}{
		{
			name:           "Valid password",
			password:       password,
			hashedPassword: hashedPassword,
			expected:       true,
		},
		{
			name:           "Invalid password",
			password:       "wrongPassword",
			hashedPassword: hashedPassword,
			expected:       false,
		},
		{
			name:           "Empty password",
			password:       "",
			hashedPassword: hashedPassword,
			expected:       false,
		},
		{
			name:           "Empty hashed password",
			password:       password,
			hashedPassword: "", // bcrypt.CompareHashAndPassword will return error for empty hash
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckPasswordHash(tt.password, tt.hashedPassword)
			if got != tt.expected {
				t.Errorf("CheckPasswordHash(%q, %q) = %v; want %v", tt.password, tt.hashedPassword, got, tt.expected)
			}
		})
	}
}
