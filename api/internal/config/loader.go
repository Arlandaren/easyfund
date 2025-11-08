package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)


const (
	EnvDevelopment Environment = "development"
	EnvProduction  Environment = "production"
)

// Load загружает конфигурацию из переменных окружения
func Load() (*Config, error) {
	env := Environment(getEnv("ENV", string(EnvDevelopment)))

	// Валидация JWT_SECRET
	jwtSecret := getEnv("JWT_SECRET", "")
	if env == EnvProduction && (jwtSecret == "" || len(jwtSecret) < 32) {
		return nil, fmt.Errorf("JWT_SECRET must be set and at least 32 characters in production")
	}
	if jwtSecret == "" {
		jwtSecret = "dev-secret-key-change-this-in-production"
		if env == EnvProduction {
			log.Fatal("JWT_SECRET is required in production")
		}
	}

	jwtExpiry, err := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	if err != nil {
		jwtExpiry = 24
	}

	dbSSLMode := getEnv("DB_SSLMODE", "disable")
	if env == EnvProduction && dbSSLMode == "disable" {
		log.Println("WARNING: DB_SSLMODE is 'disable' in production. Consider using 'require' or 'verify-full'")
	}

	return &Config{
		Env: env,
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
		},
		Database: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "lk_courses"),
			SSLMode:  dbSSLMode,
			MaxConns: parseIntWithDefault(getEnv("DB_MAX_CONNECTIONS", ""), 25),
		},
		JWT: JWTConfig{
			Secret: jwtSecret,
			Expiry: time.Hour * time.Duration(jwtExpiry),
		},
		Logger: LoggerConfig{
			Level: getEnv("LOG_LEVEL", "info"),
			File:  getEnv("LOG_FILE", "./logs/app.log"),
		},
		CORS: CORSConfig{
			AllowedOrigins: parseStringList(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")),
		},
	}, nil
}

// getEnv получает значение переменной окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseIntWithDefault парсит строку в int или возвращает значение по умолчанию
func parseIntWithDefault(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	if parsed, err := strconv.Atoi(value); err == nil {
		return parsed
	}
	return defaultValue
}

// parseStringList парсит строку в список строк
func parseStringList(value string) []string {
	if value == "" {
		return []string{}
	}
	// Парсим строку вида "http://localhost:3000,https://example.com"
	list := []string{}
	for _, item := range splitString(value, ",") {
		if trimmed := trim(item); trimmed != "" {
			list = append(list, trimmed)
		}
	}
	return list
}

func splitString(s string, sep string) []string {
	var result []string
	var current string
	for _, char := range s {
		if string(char) == sep {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

func trim(s string) string {
	return os.Expand(s, func(key string) string {
		return ""
	})
	// Упрощенная версия:
	// return strings.TrimSpace(s)
}
