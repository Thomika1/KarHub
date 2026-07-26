package config

import (
	"log/slog"
	"os"
)

var appLogger *slog.Logger

func InitLogger() *slog.Logger {
	appLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	slog.SetDefault(appLogger)

	return appLogger
}
