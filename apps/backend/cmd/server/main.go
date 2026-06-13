package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"

	"github.com/still-mvp/still/apps/backend/internal/ai"
	"github.com/still-mvp/still/apps/backend/internal/db"
	"github.com/still-mvp/still/apps/backend/internal/server"
	"github.com/still-mvp/still/apps/backend/internal/storage"
	"github.com/still-mvp/still/apps/backend/pkg/config"
)

func main() {
	loadEnv()

	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatal().Err(err).Msg("config validation failed")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := db.AutoMigrate(cfg); err != nil {
		log.Fatal().Err(err).Msg("migrate failed")
	}

	pool, err := db.NewPool(ctx, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("db pool failed")
	}
	defer pool.Close()

	if err := db.Seed(ctx, pool); err != nil {
		log.Fatal().Err(err).Msg("seed failed")
	}

	store, err := storage.NewS3Store(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("s3 store failed")
	}

	analyzer := ai.NewOpenAIAnalyzer(cfg.OpenAIKey)

	srv := server.New(":"+cfg.Port, pool, analyzer, store, cfg)

	if err := srv.Start(ctx); err != nil {
		log.Fatal().Err(err).Msg("server failed")
	}
}

func loadEnv() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	fileMap := map[string]string{
		"dev":  ".env.development",
		"prod": ".env.production",
	}
	files := []string{fileMap[env], ".env"}
	for _, f := range files {
		if f == "" {
			continue
		}
		if _, err := os.Stat(f); err == nil {
			if err := godotenv.Load(f); err != nil {
				log.Warn().Err(err).Str("file", f).Msg("failed to load env file")
			}
		}
	}
}
