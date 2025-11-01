package config

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/tapsilat/iban.im/model"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(cfg *Config) {
	var err error

	// Build DSN and open using GORM v2 drivers
	var dsn string
	var dialector gorm.Dialector
	if cfg.DB.Adapter == "mysql" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True&loc=Local",
			cfg.DB.User,
			cfg.DB.Password,
			cfg.DB.Host,
			cfg.DB.Port,
			cfg.DB.Name,
		)
		dialector = mysql.Open(dsn)
	} else if cfg.DB.Adapter == "postgres" {
		// Postgres driver supports URL dsn
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.DB.User,
			cfg.DB.Password,
			cfg.DB.Host,
			cfg.DB.Port,
			cfg.DB.Name,
		)
		dialector = postgres.Open(dsn)
	} else if cfg.DB.Adapter == "sqlite" {
		// SQLite driver - DB.Name is the file path
		dsn = cfg.DB.Name
		dialector = sqlite.Open(dsn)
	} else {
		panic(errors.New("unsupported database adapter"))
	}

	DB, err = gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	// Veritabanı bağlantısını kontrol etme
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Error getting underlying DB: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Error connecting to DB: %v", err)
	} else {
		log.Println("Successfully connected to DB")
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(30)
	sqlDB.SetConnMaxLifetime(time.Second * 60)

	DB.AutoMigrate(&model.User{}, &model.Iban{}, &model.Group{})
}
