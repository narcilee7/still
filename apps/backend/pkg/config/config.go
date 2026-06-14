package config

import (
	"fmt"
	"os"
	"strings"
)

// Config holds application configuration.
type Config struct {
	AppEnv string

	Port        string
	DatabaseURL string

	// LLM configuration (replaces the old single OpenAIKey).
	LLMProvider string
	LLMModel    string
	LLMAPIKey   string
	LLMBaseURL  string

	S3Endpoint        string
	S3Region          string
	S3Bucket          string
	S3AccessKeyID     string
	S3SecretAccessKey string

	CORSAllowedOrigins string

	ClerkSecretKey string

	SentryDSN string

	OTelExporter    string
	OTelServiceName string
}

// Load loads configuration from environment variables.
func Load() *Config {
	provider := strings.ToLower(getEnv("LLM_PROVIDER", "openai"))

	return &Config{
		AppEnv: getEnv("APP_ENV", "dev"),

		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://localhost/still?sslmode=disable"),

		LLMProvider: provider,
		LLMModel:    getEnv("LLM_MODEL", defaultLLMModel(provider)),
		LLMAPIKey:   getEnv("LLM_API_KEY", getEnv("OPENAI_API_KEY", "")),
		LLMBaseURL:  getEnv("LLM_BASE_URL", ""),

		S3Endpoint:        getEnv("S3_ENDPOINT", ""),
		S3Region:          getEnv("S3_REGION", "us-east-1"),
		S3Bucket:          getEnv("S3_BUCKET", "still-uploads"),
		S3AccessKeyID:     getEnv("S3_ACCESS_KEY_ID", ""),
		S3SecretAccessKey: getEnv("S3_SECRET_ACCESS_KEY", ""),

		CORSAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "*"),

		ClerkSecretKey: getEnv("CLERK_SECRET_KEY", ""),

		SentryDSN: getEnv("SENTRY_DSN", ""),

		OTelExporter:    getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", ""),
		OTelServiceName: getEnv("OTEL_SERVICE_NAME", "still-backend"),
	}
}

// Validate ensures required configuration is present.
func (c *Config) Validate() error {
	if c.LLMProvider == "" {
		return fmt.Errorf("LLM_PROVIDER is required")
	}
	if c.LLMModel == "" {
		return fmt.Errorf("LLM_MODEL is required")
	}
	if c.LLMAPIKey == "" {
		return fmt.Errorf("LLM_API_KEY (or legacy OPENAI_API_KEY) is required")
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

// defaultLLMModel returns the current recommended model for each provider.
// These defaults are reviewed periodically; pin LLM_MODEL in your env file if you
// need a different model or a dated snapshot.
func defaultLLMModel(provider string) string {
	switch provider {
	case "openai":
		// GPT-5.3 Instant is the current everyday ChatGPT / API default.
		// Use gpt-4o-mini if you prefer a smaller, cheaper vision-capable fallback.
		return "gpt-5.3-chat-latest"
	case "deepseek":
		// deepseek-v4-pro is the vision-capable sibling; v4-flash is text-only.
		// The legacy IDs deepseek-chat / deepseek-reasoner retire on 2026-07-24.
		return "deepseek-v4-pro"
	case "moonshot":
		// Kimi K2.6 is the latest multimodal flagship (vision + text + tool use).
		return "kimi-k2.6"
	case "qwen":
		// Qwen3.5 Plus is the native vision-language flagship on DashScope.
		return "qwen3.5-plus"
	default:
		return ""
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
