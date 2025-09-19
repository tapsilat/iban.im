package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/goccy/go-yaml"
)

type DBConfig struct {
	Name     string `env:"DB_NAME" default:"ibanim"`
	Adapter  string `env:"DB_ADAPTER" default:"postgres"`
	Host     string `env:"DB_HOST" default:"localhost"`
	Port     string `env:"DB_PORT" default:"5432"`
	User     string `env:"DB_USER" default:"ibanim"`
	Password string `env:"DB_PASSWORD" default:"ibanim"`
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
	Timeout uint `default:"60" env:"TIMEOUT"`

	// Auth Max Refresh - minutes
	MaxRefresh uint `default:"60" env:"MAX_REFRESH"`

	// Auth Key
	Key string `default:"12345678" env:"AUTH_KEY"`

	// Realm name to display to the user.
	Realm string `default:"ibanim zone" env:"REALM"`
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

	// Apply environment variable overrides
	applyEnvOverrides()
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

// applyEnvOverrides reads environment variables and overrides config fields
// using the current naming convention (UPPER_SNAKE_CASE).
func applyEnvOverrides() {
	// DB
	if v, ok := os.LookupEnv("DB_NAME"); ok {
		Config.Db.Name = v
	}
	if v, ok := os.LookupEnv("DB_ADAPTER"); ok {
		Config.Db.Adapter = v
	}
	if v, ok := os.LookupEnv("DB_HOST"); ok {
		Config.Db.Host = v
	}
	if v, ok := os.LookupEnv("DB_PORT"); ok {
		Config.Db.Port = v
	}
	if v, ok := os.LookupEnv("DB_USER"); ok {
		Config.Db.User = v
	}
	if v, ok := os.LookupEnv("DB_PASSWORD"); ok {
		Config.Db.Password = v
	}

	// App
	if v, ok := os.LookupEnv("PORT"); ok {
		if n, err := strconv.ParseUint(v, 10, 64); err == nil {
			Config.App.Port = uint(n)
		}
	}
	if v, ok := os.LookupEnv("ENV"); ok {
		Config.App.Env = v
	}
	if v, ok := os.LookupEnv("DEBUG"); ok {
		if b, err := strconv.ParseBool(v); err == nil {
			Config.App.Debug = b
		}
	}
	if v, ok := os.LookupEnv("TIMEOUT"); ok {
		if n, err := strconv.ParseUint(v, 10, 64); err == nil {
			Config.App.Timeout = uint(n)
		}
	}
	if v, ok := os.LookupEnv("MAX_REFRESH"); ok {
		if n, err := strconv.ParseUint(v, 10, 64); err == nil {
			Config.App.MaxRefresh = uint(n)
		}
	}
	if v, ok := os.LookupEnv("AUTH_KEY"); ok {
		Config.App.Key = v
	}
	if v, ok := os.LookupEnv("REALM"); ok {
		Config.App.Realm = v
	}
}
