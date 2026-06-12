package main

import (
	"log/slog"

	"github.com/Modd1e/durarun/internal/logger"
)

func main() {
	log := logger.New(logger.Config{
		Env:       "dev",
		Level:     "debug",
		AddSource: false,
	})

	slog.SetDefault(log)

	log.Info("Hello world")
}
