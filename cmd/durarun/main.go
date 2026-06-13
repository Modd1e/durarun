package main

import (
	"context"
	"log/slog"

	"github.com/Modd1e/durarun/internal/config"
	"github.com/Modd1e/durarun/internal/logger"
	"github.com/Modd1e/durarun/internal/postgres"
	"github.com/Modd1e/durarun/internal/postgres/dbgen"
)

func main() {
	log := logger.New(logger.Config{
		Env:       "dev",
		Level:     "debug",
		AddSource: false,
	})

	slog.SetDefault(log)

	config, err := config.Load()
	if err != nil {
		log.Error("load config: %v", err)
	}

	ctx := context.Background()

	pool, err := postgres.NewPool(ctx, config.DatabaseURL)
	if err != nil {
		log.Error(err.Error())
	}
	defer pool.Close()

	queries := dbgen.New(pool)

	payload := `{"task":"example"}`
	queue := int32(1)

	job, err := queries.CreateJob(ctx, dbgen.CreateJobParams{
		Queue:   &queue,
		Payload: &payload,
	})
	if err != nil {
		log.Error("create job: %w", err.Error())
	}

	slog.Info("job created", "job_id", job.ID)

	log.Info("Hello world")
}
