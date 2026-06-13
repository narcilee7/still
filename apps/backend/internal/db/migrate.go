package db

import (
	"fmt"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/still-mvp/still/apps/backend/pkg/config"
)

// AutoMigrate runs all up migrations.
func AutoMigrate(cfg *config.Config) error {
	dbURL := toPgx5MigrationURL(cfg.DatabaseURL)
	m, err := migrate.New("file://migrations", dbURL)
	if err != nil {
		return fmt.Errorf("migrate init failed: %w", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrate up failed: %w", err)
	}
	return nil
}

func toPgx5MigrationURL(databaseURL string) string {
	scheme := "pgx5://"
	s := databaseURL
	if strings.HasPrefix(s, "postgres://") {
		s = scheme + s[len("postgres://"):]
	} else if strings.HasPrefix(s, "postgresql://") {
		s = scheme + s[len("postgresql://"):]
	}
	if strings.Contains(s, "?") {
		s += "&x-migrations-table=migrations"
	} else {
		s += "?x-migrations-table=migrations"
	}
	return s
}
