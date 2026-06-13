package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	"github.com/still-mvp/still/apps/backend/internal/server"
	"github.com/still-mvp/still/apps/backend/pkg/config"
)

func main() {
	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := server.New(":" + cfg.Port)

	if err := srv.Start(ctx); err != nil {
		log.Fatal().Err(err).Msg("server failed")
	}
}
