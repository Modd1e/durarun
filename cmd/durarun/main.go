package main

import (
	"log/slog"

	"github.com/Modd1e/durarun/internal/config"
	"github.com/Modd1e/durarun/internal/logger"
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

	log.Info("Hello world")
}
