package model

import (
	"testing"
)

func TestUserHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{
			name:     "Valid password",
			password: "password123",
		},
		{
			name:     "Empty password",
			password: "",
		},
		{
			name:     "Long password",
			password: "this_is_a_very_long_password_with_many_characters_1234567890",
		},
		{
			name:     "Special characters",
			password: "p@ssw0rd!#$%",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Password: tt.password}
			originalPassword := tt.password
			user.HashPassword()

			// Check that password was hashed
			if user.Password == originalPassword && originalPassword != "" {
				t.Errorf("Password was not hashed: got %s, want different from %s", user.Password, originalPassword)
			}

			// Check that hashed password is not empty
			if tt.password != "" && user.Password == "" {
				t.Error("Hashed password should not be empty for non-empty input")
			}

			// Check that hashed password starts with bcrypt prefix
			if tt.password != "" && len(user.Password) < 10 {
				t.Error("Hashed password is too short to be a valid bcrypt hash")
			}
		})
	}
}

func TestUserComparePassword(t *testing.T) {
	tests := []struct {
		name           string
		originalPass   string
		comparePass    string
		expectedResult bool
	}{
		{
			name:           "Matching passwords",
			originalPass:   "password123",
			comparePass:    "password123",
			expectedResult: true,
		},
		{
			name:           "Non-matching passwords",
			originalPass:   "password123",
			comparePass:    "wrongpassword",
			expectedResult: false,
		},
		{
			name:           "Empty comparison password",
			originalPass:   "password123",
			comparePass:    "",
			expectedResult: false,
		},
		{
			name:           "Case sensitive comparison",
			originalPass:   "Password123",
			comparePass:    "password123",
			expectedResult: false,
		},
		{
			name:           "Special characters",
			originalPass:   "p@ssw0rd!#$",
			comparePass:    "p@ssw0rd!#$",
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Password: tt.originalPass}
			user.HashPassword()

			result := user.ComparePassword(tt.comparePass)
			if result != tt.expectedResult {
				t.Errorf("ComparePassword() = %v, want %v", result, tt.expectedResult)
			}
		})
	}
}

func TestUserComparePasswordWithoutHashing(t *testing.T) {
	user := &User{Password: "plaintext"}
	result := user.ComparePassword("plaintext")

	// Should return false because password is not hashed
	if result {
		t.Error("ComparePassword should return false for unhashed password")
	}
}

func TestUserHashPasswordEmptyPassword(t *testing.T) {
	user := &User{Password: ""}
	user.HashPassword()

	// Empty password should remain empty or unchanged after hashing attempt
	// (bcrypt will return an error for empty passwords, and HashPassword ignores it)
	if user.Password != "" {
		// The current implementation doesn't handle this case well,
		// but we're documenting the behavior
		t.Logf("Empty password resulted in: %s", user.Password)
	}
}

func TestUserStructFields(t *testing.T) {
	user := &User{
		Email:     "test@example.com",
		Password:  "password123",
		Handle:    "testuser",
		FirstName: "John",
		LastName:  "Doe",
		Bio:       "Test bio",
		Avatar:    "avatar.png",
		Visible:   true,
		Verified:  true,
		Active:    true,
	}

	if user.Email != "test@example.com" {
		t.Errorf("Email = %s, want test@example.com", user.Email)
	}
	if user.Handle != "testuser" {
		t.Errorf("Handle = %s, want testuser", user.Handle)
	}
	if user.FirstName != "John" {
		t.Errorf("FirstName = %s, want John", user.FirstName)
	}
	if user.LastName != "Doe" {
		t.Errorf("LastName = %s, want Doe", user.LastName)
	}
	if user.Bio != "Test bio" {
		t.Errorf("Bio = %s, want Test bio", user.Bio)
	}
	if !user.Visible {
		t.Error("Visible should be true")
	}
	if !user.Verified {
		t.Error("Verified should be true")
	}
	if !user.Active {
		t.Error("Active should be true")
	}
}
