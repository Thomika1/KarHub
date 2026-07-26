package config

import (
	"log/slog"
	"os"
)

var appLogger *slog.Logger

func InitLogger() *slog.Logger {
	appLogger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	return appLogger
}
