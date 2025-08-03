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
	Env      string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
	URL      string // Para suportar DATABASE_URL do Neon
}

type ServerConfig struct {
	Port int
	Host string
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

	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	config := &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     dbPort,
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "tabler_api"),
			SSLMode:  getEnv("DB_SSL_MODE", "require"),
			URL:      getEnv("DATABASE_URL", ""),
		},
		Server: ServerConfig{
			Port: port,
			Host: getEnv("SERVER_HOST", "localhost"),
		},
		Env: getEnv("ENV", "development"),
	}

	return config, nil
}

func (c *Config) GetDatabaseURL() string {
	// Se DATABASE_URL estiver definida, use ela diretamente
	if c.Database.URL != "" {
		return c.Database.URL
	}

	// Caso contrário, construa a URL a partir das variáveis individuais
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 