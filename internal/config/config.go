package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Database DBConfig
	CORS     CORSConfig
	Auth 	 AuthConfig
}

type AppConfig struct {
	Name string
	Env  string
	Port string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type CORSConfig struct {
	AllowedOrigins []string
}

type AuthConfig struct {
	JWT_SECRET string
}

func LoadConfig() (*Config, error) {
	slog.Info("Loading application config...")

	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	config := &Config{
		App: AppConfig{
			Name: getEnv("APP_NAME", "go-chat"),
			Env:  getEnv("APP_ENV", EnvDevelopment),
			Port: getEnv("PORT", "8080"),
		},
		Database: DBConfig{
			Host:     os.Getenv("DATABASE_HOST"),
			Port:     os.Getenv("DATABASE_PORT"),
			User:     os.Getenv("DATABASE_USER"),
			Password: os.Getenv("DATABASE_PASSWORD"),
			Name:     os.Getenv("DATABASE_NAME"),
		},
		CORS: CORSConfig{
			AllowedOrigins: splitCSV(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")),
		},
		Auth: AuthConfig{
			JWT_SECRET: os.Getenv("JWT_SECRET"),
		},
	}

	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func (c *Config) validate() error {
	if c.Database.Host == "" {
		return fmt.Errorf("DATABASE_HOST is required")
	}
	if c.Database.Port == "" {
		return fmt.Errorf("DATABASE_PORT is required")
	}
	if c.Database.User == "" {
		return fmt.Errorf("DATABASE_USER is required")
	}
	if c.Database.Password == "" {
		return fmt.Errorf("DATABASE_PASSWORD is required")
	}
	if c.Database.Name == "" {
		return fmt.Errorf("DATABASE_NAME is required")
	}
	if c.Auth.JWT_SECRET == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	return nil
}
