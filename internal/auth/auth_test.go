package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
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

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test_secret"
	expiresIn := time.Hour

	tests := []struct {
		name      string
		userID    uuid.UUID
		secret    string
		expiresIn time.Duration
		wantErr   bool
	}{
		{
			name:      "valid JWT creation",
			userID:    userID,
			secret:    secret,
			expiresIn: expiresIn,
			wantErr:   false,
		},
		{
			name:      "empty secret",
			userID:    userID,
			secret:    "",
			expiresIn: expiresIn,
			wantErr:   false, // Empty secret is technically valid for HMAC
		},
		{
			name:      "zero expiration",
			userID:    userID,
			secret:    secret,
			expiresIn: 0,
			wantErr:   false,
		},
		{
			name:      "negative expiration",
			userID:    userID,
			secret:    secret,
			expiresIn: -time.Hour,
			wantErr:   false, // Should create expired token
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenString, err := MakeJWT(tt.userID, tt.secret, tt.expiresIn)

			if (err != nil) != tt.wantErr {
				t.Errorf("MakeJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if tokenString == "" {
					t.Error("MakeJWT() returned empty token string")
				}

				// Token should have 3 parts separated by dots
				parts := strings.Split(tokenString, ".")
				if len(parts) != 3 {
					t.Errorf("MakeJWT() returned token with %d parts, want 3", len(parts))
				}
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test_secret"
	wrongSecret := "wrong_secret"

	// Create a valid token
	validToken, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}

	// Create an expired token
	expiredToken, err := MakeJWT(userID, secret, -time.Hour)
	if err != nil {
		t.Fatalf("Failed to create expired token: %v", err)
	}

	// Create a token with wrong secret
	wrongSecretToken, err := MakeJWT(userID, wrongSecret, time.Hour)
	if err != nil {
		t.Fatalf("Failed to create wrong secret token: %v", err)
	}

	tests := []struct {
		name       string
		token      string
		secret     string
		expectedID uuid.UUID
		wantErr    bool
	}{
		{
			name:       "valid token",
			token:      validToken,
			secret:     secret,
			expectedID: userID,
			wantErr:    false,
		},
		{
			name:       "expired token",
			token:      expiredToken,
			secret:     secret,
			expectedID: uuid.UUID{},
			wantErr:    true,
		},
		{
			name:       "wrong secret",
			token:      wrongSecretToken,
			secret:     secret,
			expectedID: uuid.UUID{},
			wantErr:    true,
		},
		{
			name:       "malformed token",
			token:      "invalid.token.format",
			secret:     secret,
			expectedID: uuid.UUID{},
			wantErr:    true,
		},
		{
			name:       "empty token",
			token:      "",
			secret:     secret,
			expectedID: uuid.UUID{},
			wantErr:    true,
		},
		{
			name:       "valid token with wrong secret",
			token:      validToken,
			secret:     wrongSecret,
			expectedID: uuid.UUID{},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, err := ValidateJWT(tt.token, tt.secret)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && gotID != tt.expectedID {
				t.Errorf("ValidateJWT() = %v, want %v", gotID, tt.expectedID)
			}
		})
	}
}

func TestJWTRoundTrip(t *testing.T) {
	userID := uuid.New()
	secret := "round_trip_secret"
	expiresIn := time.Hour

	// Create a JWT
	tokenString, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT() failed: %v", err)
	}

	// Validate the JWT
	retrievedUserID, err := ValidateJWT(tokenString, secret)
	if err != nil {
		t.Fatalf("ValidateJWT() failed: %v", err)
	}

	// Should get back the same user ID
	if retrievedUserID != userID {
		t.Errorf("Round trip failed: got %v, want %v", retrievedUserID, userID)
	}
}

func TestJWTExpiration(t *testing.T) {
	userID := uuid.New()
	secret := "expiration_test_secret"

	// Test 1: Create a valid token
	validToken, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT() failed for valid token: %v", err)
	}

	// Should be valid
	_, err = ValidateJWT(validToken, secret)
	if err != nil {
		t.Errorf("ValidateJWT() failed for valid token: %v", err)
	}

	// Test 2: Create an already expired token
	expiredToken, err := MakeJWT(userID, secret, -time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT() failed for expired token: %v", err)
	}

	// Should be expired
	_, err = ValidateJWT(expiredToken, secret)
	if err == nil {
		t.Error("ValidateJWT() should have failed for expired token")
	}
}

func TestJWTDifferentSecrets(t *testing.T) {
	userID := uuid.New()
	secret1 := "secret_one"
	secret2 := "secret_two"
	expiresIn := time.Hour

	// Create token with secret1
	tokenString, err := MakeJWT(userID, secret1, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT() failed: %v", err)
	}

	// Should validate with secret1
	_, err = ValidateJWT(tokenString, secret1)
	if err != nil {
		t.Errorf("ValidateJWT() failed with correct secret: %v", err)
	}

	// Should fail with secret2
	_, err = ValidateJWT(tokenString, secret2)
	if err == nil {
		t.Error("ValidateJWT() should have failed with wrong secret")
	}
}

func TestJWTWithDifferentUserIDs(t *testing.T) {
	userID1 := uuid.New()
	userID2 := uuid.New()
	secret := "user_id_test_secret"
	expiresIn := time.Hour

	// Create tokens for different users
	token1, err := MakeJWT(userID1, secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT() failed for user 1: %v", err)
	}

	token2, err := MakeJWT(userID2, secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT() failed for user 2: %v", err)
	}

	// Validate first token returns first user ID
	retrievedID1, err := ValidateJWT(token1, secret)
	if err != nil {
		t.Fatalf("ValidateJWT() failed for token 1: %v", err)
	}
	if retrievedID1 != userID1 {
		t.Errorf("ValidateJWT() returned %v for token 1, want %v", retrievedID1, userID1)
	}

	// Validate second token returns second user ID
	retrievedID2, err := ValidateJWT(token2, secret)
	if err != nil {
		t.Fatalf("ValidateJWT() failed for token 2: %v", err)
	}
	if retrievedID2 != userID2 {
		t.Errorf("ValidateJWT() returned %v for token 2, want %v", retrievedID2, userID2)
	}

	// Tokens should be different
	if token1 == token2 {
		t.Error("Different user IDs generated identical tokens")
	}
}
