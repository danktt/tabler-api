package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Auth0    Auth0Config
	Env      string
}

type DatabaseConfig struct {
	URL string // DATABASE_URL do Neon
}

type ServerConfig struct {
	Port int
	Host string
}

type Auth0Config struct {
	Domain        string
	Audience      string
	Issuer        string
	JWKSEndpoint  string
}

func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// Don't return error if .env doesn't exist, use environment variables
	}

	port, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
	}

	auth0Domain := getEnv("AUTH0_DOMAIN", "")
	if auth0Domain == "" {
		return nil, fmt.Errorf("AUTH0_DOMAIN is required")
	}

	config := &Config{
		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", ""),
		},
		Server: ServerConfig{
			Port: port,
			Host: getEnv("SERVER_HOST", "localhost"),
		},
		Auth0: Auth0Config{
			Domain:       auth0Domain,
			Audience:     getEnv("AUTH0_AUDIENCE", ""),
			Issuer:       fmt.Sprintf("https://%s/", auth0Domain),
			JWKSEndpoint: fmt.Sprintf("https://%s/.well-known/jwks.json", auth0Domain),
		},
		Env: getEnv("ENV", "development"),
	}

	// Validate that DATABASE_URL is provided
	if config.Database.URL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required. Please configure your Neon database URL")
	}

	return config, nil
}

func (c *Config) GetDatabaseURL() string {
	return c.Database.URL
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 