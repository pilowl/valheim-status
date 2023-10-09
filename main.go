package main

import (
	"os"

	"golang.org/x/exp/slog"
)

func main() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logHandler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(logHandler)

	config, err := NewConfig()
	if err != nil {
		logger.With(err).Info("Failed to initialize config")
		os.Exit(1)
	}
	logger.Info("Service config: " + config.String())

	logger.Info("Starting server now...")
	if err := run(logger, config); err != nil {
		logger.With(map[string]interface{}{
			"error":  err,
			"config": config,
		}).Error("Error occured while running server")
	}

	logger.Info("Server exited")
}
