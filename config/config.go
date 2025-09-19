package config

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

type DBConfig struct {
	Name     string `env:"DBName" default:"ibanim"`
	Adapter  string `env:"DBAdapter" default:"postgres"`
	Host     string `env:"DBHost" default:"localhost"`
	Port     string `env:"DBPort" default:"5432"`
	User     string `env:"DBUser" default:"ibanim"`
	Password string `env:"DBPassword" default:"ibanim"`
}

type SMTPConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

type AppConfig struct {
	Port  uint   `default:"7000" env:"PORT"`
	Env   string `default:"localhost" env:"ENV"`
	Debug bool   `default:"false" env:"DEBUG"`

	// Session Timeout - minutes
	Timeout uint `default:"60"`

	// Auth Max Refresh - minutes
	MaxRefresh uint `default:"60"`

	// Auth Key
	Key string `default:"12345678" env:"AUTH_KEY"`

	// Realm name to display to the user.
	Realm string `default:"ibanim zone"`
}

var Config = struct {
	Db   DBConfig
	Smtp SMTPConfig
	App  AppConfig
}{}

func init() {
	// Set sane defaults (previously provided by struct tags via configor)
	Config.Db = DBConfig{
		Name:     "ibanim",
		Adapter:  "postgres",
		Host:     "localhost",
		Port:     "5432",
		User:     "ibanim",
		Password: "ibanim",
	}
	Config.App = AppConfig{
		Port:       7000,
		Env:        "localhost",
		Debug:      false,
		Timeout:    60,
		MaxRefresh: 60,
		Key:        "12345678",
		Realm:      "ibanim zone",
	}

	// Load YAML files if present; missing files are ignored
	for _, path := range []string{
		"config/database.yml",
		"config/smtp.yml",
		"config/application.yml",
	} {
		if err := loadYAML(path); err != nil {
			panic(err)
		}
	}
}

// loadYAML unmarshals the YAML file into the global Config, updating only fields present in the file.
func loadYAML(path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read %s: %w", path, err)
	}
	if err := yaml.Unmarshal(b, &Config); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	return nil
}
