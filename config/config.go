package config

import (
	"log"
	"os"
	"sync"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type DBase struct {
	Name     string `env:"DB_NAME" default:"ibanim"`
	Adapter  string `env:"DB_ADAPTER" default:"postgres"`
	Host     string `env:"DB_HOST" default:"localhost"`
	Port     string `env:"DB_PORT" default:"5432"`
	User     string `env:"DB_USER" default:"ibanim"`
	Password string `env:"DB_PASSWORD" default:"ibanim"`
}

type SMTP struct {
	Host     string `env:"SMTP_HOST" envDefault:"localhost"`
	Port     string `env:"SMTP_PORT" envDefault:"2525"`
	User     string `env:"SMTP_USER"`
	Password string `env:"SMTP_PASSWORD"`
}

type App struct {
	Port       string `env:"APP_PORT" envDefault:"8080"`
	Env        string `env:"APP_ENV" envDefault:"development"`
	Debug      bool   `env:"APP_DEBUG" envDefault:"false"`
	Timeout    uint   `env:"APP_TIMEOUT" envDefault:"5"`
	MaxRefresh uint   `env:"APP_MAX_REFRESH" envDefault:"5"`
	Key        string `env:"APP_KEY"`
	Realm      string `env:"APP_REALM"`
}

type Config struct {
	DB   DBase
	Smtp SMTP
	App  App
}

var (
	globalConfig *Config
	configOnce   sync.Once
)

func GetGlobalConfig() *Config {
	configOnce.Do(func() {
		var err error
		globalConfig, err = GetConfig()
		if err != nil {
			log.Fatal(err)
		}
	})
	return globalConfig
}

func GetConfig() (*Config, error) {
	envFile := ".env"
	if env := os.Getenv("ENV"); env != "" {
		envFile = ".env." + env
	}

	if err := godotenv.Load(envFile); err != nil {
		// Not a critical error, continue with environment variables only
	}

	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
