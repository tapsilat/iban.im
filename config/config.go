package config

import (
	"os"
	"strconv"
)

type DBConfig struct {
	Name     string
	Adapter  string
	Host     string
	Port     string
	User     string
	Password string
}

type SMTPConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

type AppConfig struct {
	Port       uint
	Env        string
	Debug      bool
	Timeout    uint
	MaxRefresh uint
	Key        string
	Realm      string
}

var Config = struct {
	Db   DBConfig
	Smtp SMTPConfig
	App  AppConfig
}{}

// init loads configuration strictly from environment variables.
// Any missing values remain zero-values.
func init() {
	// DB
	if v := os.Getenv("DB_NAME"); v != "" {
		Config.Db.Name = v
	}
	if v := os.Getenv("DB_ADAPTER"); v != "" {
		Config.Db.Adapter = v
	}
	if v := os.Getenv("DB_HOST"); v != "" {
		Config.Db.Host = v
	}
	if v := os.Getenv("DB_PORT"); v != "" {
		Config.Db.Port = v
	}
	if v := os.Getenv("DB_USER"); v != "" {
		Config.Db.User = v
	}
	if v := os.Getenv("DB_PASSWORD"); v != "" {
		Config.Db.Password = v
	}

	// App
	if v := os.Getenv("PORT"); v != "" {
		if n, err := strconv.ParseUint(v, 10, 64); err == nil {
			Config.App.Port = uint(n)
		}
	}
	if v := os.Getenv("ENV"); v != "" {
		Config.App.Env = v
	}
	if v := os.Getenv("DEBUG"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			Config.App.Debug = b
		}
	}
	if v := os.Getenv("TIMEOUT"); v != "" {
		if n, err := strconv.ParseUint(v, 10, 64); err == nil {
			Config.App.Timeout = uint(n)
		}
	}
	if v := os.Getenv("MAX_REFRESH"); v != "" {
		if n, err := strconv.ParseUint(v, 10, 64); err == nil {
			Config.App.MaxRefresh = uint(n)
		}
	}
	if v := os.Getenv("AUTH_KEY"); v != "" {
		Config.App.Key = v
	}
	if v := os.Getenv("REALM"); v != "" {
		Config.App.Realm = v
	}
}
