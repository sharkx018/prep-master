package config

import (
	"os"
	"strings"
)

// Config holds all configuration for the application
type Config struct {
	DatabaseURL   string
	Port          string
	Environment   string
	AuthUsername  string
	AuthPassword  string
	AuthUsers     string // Comma-separated list of usernames
	AuthPasswords string // Comma-separated list of passwords
	JWTSecret     string
}

// Load reads configuration from environment variables
func Load() *Config {
	return &Config{
		DatabaseURL:   getEnv("DATABASE_URL", ""),
		Port:          getEnv("PORT", "8080"),
		Environment:   getEnv("NODE_ENV", "development"),
		AuthUsername:  getEnv("AUTH_USERNAME", "admin"),
		AuthPassword:  getEnv("AUTH_PASSWORD", "password"),
		AuthUsers:     getEnv("AUTH_USERS", ""),
		AuthPasswords: getEnv("AUTH_PASSWORDS", ""),
		JWTSecret:     getEnv("JWT_SECRET", "default_secret_key"),
	}
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// ValidateCredentials checks if the provided username and password are valid
// This method combines both multi-user and single-user authentication
func (c *Config) ValidateCredentials(username, password string) bool {
	// First, check multi-user credentials if they exist
	if c.AuthUsers != "" && c.AuthPasswords != "" {
		users := strings.Split(c.AuthUsers, ",")
		passwords := strings.Split(c.AuthPasswords, ",")

		// Trim spaces
		for i := range users {
			users[i] = strings.TrimSpace(users[i])
		}
		for i := range passwords {
			passwords[i] = strings.TrimSpace(passwords[i])
		}

		// Check if counts match for multi-user
		if len(users) == len(passwords) {
			for i, user := range users {
				if user == username && passwords[i] == password {
					return true
				}
			}
		}
	}

	// Also check single-user credentials (always available as fallback)
	if c.AuthUsername != "" && c.AuthPassword != "" {
		if username == c.AuthUsername && password == c.AuthPassword {
			return true
		}
	}

	return false
}

// GetUsers returns a slice of all valid usernames (for informational purposes)
func (c *Config) GetUsers() []string {
	var users []string

	// Add users from comma-separated AUTH_USERS if present
	if c.AuthUsers != "" {
		multiUsers := strings.Split(c.AuthUsers, ",")
		for _, user := range multiUsers {
			trimmed := strings.TrimSpace(user)
			if trimmed != "" {
				users = append(users, trimmed)
			}
		}
	}

	// Always add the single AUTH_USERNAME as well (if it's not empty and not already in the list)
	if c.AuthUsername != "" {
		// Check if AUTH_USERNAME is already in the multi-user list
		found := false
		for _, user := range users {
			if user == c.AuthUsername {
				found = true
				break
			}
		}
		if !found {
			users = append(users, c.AuthUsername)
		}
	}

	// If no users found, return default
	if len(users) == 0 {
		return []string{"admin"}
	}

	return users
}

// GetPasswords returns a slice of all passwords (for informational purposes)
// Note: This is mainly for testing - in production you shouldn't expose passwords
func (c *Config) GetPasswords() []string {
	var passwords []string

	// Add passwords from comma-separated AUTH_PASSWORDS if present
	if c.AuthPasswords != "" {
		multiPasswords := strings.Split(c.AuthPasswords, ",")
		for _, password := range multiPasswords {
			trimmed := strings.TrimSpace(password)
			if trimmed != "" {
				passwords = append(passwords, trimmed)
			}
		}
	}

	// Always add the single AUTH_PASSWORD as well (if it's not empty)
	if c.AuthPassword != "" {
		passwords = append(passwords, c.AuthPassword)
	}

	// If no passwords found, return default
	if len(passwords) == 0 {
		return []string{"password"}
	}

	return passwords
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
