package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Modd1e/durarun/internal/config"
	"github.com/Modd1e/durarun/internal/logger"
	"github.com/Modd1e/durarun/internal/postgres"
	"github.com/Modd1e/durarun/internal/postgres/dbgen"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	log, cfg, err := run()
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx := context.Background()

	pool, _, err := initDB(ctx, cfg)
	if err != nil {
		log.Error("init db", "error", err)
		return
	}

	log.Info("application started")
	// code here

	defer pool.Close()
}

func run() (*slog.Logger, *config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, nil, fmt.Errorf("load config: %w", err)
	}

	log := logger.New(logger.Config{
		Env:       string(cfg.Environment),
		Level:     "debug",
		AddSource: false,
	})

	slog.SetDefault(log)

	return log, &cfg, nil
}

func initDB(
	ctx context.Context,
	cfg *config.Config,
) (*pgxpool.Pool, *dbgen.Queries, error) {
	pool, err := postgres.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("create postgres pool: %w", err)
	}

	queries := dbgen.New(pool)

	return pool, queries, nil
}
