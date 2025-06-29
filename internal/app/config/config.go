package config

import (
	"log"
	"os"
)

var (
	defaultHTTPPort       = "8080"
	defaultDBPath         = "tasks.db"
	defaultJWTPrivateKey  = "123"
	defaultJWTIssuer      = "123"
	defaultJWTTokenExpiry = "3600s"
)

type Config struct {
	Port           string
	DBPath         string
	JWTPrivateKey  string
	JWTIssuer      string
	JWTTokenExpiry string
}

func Read() *Config {
	cfg := &Config{
		Port:           getEnvOrDefault("PORT", defaultHTTPPort),
		DBPath:         getEnvOrDefault("DB_PATH", defaultDBPath),
		JWTPrivateKey:  getJWTPrivateKey(),
		JWTIssuer:      getEnvOrDefault("JWT_ISSUER", defaultJWTIssuer),
		JWTTokenExpiry: getEnvOrDefault("JWT_TOKEN_EXPIRY", defaultJWTTokenExpiry),
	}

	return cfg
}

func getJWTPrivateKey() string {
	if keyFile := os.Getenv("JWT_PRIVATE_KEY_FILE"); keyFile != "" {
		keyData, err := os.ReadFile(keyFile)
		if err != nil {
			log.Printf("Warning: Failed to read JWT private key file %s: %v", keyFile, err)
			return getEnvOrDefault("JWT_PRIVATE_KEY", defaultJWTPrivateKey)
		}
		return string(keyData)
	}

	return getEnvOrDefault("JWT_PRIVATE_KEY", defaultJWTPrivateKey)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
