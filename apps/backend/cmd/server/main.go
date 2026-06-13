package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	"github.com/still-mvp/still/apps/backend/internal/ai"
	"github.com/still-mvp/still/apps/backend/internal/db"
	"github.com/still-mvp/still/apps/backend/internal/server"
	"github.com/still-mvp/still/apps/backend/internal/storage"
	"github.com/still-mvp/still/apps/backend/pkg/config"
)

func main() {
	cfg := config.Load()

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

	var analyzer ai.Analyzer = &ai.StubAnalyzer{}
	if cfg.OpenAIKey != "" {
		analyzer = ai.NewOpenAIAnalyzer(cfg.OpenAIKey)
	}

	store := storage.NewLocalStore(cfg.StorageBaseURL, cfg.UploadDir)

	srv := server.New(":"+cfg.Port, pool, analyzer, store, cfg)

	if err := srv.Start(ctx); err != nil {
		log.Fatal().Err(err).Msg("server failed")
	}
}
