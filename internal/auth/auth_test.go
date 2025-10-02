package auth

import (
	"strings"
	"testing"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  false, // argon2id allows empty passwords
		},
		{
			name:     "long password",
			password: strings.Repeat("a", 1000),
			wantErr:  false,
		},
		{
			name:     "password with special characters",
			password: "p@ssw0rd!#$%^&*()",
			wantErr:  false,
		},
		{
			name:     "unicode password",
			password: "пароль123",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)

			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Hash should not be empty
				if hash == "" {
					t.Error("HashPassword() returned empty hash")
				}

				// Hash should be different from the original password
				if hash == tt.password {
					t.Error("HashPassword() returned the original password as hash")
				}

				// Hash should start with argon2id prefix
				if !strings.HasPrefix(hash, "$argon2id$") {
					t.Error("HashPassword() returned hash without proper argon2id prefix")
				}
			}
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	// First, create some test hashes
	password1 := "password123"
	password2 := "different_password"

	hash1, err := HashPassword(password1)
	if err != nil {
		t.Fatalf("Failed to create test hash: %v", err)
	}

	hash2, err := HashPassword(password2)
	if err != nil {
		t.Fatalf("Failed to create test hash: %v", err)
	}

	tests := []struct {
		name     string
		password string
		hash     string
		want     bool
		wantErr  bool
	}{
		{
			name:     "correct password and hash",
			password: password1,
			hash:     hash1,
			want:     true,
			wantErr:  false,
		},
		{
			name:     "incorrect password",
			password: "wrong_password",
			hash:     hash1,
			want:     false,
			wantErr:  false,
		},
		{
			name:     "different correct password and hash",
			password: password2,
			hash:     hash2,
			want:     true,
			wantErr:  false,
		},
		{
			name:     "password for different hash",
			password: password1,
			hash:     hash2,
			want:     false,
			wantErr:  false,
		},
		{
			name:     "empty password with valid hash",
			password: "",
			hash:     hash1,
			want:     false,
			wantErr:  false,
		},
		{
			name:     "invalid hash format",
			password: password1,
			hash:     "invalid_hash",
			want:     false,
			wantErr:  true,
		},
		{
			name:     "empty hash",
			password: password1,
			hash:     "",
			want:     false,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckPasswordHash(tt.password, tt.hash)

			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("CheckPasswordHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHashPasswordConsistency(t *testing.T) {
	password := "test_password"

	// Generate multiple hashes for the same password
	hash1, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	hash2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	// Hashes should be different (due to salt)
	if hash1 == hash2 {
		t.Error("HashPassword() generated identical hashes for the same password")
	}

	// But both should validate against the original password
	match1, err := CheckPasswordHash(password, hash1)
	if err != nil {
		t.Fatalf("CheckPasswordHash() error = %v", err)
	}
	if !match1 {
		t.Error("CheckPasswordHash() failed to validate first hash")
	}

	match2, err := CheckPasswordHash(password, hash2)
	if err != nil {
		t.Fatalf("CheckPasswordHash() error = %v", err)
	}
	if !match2 {
		t.Error("CheckPasswordHash() failed to validate second hash")
	}
}

func TestEmptyPasswordHash(t *testing.T) {
	// Test that empty password can be hashed and verified
	emptyPassword := ""

	hash, err := HashPassword(emptyPassword)
	if err != nil {
		t.Fatalf("HashPassword() with empty password failed: %v", err)
	}

	match, err := CheckPasswordHash(emptyPassword, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash() with empty password failed: %v", err)
	}

	if !match {
		t.Error("CheckPasswordHash() failed to validate empty password")
	}

	// Verify that non-empty password doesn't match empty password hash
	match, err = CheckPasswordHash("not_empty", hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash() error = %v", err)
	}

	if match {
		t.Error("CheckPasswordHash() incorrectly validated non-empty password against empty password hash")
	}
}
