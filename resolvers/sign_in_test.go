package resolvers

import (
	"testing"

	"github.com/tapsilat/iban.im/model"
)

func TestSignIn(t *testing.T) {
	_, db, cleanup := setupTestResolverWithDB(t)
	defer cleanup()

	// Create a test user
	email := "test@example.com"
	password := "testpassword"
	user := createTestUser(t, db, email, password, "testuser", "Test", "User")

	// Test that we can query the user from the database
	if user.UserID == 0 {
		t.Error("User ID should not be 0")
	}

	// Test querying for a non-existing user
	var foundUser model.User
	result := db.Where("email = ?", "notexisting@test.com").First(&foundUser)
	if result.Error == nil {
		t.Error("Should not find a non-existing user")
	}

	// Test querying for the existing user
	var existingUser model.User
	result = db.Where("email = ?", email).First(&existingUser)
	if result.Error != nil {
		t.Errorf("Should find existing user: %v", result.Error)
	}
	if existingUser.Email != email {
		t.Errorf("Found user email = %s, want %s", existingUser.Email, email)
	}

	// Test password validation
	if !existingUser.ComparePassword(password) {
		t.Error("Password should match")
	}
	if existingUser.ComparePassword("wrongpassword") {
		t.Error("Wrong password should not match")
	}
}
