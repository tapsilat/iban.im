package main

import (
	"fmt"
	"log"

	"github.com/tapsilat/iban.im/config"
	"github.com/tapsilat/iban.im/model"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	config.InitDB(cfg)

	// Create test user
	user := &model.User{
		Email:     "test@example.com",
		Password:  "password123",
		Handle:    "testuser",
		FirstName: "Test",
		LastName:  "User",
		Active:    true,
		Verified:  true,
	}
	user.HashPassword()

	if err := config.DB.Create(user).Error; err != nil {
		log.Printf("User might already exist: %v", err)
		// Try to find existing user
		if err := config.DB.Where("handle = ?", "testuser").First(user).Error; err != nil {
			log.Fatalf("Failed to find or create user: %v", err)
		}
	}

	fmt.Printf("User created/found: %s (ID: %d)\n", user.Handle, user.UserID)

	// Create test IBANs
	ibans := []*model.Iban{
		{
			Text:        "TR320010009999901234567890",
			Description: "My primary bank account",
			Handle:      "primary",
			OwnerID:     user.UserID,
			OwnerType:   "User",
			Active:      true,
			IsPrivate:   false,
		},
		{
			Text:        "TR420010009999901234567891",
			Description: "Savings account",
			Handle:      "savings",
			OwnerID:     user.UserID,
			OwnerType:   "User",
			Active:      true,
			IsPrivate:   false,
		},
		{
			Text:        "TR520010009999901234567892",
			Description: "Private account",
			Handle:      "private",
			Password:    "secret123",
			OwnerID:     user.UserID,
			OwnerType:   "User",
			Active:      true,
			IsPrivate:   true,
		},
	}

	for _, iban := range ibans {
		if iban.IsPrivate {
			iban.HashPassword()
		}
		if err := config.DB.Create(iban).Error; err != nil {
			log.Printf("IBAN %s might already exist: %v", iban.Handle, err)
		} else {
			fmt.Printf("IBAN created: %s/%s\n", user.Handle, iban.Handle)
		}
	}

	fmt.Println("\nTest data created successfully!")
	fmt.Printf("You can access IBANs at:\n")
	fmt.Printf("  http://localhost:8080/%s/primary\n", user.Handle)
	fmt.Printf("  http://localhost:8080/%s/savings\n", user.Handle)
	fmt.Printf("  http://localhost:8080/%s/private (should return 404 - it's private)\n", user.Handle)
}
