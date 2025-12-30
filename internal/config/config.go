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
}

func Load() Config {
	port, _ := strconv.Atoi(getEnv("GOLLAB_DB_PORT", "5432"))
	return Config{
		DBHost:    getEnv("GOLLAB_DB_HOST", "localhost"),
		DBPort:    port,
		DBUser:    getEnv("GOLLAB_DB_USERNAME", "postgres"),
		DBPass:    getEnv("GOLLAB_DB_PASS", "postgres"),
		DBName:    getEnv("GOLLAB_DB_NAME", "gollab_db"),
		DBSSLMode: getEnv("GOLLAB_DB_SSL_MODE", "disable"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
