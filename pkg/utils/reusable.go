package utils

import (
	"log/slog"
	"os"
)

// SetupLogger Устанавливает логгер
func SetupLogger() *slog.Logger {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	return log
}
