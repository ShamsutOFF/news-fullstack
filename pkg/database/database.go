package database

import (
	"context"
	"log/slog"
	"news-fullstack/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateDbPool(config *config.DatabaseConfig, logger *slog.Logger) *pgxpool.Pool {
	logger.Info("Connecting to database...")

	dbpool, err := pgxpool.New(context.Background(), config.URL)
	if err != nil {
		logger.Error("Failed to connect to database", err)
		panic(err)
	}

	// Проверим подключение
	if err := dbpool.Ping(context.Background()); err != nil {
		logger.Error("Database ping failed", err)
		panic(err)
	}

	logger.Info("Database connection established successfully")
	return dbpool
}
