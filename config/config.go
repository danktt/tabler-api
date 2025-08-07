package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Database   DatabaseConfig
	Server     ServerConfig
	BetterAuth BetterAuthConfig
	CORS       CORSConfig
	Env        string
}

type DatabaseConfig struct {
	URL string // DATABASE_URL do Neon
}

type ServerConfig struct {
	Port int
	Host string
}

type BetterAuthConfig struct {
	BaseURL     string
	Audience    string
	Issuer      string
	JWKSEndpoint string
}

type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
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

	betterAuthURL := getEnv("BETTER_AUTH_URL", "http://localhost:3000")
	if betterAuthURL == "" {
		return nil, fmt.Errorf("BETTER_AUTH_URL is required")
	}

	// Configuração de CORS
	corsOrigins := getEnv("CORS_ALLOWED_ORIGINS", "*")
	corsMethods := getEnv("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,OPTIONS,PATCH")
	corsHeaders := getEnv("CORS_ALLOWED_HEADERS", "Origin,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization,X-Requested-With")

	config := &Config{
		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", ""),
		},
		Server: ServerConfig{
			Port: port,
			Host: getEnv("SERVER_HOST", "localhost"),
		},
		BetterAuth: BetterAuthConfig{
			BaseURL:      betterAuthURL,
			Audience:     betterAuthURL,
			Issuer:       betterAuthURL,
			JWKSEndpoint: fmt.Sprintf("%s/api/auth/jwks", betterAuthURL),
		},
		CORS: CORSConfig{
			AllowedOrigins: strings.Split(corsOrigins, ","),
			AllowedMethods: strings.Split(corsMethods, ","),
			AllowedHeaders: strings.Split(corsHeaders, ","),
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