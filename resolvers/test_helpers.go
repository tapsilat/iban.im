package resolvers

import (
	"context"
	"testing"

	"github.com/tapsilat/iban.im/config"
	"github.com/tapsilat/iban.im/handler"
	"github.com/tapsilat/iban.im/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto-migrate all models
	if err := db.AutoMigrate(&model.User{}, &model.Iban{}, &model.Group{}); err != nil {
		t.Fatalf("Failed to auto-migrate: %v", err)
	}

	return db
}

// createTestUser creates a test user in the database
func createTestUser(t *testing.T, db *gorm.DB, email, password, handle, firstName, lastName string) *model.User {
	user := &model.User{
		Email:     email,
		Password:  password,
		Handle:    handle,
		FirstName: firstName,
		LastName:  lastName,
		Active:    true,
		Verified:  false,
	}
	user.HashPassword()

	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	return user
}

// createTestIban creates a test IBAN in the database
func createTestIban(t *testing.T, db *gorm.DB, ownerID uint, text, handle, password string, isPrivate bool) *model.Iban {
	iban := &model.Iban{
		Text:      text,
		Handle:    handle,
		Password:  password,
		OwnerID:   ownerID,
		OwnerType: "User",
		Active:    true,
		IsPrivate: isPrivate,
	}
	if isPrivate && password != "" {
		iban.HashPassword()
	}

	if err := db.Create(iban).Error; err != nil {
		t.Fatalf("Failed to create test IBAN: %v", err)
	}

	return iban
}

// contextWithUserID creates a context with a user ID for authenticated requests
func contextWithUserID(userID int) context.Context {
	ctx := context.Background()
	return context.WithValue(ctx, handler.ContextKey("UserID"), userID)
}

// setupTestResolverWithDB sets up a test resolver with a test database
func setupTestResolverWithDB(t *testing.T) (*Resolvers, *gorm.DB, func()) {
	db := setupTestDB(t)
	
	// Store original DB and replace it for testing
	originalDB := config.DB
	config.DB = db
	
	resolver := &Resolvers{}
	
	cleanup := func() {
		config.DB = originalDB
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}
	
	return resolver, db, cleanup
}
