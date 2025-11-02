package resolvers

import (
	"context"
	"strings"
	"testing"

	"gorm.io/gorm"
)

func TestChangeProfile(t *testing.T) {
	tests := []struct {
		name          string
		args          changeProfileMutationArgs
		setupDB       func(*gorm.DB) *uint
		withContext   bool
		expectSuccess bool
		expectError   string
	}{
		{
			name: "Update bio",
			args: changeProfileMutationArgs{
				Bio: strPtr("New bio description"),
			},
			setupDB: func(db *gorm.DB) *uint {
				user := createTestUser(t, db, "test@example.com", "pass", "testuser", "Test", "User")
				return &user.UserID
			},
			withContext:   true,
			expectSuccess: true,
		},
		{
			name: "Update handle",
			args: changeProfileMutationArgs{
				Handle: strPtr("newhandle"),
			},
			setupDB: func(db *gorm.DB) *uint {
				user := createTestUser(t, db, "test@example.com", "pass", "oldhandle", "Test", "User")
				return &user.UserID
			},
			withContext:   true,
			expectSuccess: true,
		},
		{
			name: "Update both bio and handle",
			args: changeProfileMutationArgs{
				Bio:    strPtr("Updated bio"),
				Handle: strPtr("updatedhandle"),
			},
			setupDB: func(db *gorm.DB) *uint {
				user := createTestUser(t, db, "test@example.com", "pass", "testuser", "Test", "User")
				return &user.UserID
			},
			withContext:   true,
			expectSuccess: true,
		},
		{
			name: "Not authenticated",
			args: changeProfileMutationArgs{
				Bio: strPtr("New bio"),
			},
			withContext:   false,
			expectSuccess: false,
			expectError:   "Not Authorized",
		},
		{
			name: "Non-existing user",
			args: changeProfileMutationArgs{
				Bio: strPtr("New bio"),
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
			name: "Update with empty values",
			args: changeProfileMutationArgs{
				Bio:    strPtr(""),
				Handle: strPtr("testuser"),
			},
			setupDB: func(db *gorm.DB) *uint {
				user := createTestUser(t, db, "test@example.com", "pass", "testuser", "Test", "User")
				return &user.UserID
			},
			withContext:   true,
			expectSuccess: true,
		},
		{
			name: "Handle case insensitivity",
			args: changeProfileMutationArgs{
				Handle: strPtr("NewHandle"),
			},
			setupDB: func(db *gorm.DB) *uint {
				user := createTestUser(t, db, "test@example.com", "pass", "oldhandle", "Test", "User")
				return &user.UserID
			},
			withContext:   true,
			expectSuccess: true,
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

			resp, err := resolver.ChangeProfile(ctx, tt.args)
			if err != nil {
				t.Fatalf("ChangeProfile returned unexpected error: %v", err)
			}

			if resp == nil {
				t.Fatal("ChangeProfile returned nil response")
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
				if tt.args.Bio != nil && *resp.User.Bio() != *tt.args.Bio {
					t.Errorf("User bio = %s, want %s", *resp.User.Bio(), *tt.args.Bio)
				}
				if tt.args.Handle != nil {
					// Handle should be lowercased by the resolver
					expectedHandle := strings.ToLower(*tt.args.Handle)
					if resp.User.Handle() != expectedHandle {
						t.Errorf("User handle = %s, want %s", resp.User.Handle(), expectedHandle)
					}
				}
			}
		})
	}
}

func TestDeleteProfile(t *testing.T) {
	tests := []struct {
		name          string
		args          deleteProfileMutationArgs
		setupDB       func(*gorm.DB) (*uint, string)
		withContext   bool
		expectSuccess bool
		expectError   string
	}{
		{
			name: "Successful deletion with correct password",
			args: deleteProfileMutationArgs{
				ConfirmPassword: "correctpassword",
			},
			setupDB: func(db *gorm.DB) (*uint, string) {
				password := "correctpassword"
				user := createTestUser(t, db, "test@example.com", password, "testuser", "Test", "User")
				return &user.UserID, password
			},
			withContext:   true,
			expectSuccess: true,
		},
		{
			name: "Deletion with IBANs",
			args: deleteProfileMutationArgs{
				ConfirmPassword: "password123",
			},
			setupDB: func(db *gorm.DB) (*uint, string) {
				password := "password123"
				user := createTestUser(t, db, "test@example.com", password, "testuser", "Test", "User")
				createTestIban(t, db, user.UserID, "TR320010009999901234567890", "iban1", "", false)
				createTestIban(t, db, user.UserID, "TR420010009999901234567891", "iban2", "", false)
				return &user.UserID, password
			},
			withContext:   true,
			expectSuccess: true,
		},
		{
			name: "Incorrect password",
			args: deleteProfileMutationArgs{
				ConfirmPassword: "wrongpassword",
			},
			setupDB: func(db *gorm.DB) (*uint, string) {
				password := "correctpassword"
				user := createTestUser(t, db, "test@example.com", password, "testuser", "Test", "User")
				return &user.UserID, password
			},
			withContext:   true,
			expectSuccess: false,
			expectError:   "Invalid password confirmation",
		},
		{
			name: "Not authenticated",
			args: deleteProfileMutationArgs{
				ConfirmPassword: "password",
			},
			withContext:   false,
			expectSuccess: false,
			expectError:   "Not Authorized",
		},
		{
			name: "Non-existing user",
			args: deleteProfileMutationArgs{
				ConfirmPassword: "password",
			},
			setupDB: func(db *gorm.DB) (*uint, string) {
				userID := uint(999)
				return &userID, "password"
			},
			withContext:   true,
			expectSuccess: false,
			expectError:   "User not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver, db, cleanup := setupTestResolverWithDB(t)
			defer cleanup()

			var ctx context.Context
			if tt.withContext && tt.setupDB != nil {
				userID, _ := tt.setupDB(db)
				ctx = contextWithUserID(int(*userID))
			} else if tt.withContext {
				ctx = contextWithUserID(1)
			} else {
				ctx = context.Background()
			}

			if tt.setupDB != nil && !tt.withContext {
				tt.setupDB(db)
			}

			resp, err := resolver.DeleteProfile(ctx, tt.args)
			if err != nil {
				t.Fatalf("DeleteProfile returned unexpected error: %v", err)
			}

			if resp == nil {
				t.Fatal("DeleteProfile returned nil response")
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

			if tt.expectSuccess && resp.Message() != nil {
				// Verify success message exists
				if *resp.Message() == "" {
					t.Error("Expected success message to be non-empty")
				}
			}
		})
	}
}
