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
	config, err := NewMockConfig()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	logger.Info("Service config: " + config.String())
	logger.Info("Starting service now...")
}
