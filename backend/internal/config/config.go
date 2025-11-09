package config

import "time"

type Environment string

type Config struct {
	Env      Environment
	Server   ServerConfig
	Database DBConfig
	JWT      JWTConfig
	Logger   LoggerConfig
	CORS     CORSConfig
}

type ServerConfig struct {
	Port string
	Host string
}

type DBConfig struct {
	Host        string
	Port        string
	User        string
	Password    string
	DBName      string
	SSLMode     string
	MaxConns    int // максимум соединений в пуле
}

type JWTConfig struct {
	Secret string
	Expiry time.Duration
}

type LoggerConfig struct {
	Level string // "debug", "info", "warn", "error"
	File  string // путь к логу
}

type CORSConfig struct {
	AllowedOrigins []string
}
