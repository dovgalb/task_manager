package custom_loggers

import (
	"log/slog"
	"os"
)

func SetupLogger() *slog.Logger {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	return log
}
