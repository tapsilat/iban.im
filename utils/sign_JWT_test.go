package utils

import (
	"testing"
)

// TestSignJWT is skipped because it requires a running server
// This is an integration test that should be run separately
func TestSignJWT(t *testing.T) {
	t.Skip("Skipping test: requires running server on localhost:8080")
	
	// This test needs the full application stack running
	// including the /api/login endpoint
	// For actual testing, use integration tests or mock the HTTP client
	
	userMail := "test@example.com"
	userPass := "testpassword"
	token, err := SignJWT(&userMail, &userPass)
	if err != nil {
		t.Error(err)
	}
	t.Log(*token)
}
