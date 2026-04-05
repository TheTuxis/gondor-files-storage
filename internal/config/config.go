package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port          string
	DatabaseURL   string
	JWTSecret     string
	RedisURL      string
	NatsURL       string
	S3Endpoint    string
	S3AccessKey   string
	S3SecretKey   string
	S3Bucket      string
	S3Region      string
	MaxUploadSize int64
	LogLevel      string
	Environment   string
}

func Load() *Config {
	return &Config{
		Port:          getEnv("PORT", "8006"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://gondor:gondor_dev@localhost:5432/gondor_files?sslmode=disable"),
		JWTSecret:     getEnv("JWT_SECRET", ""),
		RedisURL:      getEnv("REDIS_URL", "redis://localhost:6379"),
		NatsURL:       getEnv("NATS_URL", "nats://localhost:4222"),
		S3Endpoint:    getEnv("S3_ENDPOINT", "http://localhost:9000"),
		S3AccessKey:   getEnv("S3_ACCESS_KEY", "minioadmin"),
		S3SecretKey:   getEnv("S3_SECRET_KEY", "minioadmin"),
		S3Bucket:      getEnv("S3_BUCKET", "gondor-files"),
		S3Region:      getEnv("S3_REGION", "us-east-1"),
		MaxUploadSize: getEnvInt64("MAX_UPLOAD_SIZE", 104857600),
		LogLevel:      getEnv("LOG_LEVEL", "info"),
		Environment:   getEnv("ENVIRONMENT", "development"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt64(key string, fallback int64) int64 {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return i
		}
	}
	return fallback
}
