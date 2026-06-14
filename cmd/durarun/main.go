package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Modd1e/durarun/internal/api"
	"github.com/Modd1e/durarun/internal/config"
	"github.com/Modd1e/durarun/internal/logger"
	"github.com/Modd1e/durarun/internal/postgres"
	"github.com/Modd1e/durarun/internal/postgres/dbgen"
)

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("load config: %w", err)
		return
	}

	// Initialize logger
	log := logger.New(logger.Config{
		Env:       string(cfg.Environment),
		Level:     "debug",
		AddSource: false,
	})

	slog.SetDefault(log)

	// Initialize database
	ctx := context.Background()

	pool, err := postgres.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		fmt.Errorf("create postgres pool: %w", err)
		return
	}

	queries := dbgen.New(pool)

	// Initialize API
	serverAPI := api.New(queries, log)

	server := &http.Server{
		Addr:    ":3000",
		Handler: serverAPI.Handler(),
	}

	log.Info("application started", "address", server.Addr)

	if err := server.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {
		log.Error("HTTP server stopped", "error", err)
	}

	defer pool.Close()
}
