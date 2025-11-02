package resolvers

import (
	"testing"

	"gorm.io/gorm"
)

func TestSignUp(t *testing.T) {
	tests := []struct {
		name          string
		args          signUpMutationArgs
		setupDB       func(*gorm.DB)
		expectSuccess bool
		expectError   string
	}{
		{
			name: "Successful signup",
			args: signUpMutationArgs{
				Email:     "newuser@example.com",
				Password:  "password123",
				Handle:    "newuser",
				FirstName: "John",
				LastName:  "Doe",
			},
			expectSuccess: true,
		},
		{
			name: "Duplicate email",
			args: signUpMutationArgs{
				Email:     "existing@example.com",
				Password:  "password123",
				Handle:    "newhandle",
				FirstName: "John",
				LastName:  "Doe",
			},
			setupDB: func(db *gorm.DB) {
				createTestUser(t, db, "existing@example.com", "pass123", "existinguser", "Jane", "Smith")
			},
			expectSuccess: false,
			expectError:   "Already signed up",
		},
		{
			name: "Duplicate handle",
			args: signUpMutationArgs{
				Email:     "newemail@example.com",
				Password:  "password123",
				Handle:    "existinghandle",
				FirstName: "John",
				LastName:  "Doe",
			},
			setupDB: func(db *gorm.DB) {
				createTestUser(t, db, "other@example.com", "pass123", "existinghandle", "Jane", "Smith")
			},
			expectSuccess: false,
			expectError:   "Already signed up",
		},
		{
			name: "Empty email",
			args: signUpMutationArgs{
				Email:     "",
				Password:  "password123",
				Handle:    "testhandle",
				FirstName: "John",
				LastName:  "Doe",
			},
			expectSuccess: true, // Note: The current implementation doesn't validate this
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver, db, cleanup := setupTestResolverWithDB(t)
			defer cleanup()

			if tt.setupDB != nil {
				tt.setupDB(db)
			}

			resp, err := resolver.SignUp(tt.args)
			if err != nil {
				t.Fatalf("SignUp returned unexpected error: %v", err)
			}

			if resp == nil {
				t.Fatal("SignUp returned nil response")
			}

			if resp.Ok() != tt.expectSuccess {
				t.Errorf("Ok() = %v, want %v", resp.Ok(), tt.expectSuccess)
			}

			if tt.expectError != "" {
				if resp.Error() == nil {
					t.Errorf("Expected error '%s', got nil", tt.expectError)
				} else if *resp.Error() != tt.expectError {
					t.Errorf("Error() = %s, want %s", *resp.Error(), tt.expectError)
				}
			}

			if tt.expectSuccess {
				if resp.User == nil {
					t.Error("Expected user in response, got nil")
				} else {
					if resp.User.Email() != tt.args.Email {
						t.Errorf("User email = %s, want %s", resp.User.Email(), tt.args.Email)
					}
					if resp.User.Handle() != tt.args.Handle {
						t.Errorf("User handle = %s, want %s", resp.User.Handle(), tt.args.Handle)
					}
					if resp.User.FirstName() != tt.args.FirstName {
						t.Errorf("User firstName = %s, want %s", resp.User.FirstName(), tt.args.FirstName)
					}
					if resp.User.LastName() != tt.args.LastName {
						t.Errorf("User lastName = %s, want %s", resp.User.LastName(), tt.args.LastName)
					}
				}
			}
		})
	}
}

func TestSignUpPasswordHashing(t *testing.T) {
	resolver, _, cleanup := setupTestResolverWithDB(t)
	defer cleanup()

	args := signUpMutationArgs{
		Email:     "test@example.com",
		Password:  "plainpassword",
		Handle:    "testuser",
		FirstName: "Test",
		LastName:  "User",
	}

	resp, err := resolver.SignUp(args)
	if err != nil {
		t.Fatalf("SignUp failed: %v", err)
	}

	if !resp.Ok() {
		t.Fatal("SignUp should succeed")
	}

	// Verify password was hashed (should not be plain text)
	if resp.User.u.Password == "plainpassword" {
		t.Error("Password should be hashed, not stored as plain text")
	}

	// Verify hashed password can be validated
	if !resp.User.u.ComparePassword("plainpassword") {
		t.Error("Hashed password should match original password")
	}
}
