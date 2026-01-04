package config

import (
	"os"
	"strconv"
)

type Config struct {
	DBHost    string
	DBPort    int
	DBUser    string
	DBPass    string
	DBName    string
	DBSSLMode string
	ApiPort   int
}

func Load() Config {
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))
	apiPort, _ := strconv.Atoi(getEnv("GOLLAB_API_PORT", "8080"))
	return Config{
		DBHost:    getEnv("DB_HOST", "localhost"),
		DBPort:    dbPort,
		DBUser:    getEnv("DB_USER", "postgres"),
		DBPass:    getEnv("DB_PASS", "postgres"),
		DBName:    getEnv("DB_NAME", "gollab_db"),
		DBSSLMode: getEnv("DB_SSL_MODE", "disable"),
		ApiPort:   apiPort,
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
