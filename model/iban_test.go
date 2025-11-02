package model

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestIbanHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{
			name:     "Valid password",
			password: "secret123",
		},
		{
			name:     "Empty password",
			password: "",
		},
		{
			name:     "Long password",
			password: "this_is_a_very_long_iban_password_with_many_characters_1234567890",
		},
		{
			name:     "Special characters",
			password: "s3cr3t!@#$",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iban := &Iban{Password: tt.password}
			originalPassword := tt.password
			iban.HashPassword()

			// Check that password was hashed
			if iban.Password == originalPassword && originalPassword != "" {
				t.Errorf("Password was not hashed: got %s, want different from %s", iban.Password, originalPassword)
			}

			// Check that hashed password is not empty
			if tt.password != "" && iban.Password == "" {
				t.Error("Hashed password should not be empty for non-empty input")
			}

			// Check that hashed password starts with bcrypt prefix
			if tt.password != "" && len(iban.Password) < 10 {
				t.Error("Hashed password is too short to be a valid bcrypt hash")
			}
		})
	}
}

func TestIbanComparePassword(t *testing.T) {
	tests := []struct {
		name           string
		originalPass   string
		comparePass    string
		expectedResult bool
	}{
		{
			name:           "Matching passwords",
			originalPass:   "secret123",
			comparePass:    "secret123",
			expectedResult: true,
		},
		{
			name:           "Non-matching passwords",
			originalPass:   "secret123",
			comparePass:    "wrongsecret",
			expectedResult: false,
		},
		{
			name:           "Empty comparison password",
			originalPass:   "secret123",
			comparePass:    "",
			expectedResult: false,
		},
		{
			name:           "Case sensitive comparison",
			originalPass:   "Secret123",
			comparePass:    "secret123",
			expectedResult: false,
		},
		{
			name:           "Special characters",
			originalPass:   "s3cr3t!@#",
			comparePass:    "s3cr3t!@#",
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iban := &Iban{Password: tt.originalPass}
			iban.HashPassword()

			result := iban.ComparePassword(tt.comparePass)
			if result != tt.expectedResult {
				t.Errorf("ComparePassword() = %v, want %v", result, tt.expectedResult)
			}
		})
	}
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto-migrate the Iban model
	if err := db.AutoMigrate(&Iban{}); err != nil {
		t.Fatalf("Failed to auto-migrate: %v", err)
	}

	return db
}

func TestIbanCheckHandle(t *testing.T) {
	db := setupTestDB(t)

	tests := []struct {
		name           string
		existingIbans  []Iban
		checkIban      Iban
		expectedExists bool
	}{
		{
			name: "No duplicate handle",
			existingIbans: []Iban{
				{Handle: "existing1", Text: "TR320010009999901234567890", OwnerID: 1},
			},
			checkIban:      Iban{Handle: "newhandle", Text: "TR420010009999901234567891", OwnerID: 1},
			expectedExists: false,
		},
		{
			name: "Duplicate handle same owner",
			existingIbans: []Iban{
				{Handle: "duplicate", Text: "TR320010009999901234567890", OwnerID: 1},
			},
			checkIban:      Iban{Handle: "duplicate", Text: "TR420010009999901234567891", OwnerID: 1},
			expectedExists: true,
		},
		{
			name: "Same handle different owner",
			existingIbans: []Iban{
				{Handle: "samename", Text: "TR320010009999901234567890", OwnerID: 1},
			},
			checkIban:      Iban{Handle: "samename", Text: "TR420010009999901234567891", OwnerID: 2},
			expectedExists: false,
		},
		{
			name: "Update same IBAN with same handle",
			existingIbans: []Iban{
				{IbanID: 1, Handle: "handle1", Text: "TR320010009999901234567890", OwnerID: 1},
			},
			checkIban:      Iban{IbanID: 1, Handle: "handle1", Text: "TR320010009999901234567890", OwnerID: 1},
			expectedExists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear the database
			db.Exec("DELETE FROM ibans")

			// Insert existing IBANs
			for _, iban := range tt.existingIbans {
				if err := db.Create(&iban).Error; err != nil {
					t.Fatalf("Failed to create test IBAN: %v", err)
				}
			}

			// Check if handle exists
			exists := tt.checkIban.CheckHandle(db)
			if exists != tt.expectedExists {
				t.Errorf("CheckHandle() = %v, want %v", exists, tt.expectedExists)
			}
		})
	}
}

