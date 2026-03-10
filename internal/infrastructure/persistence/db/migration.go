package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/thanvuc/go-core-lib/log"
	"go.uber.org/zap"
)

func NewMigrations(
	db *pgxpool.Pool,
	logger log.LoggerV2,
) {
	// Folder where your migration files (.sql) are stored
	migrationsDir := "internal/infrastructure/persistence/db/sql/schema"

	// Convert *pgxpool.Pool to *sql.DB using stdlib
	sqlDB := stdlib.OpenDBFromPool(db)

	if err := goose.Up(sqlDB, migrationsDir); err != nil {
		logger.Error("Failed to apply migrations:", log.WithFields(zap.Error(err)))
		return
	}

	logger.Info("Migrations applied successfully from")
}
