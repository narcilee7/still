package config

import (
	"os"
)

// Config holds application configuration.
type Config struct {
	Port           string
	DatabaseURL    string
	R2Endpoint     string
	R2Bucket       string
	OpenAIKey      string
	StorageBaseURL string
	UploadDir      string
}

// Load loads configuration from environment variables.
func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", "8080"),
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://localhost/still?sslmode=disable"),
		R2Endpoint:     getEnv("R2_ENDPOINT", ""),
		R2Bucket:       getEnv("R2_BUCKET", ""),
		OpenAIKey:      getEnv("OPENAI_API_KEY", ""),
		StorageBaseURL: getEnv("STORAGE_BASE_URL", "http://localhost:8080"),
		UploadDir:      getEnv("UPLOAD_DIR", "uploads"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
