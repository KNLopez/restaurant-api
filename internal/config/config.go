package config

import (
	"fmt"
	"os"
)

type Config struct {
	Server     ServerConfig
	Database   DatabaseConfig
	Auth       AuthConfig
	BaseURL    string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	Cloudinary struct {
		CloudName string `env:"CLOUDINARY_CLOUD_NAME"`
		APIKey    string `env:"CLOUDINARY_API_KEY"`
		APISecret string `env:"CLOUDINARY_API_SECRET"`
	}
}

type ServerConfig struct {
	Port int
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type AuthConfig struct {
	JWTSecret         string
	GoogleClientID    string
	GoogleClientKey   string
	FacebookClientID  string
	FacebookClientKey string
}

func Load() (*Config, error) {
	dbPort := 5432     // default postgres port
	serverPort := 8080 // default server port

	return &Config{
		Server: ServerConfig{
			Port: serverPort,
		},
		Database: DatabaseConfig{
			Host:     getEnvOrDefault("DB_HOST", "localhost"),
			Port:     dbPort,
			User:     getEnvOrDefault("DB_USER", "postgres"),
			Password: getEnvOrDefault("DB_PASSWORD", "postgres"),
			DBName:   getEnvOrDefault("DB_NAME", "restaurant_db"),
			SSLMode:  getEnvOrDefault("DB_SSLMODE", "disable"),
		},
		Auth: AuthConfig{
			JWTSecret:         getEnvOrDefault("JWT_SECRET", "your-secret-key"),
			GoogleClientID:    os.Getenv("GOOGLE_CLIENT_ID"),
			GoogleClientKey:   os.Getenv("GOOGLE_CLIENT_SECRET"),
			FacebookClientID:  os.Getenv("FACEBOOK_CLIENT_ID"),
			FacebookClientKey: os.Getenv("FACEBOOK_CLIENT_SECRET"),
		},
	}, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}
