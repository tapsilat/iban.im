package resolvers

import (
	"context"
	"testing"

	"gorm.io/gorm"
)

func TestChangePassword(t *testing.T) {
	tests := []struct {
		name          string
		args          changePasswordMutationArgs
		setupDB       func(*gorm.DB) *uint
		withContext   bool
		expectSuccess bool
		expectError   string
	}{
		{
			name: "Successful password change",
			args: changePasswordMutationArgs{
				Password: "newpassword123",
			},
			setupDB: func(db *gorm.DB) *uint {
				user := createTestUser(t, db, "test@example.com", "oldpass", "testuser", "Test", "User")
				return &user.UserID
			},
			withContext:   true,
			expectSuccess: true,
		},
		{
			name: "Not authenticated",
			args: changePasswordMutationArgs{
				Password: "newpassword123",
			},
			withContext:   false,
			expectSuccess: false,
			expectError:   "Not Authorized",
		},
		{
			name: "Non-existing user",
			args: changePasswordMutationArgs{
				Password: "newpassword123",
			},
			setupDB: func(db *gorm.DB) *uint {
				userID := uint(999)
				return &userID
			},
			withContext:   true,
			expectSuccess: false,
			expectError:   "Not existing user",
		},
		{
			name: "Empty new password",
			args: changePasswordMutationArgs{
				Password: "",
			},
			setupDB: func(db *gorm.DB) *uint {
				user := createTestUser(t, db, "test@example.com", "oldpass", "testuser", "Test", "User")
				return &user.UserID
			},
			withContext:   true,
			expectSuccess: true, // Note: Current implementation doesn't validate empty password
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver, db, cleanup := setupTestResolverWithDB(t)
			defer cleanup()

			var ctx context.Context
			if tt.withContext && tt.setupDB != nil {
				userID := tt.setupDB(db)
				ctx = contextWithUserID(int(*userID))
			} else if tt.withContext {
				ctx = contextWithUserID(1)
			} else {
				ctx = context.Background()
			}

			if tt.setupDB != nil && !tt.withContext {
				tt.setupDB(db)
			}

			resp, err := resolver.ChangePassword(ctx, tt.args)
			if err != nil {
				t.Fatalf("ChangePassword returned unexpected error: %v", err)
			}

			if resp == nil {
				t.Fatal("ChangePassword returned nil response")
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

			if tt.expectSuccess && resp.User != nil {
				// Verify the new password was hashed
				if resp.User.u.Password == tt.args.Password {
					t.Error("Password should be hashed, not stored as plain text")
				}

				// Verify the hashed password can be validated
				if !resp.User.u.ComparePassword(tt.args.Password) {
					t.Error("New hashed password should match the provided password")
				}
			}
		})
	}
}

func TestChangePasswordHashingAndValidation(t *testing.T) {
	resolver, db, cleanup := setupTestResolverWithDB(t)
	defer cleanup()

	oldPassword := "oldpassword123"
	newPassword := "newpassword456"

	user := createTestUser(t, db, "test@example.com", oldPassword, "testuser", "Test", "User")
	ctx := contextWithUserID(int(user.UserID))

	args := changePasswordMutationArgs{
		Password: newPassword,
	}

	resp, err := resolver.ChangePassword(ctx, args)
	if err != nil {
		t.Fatalf("ChangePassword failed: %v", err)
	}

	if !resp.Ok() {
		t.Fatal("ChangePassword should succeed")
	}

	// Verify old password no longer works
	if resp.User.u.ComparePassword(oldPassword) {
		t.Error("Old password should no longer be valid")
	}

	// Verify new password works
	if !resp.User.u.ComparePassword(newPassword) {
		t.Error("New password should be valid")
	}
}
