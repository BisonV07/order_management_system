package config

import (
	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	JWT      JWTConfig
	Logging  LoggingConfig
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret string
	Expiry string
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level string
}

// Load loads configuration from environment variables and config files
func Load() (*Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 5432)
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("DB_NAME", "oms_db")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("LOG_LEVEL", "info")

	if err := viper.ReadInConfig(); err != nil {
		// Config file not found; use defaults and env vars
	}

	return &Config{
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetInt("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Name:     viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSLMODE"),
		},
		Server: ServerConfig{
			Port: viper.GetString("SERVER_PORT"),
		},
		JWT: JWTConfig{
			Secret: viper.GetString("JWT_SECRET"),
			Expiry: viper.GetString("JWT_EXPIRY"),
		},
		Logging: LoggingConfig{
			Level: viper.GetString("LOG_LEVEL"),
		},
	}, nil
}

