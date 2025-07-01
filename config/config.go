package config

import (
	"log"
	"os"
)

type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	AppPort    string
}

func LoadConfig() *Config {
	// Set default values
	cfg := &Config{
		DBHost:     getEnv("DB_HOST", "PostgreSQL"),
		DBUser:     getEnv("DB_USER", "admin"),
		DBPassword: getEnv("DB_PASSWORD", "admin@123456"),
		DBName:     getEnv("DB_NAME", "PhotoKit"),
		DBPort:     getEnv("DB_PORT", "5432"),
		AppPort:    getEnv("PORT", "8080"),
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		if defaultValue == "" {
			log.Fatalf("Environment variable %s not set", key)
		}
		return defaultValue
	}
	return value
}
