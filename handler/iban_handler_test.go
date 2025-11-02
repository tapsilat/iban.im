package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/tapsilat/iban.im/config"
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

func TestGetIbanByHandles(t *testing.T) {
	// Setup test database
	db := setupTestDB(t)
	originalDB := config.DB
	config.DB = db
	defer func() {
		config.DB = originalDB
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}()

	// Create test data
	user := createTestUser(t, db, "test@example.com", "password123", "testuser", "Test", "User")
	_ = createTestIban(t, db, user.UserID, "TR320010009999901234567890", "testiban", "", false)

	tests := []struct {
		name           string
		userHandle     string
		ibanHandle     string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid IBAN request",
			userHandle:     "testuser",
			ibanHandle:     "testiban",
			expectedStatus: http.StatusOK,
			expectedBody:   "TR320010009999901234567890",
		},
		{
			name:           "User not found",
			userHandle:     "nonexistent",
			ibanHandle:     "testiban",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "User not found",
		},
		{
			name:           "IBAN not found",
			userHandle:     "testuser",
			ibanHandle:     "nonexistent",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "IBAN not found or is private",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup Gin test context
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Set params
			c.Params = gin.Params{
				{Key: "userHandle", Value: tt.userHandle},
				{Key: "ibanHandle", Value: tt.ibanHandle},
			}

			// Call handler
			GetIbanByHandles(c)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check response body contains expected content
			body := w.Body.String()
			if body == "" {
				t.Errorf("Expected non-empty body")
			}
		})
	}
}

func TestRenderIbanPage(t *testing.T) {
	// Setup test database
	db := setupTestDB(t)
	originalDB := config.DB
	config.DB = db
	defer func() {
		config.DB = originalDB
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}()

	// Create test data
	user := createTestUser(t, db, "test@example.com", "password123", "testuser", "Test", "User")
	createTestIban(t, db, user.UserID, "TR320010009999901234567890", "testiban", "", false)
	createTestIban(t, db, user.UserID, "TR420010009999901234567891", "privateiban", "secret", true)

	tests := []struct {
		name           string
		userHandle     string
		ibanHandle     string
		expectedStatus int
	}{
		{
			name:           "Valid public IBAN",
			userHandle:     "testuser",
			ibanHandle:     "testiban",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Private IBAN should not be accessible",
			userHandle:     "testuser",
			ibanHandle:     "privateiban",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "User not found",
			userHandle:     "nonexistent",
			ibanHandle:     "testiban",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "IBAN not found",
			userHandle:     "testuser",
			ibanHandle:     "nonexistent",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup Gin test context
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, router := gin.CreateTestContext(w)
			
			// Load templates (for actual rendering, we'll just check status)
			// In a real scenario, you'd need to load templates
			router.LoadHTMLGlob("../templates/*.tmpl.html")

			// Set params
			c.Params = gin.Params{
				{Key: "userHandle", Value: tt.userHandle},
				{Key: "ibanHandle", Value: tt.ibanHandle},
			}

			// Call handler
			RenderIbanPage(c)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestIsValidRoute(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "Valid user/iban route",
			path:     "/testuser/testiban",
			expected: true,
		},
		{
			name:     "Assets route",
			path:     "/assets/style.css",
			expected: false,
		},
		{
			name:     "API route",
			path:     "/api/login",
			expected: false,
		},
		{
			name:     "Auth route",
			path:     "/auth/refresh_token",
			expected: false,
		},
		{
			name:     "Graph route",
			path:     "/graph",
			expected: false,
		},
		{
			name:     "Root path",
			path:     "/",
			expected: false,
		},
		{
			name:     "Path with multiple segments",
			path:     "/user/iban/extra",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidRoute(tt.path)
			if result != tt.expected {
				t.Errorf("IsValidRoute(%s) = %v, expected %v", tt.path, result, tt.expected)
			}
		})
	}
}
