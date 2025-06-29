package config

import (
	"os"
)

// Config holds all configuration for the application
type Config struct {
	DatabaseURL  string
	Port         string
	Environment  string
	AuthUsername string
	AuthPassword string
	JWTSecret    string
}

// Load reads configuration from environment variables
func Load() *Config {
	return &Config{
		DatabaseURL:  getEnv("DATABASE_URL", ""),
		Port:         getEnv("PORT", "8080"),
		Environment:  getEnv("NODE_ENV", "development"),
		AuthUsername: getEnv("AUTH_USERNAME", "admin"),
		AuthPassword: getEnv("AUTH_PASSWORD", "password"),
		JWTSecret:    getEnv("JWT_SECRET", "default_secret_key"),
	}
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
