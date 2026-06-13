package config

import (
	"fmt"
	"os"
)

// Config holds application configuration.
type Config struct {
	AppEnv string

	Port        string
	DatabaseURL string

	OpenAIKey string

	S3Endpoint        string
	S3Region          string
	S3Bucket          string
	S3AccessKeyID     string
	S3SecretAccessKey string

	CORSAllowedOrigins string

	ClerkSecretKey string
}

// Load loads configuration from environment variables.
func Load() *Config {
	return &Config{
		AppEnv: getEnv("APP_ENV", "dev"),

		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://localhost/still?sslmode=disable"),

		OpenAIKey: getEnv("OPENAI_API_KEY", ""),

		S3Endpoint:        getEnv("S3_ENDPOINT", ""),
		S3Region:          getEnv("S3_REGION", "us-east-1"),
		S3Bucket:          getEnv("S3_BUCKET", "still-uploads"),
		S3AccessKeyID:     getEnv("S3_ACCESS_KEY_ID", ""),
		S3SecretAccessKey: getEnv("S3_SECRET_ACCESS_KEY", ""),

		CORSAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "*"),

		ClerkSecretKey: getEnv("CLERK_SECRET_KEY", ""),
	}
}

// Validate ensures required configuration is present.
func (c *Config) Validate() error {
	if c.OpenAIKey == "" {
		return fmt.Errorf("OPENAI_API_KEY is required")
	}
	if c.S3Bucket == "" {
		return fmt.Errorf("S3_BUCKET is required")
	}
	if c.S3Region == "" {
		return fmt.Errorf("S3_REGION is required")
	}
	if c.ClerkSecretKey == "" {
		return fmt.Errorf("CLERK_SECRET_KEY is required")
	}
	if c.AppEnv == "prod" {
		if c.S3AccessKeyID == "" || c.S3SecretAccessKey == "" {
			return fmt.Errorf("S3_ACCESS_KEY_ID and S3_SECRET_ACCESS_KEY are required in prod")
		}
	}
	return nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
