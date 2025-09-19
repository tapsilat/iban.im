package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/tapsilat/iban.im/model"
)

var DB *gorm.DB

func InitDB() {
	var err error

	adapter := os.Getenv("DB_ADAPTER")
	if adapter == "" {
		adapter = "postgres"
	}

	var connStr string
	if adapter == "mysql" {
		connStr = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True&loc=Local",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		)
	} else if adapter == "postgres" {
		connStr = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		)
	} else {
		panic(errors.New("unsupported database adapter"))
	}

	DB, err = gorm.Open(adapter, connStr)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	// Veritabanı bağlantısını kontrol etme
	if err := DB.DB().Ping(); err != nil {
		log.Fatalf("Error connecting to DB: %v", err)
	} else {
		log.Println("Successfully connected to DB")
	}

	DB.LogMode(true)
	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(30)
	DB.DB().SetConnMaxLifetime(time.Second * 60)

	DB.AutoMigrate(&model.User{}, &model.Iban{}, &model.Group{})
}
