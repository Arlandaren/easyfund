package config

import (
	"os"
	"strconv"
)

type Config struct {
	// Server
	Port        string
	Environment string

	// Virtual Bank API
	VBankBaseURL       string
	VBankClientID      string
	VBankClientSecret  string
	RequestingBankName string

	// JWT
	JWTSecret      string
	JWTExpiryHours int

	// CORS
	CORSAllowedOrigins []string
	CORSAllowedMethods []string
	CORSAllowedHeaders []string
}

func LoadConfig() *Config {
	return &Config{
		// Server
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),

		// Virtual Bank API
		VBankBaseURL:       getEnv("VBANK_BASE_URL", "https://vbank.open.bankingapi.ru"),
		VBankClientID:      getEnv("VBANK_CLIENT_ID", "team080"),
		VBankClientSecret:  getEnv("VBANK_CLIENT_SECRET", "dpzKsUNyd5PSMmqk8vEA4FBA4lu1XdzK"),
		RequestingBankName: getEnv("REQUESTING_BANK_NAME", "EasyFund Consortium Platform"),

		// JWT
		JWTSecret:      getEnv("JWT_SECRET", "your-super-secret-jwt-key-for-easyfund-2024"),
		JWTExpiryHours: getEnvAsInt("JWT_EXPIRY_HOURS", 24),

		// CORS
		CORSAllowedOrigins: []string{"*"},
		CORSAllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		CORSAllowedHeaders: []string{"Content-Type", "Authorization", "X-Requested-With"},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