func TestIbanBeforeSave(t *testing.T) {
	db := setupTestDB(t)

	tests := []struct {
		name          string
		existingIbans []Iban
		newIban       Iban
		expectError   bool
	}{
		{
			name: "No conflict - should succeed",
			existingIbans: []Iban{
				{Handle: "existing1", Text: "TR320010009999901234567890", OwnerID: 1},
			},
			newIban:     Iban{Handle: "newhandle", Text: "TR420010009999901234567891", OwnerID: 1},
			expectError: false,
		},
		{
			name: "Duplicate handle - should fail",
			existingIbans: []Iban{
				{Handle: "duplicate", Text: "TR320010009999901234567890", OwnerID: 1},
			},
			newIban:     Iban{Handle: "duplicate", Text: "TR420010009999901234567891", OwnerID: 1},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear the database
			db.Exec("DELETE FROM ibans")

			// Insert existing IBANs
			for _, iban := range tt.existingIbans {
				if err := db.Create(&iban).Error; err != nil {
					t.Fatalf("Failed to create test IBAN: %v", err)
				}
			}

			// Try to create new IBAN
			err := db.Create(&tt.newIban).Error
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestIbanValidate(t *testing.T) {
	db := setupTestDB(t)

	tests := []struct {
		name        string
		iban        Iban
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid IBAN with all fields",
			iban:        Iban{Text: "TR320010009999901234567890", Handle: "myiban", IsPrivate: false},
			expectError: false,
		},
		{
			name:        "Valid private IBAN with password",
			iban:        Iban{Text: "TR320010009999901234567890", Handle: "myiban", IsPrivate: true, Password: "secret"},
			expectError: false,
		},
		{
			name:        "Empty IBAN text",
			iban:        Iban{Text: "", Handle: "myiban", IsPrivate: false},
			expectError: true,
			errorMsg:    "you have to provide IBAN",
		},
		{
			name:        "Whitespace only IBAN text",
			iban:        Iban{Text: "   ", Handle: "myiban", IsPrivate: false},
			expectError: true,
			errorMsg:    "you have to provide IBAN",
		},
		{
			name:        "Empty handle",
			iban:        Iban{Text: "TR320010009999901234567890", Handle: "", IsPrivate: false},
			expectError: true,
			errorMsg:    "you have to provide handle",
		},
		{
			name:        "Whitespace only handle",
			iban:        Iban{Text: "TR320010009999901234567890", Handle: "   ", IsPrivate: false},
			expectError: true,
			errorMsg:    "you have to provide handle",
		},
		{
			name:        "Private IBAN without password",
			iban:        Iban{Text: "TR320010009999901234567890", Handle: "myiban", IsPrivate: true, Password: ""},
			expectError: true,
			errorMsg:    "you have to provide password",
		},
		{
			name:        "Private IBAN with whitespace password",
			iban:        Iban{Text: "TR320010009999901234567890", Handle: "myiban", IsPrivate: true, Password: "   "},
			expectError: true,
			errorMsg:    "you have to provide password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new DB instance for each test
			testDB := db.Session(&gorm.Session{})
			
			tt.iban.Validate(testDB)
			
			err := testDB.Error
			if tt.expectError && err == nil {
				t.Error("Expected validation error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected validation error: %v", err)
			}
			if tt.expectError && err != nil && err.Error() != tt.errorMsg {
				t.Errorf("Expected error message '%s', got '%s'", tt.errorMsg, err.Error())
			}
		})
	}
}

func TestIbanStructFields(t *testing.T) {
	iban := &Iban{
		Text:        "TR320010009999901234567890",
		Description: "Test IBAN",
		Password:    "secret",
		Handle:      "myiban",
		Active:      true,
		IsPrivate:   true,
		OwnerID:     1,
		OwnerType:   "User",
	}

	if iban.Text != "TR320010009999901234567890" {
		t.Errorf("Text = %s, want TR320010009999901234567890", iban.Text)
	}
	if iban.Description != "Test IBAN" {
		t.Errorf("Description = %s, want Test IBAN", iban.Description)
	}
	if iban.Handle != "myiban" {
		t.Errorf("Handle = %s, want myiban", iban.Handle)
	}
	if !iban.Active {
		t.Error("Active should be true")
	}
	if !iban.IsPrivate {
		t.Error("IsPrivate should be true")
	}
	if iban.OwnerID != 1 {
		t.Errorf("OwnerID = %d, want 1", iban.OwnerID)
	}
}
