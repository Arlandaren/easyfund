package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	PostgresURL  string
	RedisAddr    string
	RedisPass    string
	GRPCAddress  string
	HTTPAddress  string
	JWTSecret    string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		PostgresURL: getEnv("PG_STRING", "postgresql://easyfund:easyfund123@localhost:5434/easyfund_db?sslmode=disable"),
		RedisAddr:   getEnv("REDIS_ADDRESS", "localhost:6378"),
		RedisPass:   getEnv("REDIS_PASSWORD", "redis123"),
		GRPCAddress: getEnv("GRPC_ADDRESS", ":9000"),
		HTTPAddress: getEnv("HTTP_ADDRESS", ":8080"),
		JWTSecret:   getEnv("JWT_SECRET", "your-super-secret-jwt-key"),
	}, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
