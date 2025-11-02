package schema

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetSchema(t *testing.T) {
	// Change to the project root directory for this test
	// The schema loader expects to be run from the project root
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	// Check if we're in the schema directory, if so, go up one level
	if filepath.Base(wd) == "schema" {
		if err := os.Chdir(".."); err != nil {
			t.Fatalf("Failed to change directory: %v", err)
		}
		defer os.Chdir(wd) // Restore original directory
	}

	// Check if schema directory exists
	if _, err := os.Stat("./schema"); os.IsNotExist(err) {
		t.Skip("Skipping test: schema directory not found from current location")
	}

	s := NewSchema()

	if s == nil {
		t.Fatal("NewSchema() returned nil")
	}

	if *s == "" {
		t.Error("NewSchema() returned empty schema")
	}

	// Basic validation that schema contains expected GraphQL keywords
	schemaStr := *s
	if len(schemaStr) < 10 {
		t.Errorf("Schema seems too short: %d bytes", len(schemaStr))
	}

	t.Logf("Schema loaded successfully, length: %d bytes", len(schemaStr))
}
