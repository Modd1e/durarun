package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Modd1e/durarun/internal/config"
	"github.com/Modd1e/durarun/internal/logger"
	"github.com/Modd1e/durarun/internal/postgres"
	"github.com/Modd1e/durarun/internal/postgres/dbgen"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/jobs/count/", JobsCountHandler(queries, ctx))
	r.Post("/jobs/create/", CreateJobHandler(queries, ctx))

	log.Info("application started")
	http.ListenAndServe(":3000", r)

	defer pool.Close()
}

func CreateJobHandler(queries *dbgen.Queries, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload dbgen.CreateJobParams
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		queries.CreateJob(ctx, payload)
	}
}

func JobsCountHandler(queries *dbgen.Queries, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		count, _ := queries.CountJobs(ctx)
		w.Write([]byte(fmt.Sprintf("%d", count)))
	}
}
