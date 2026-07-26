package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Thomika1/KarHub/internal/api"
	"github.com/Thomika1/KarHub/internal/container"
)

func main() {
	startedAt := time.Now()

	ctx := context.Background()
	ctx = context.WithValue(ctx, "startedAt", startedAt)

	// Init container
	dep, err := container.New(ctx)
	if err != nil {
		dep.Components.Logger.Error("Cannot initialize container", "err", err)
		os.Exit(1)
	}

	logger := dep.Components.Logger
	logger.Info("Starting API...", "startedAt", startedAt)

	// Get port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		logger.Warn("No PORT environment variable found, using default port", "port", port)
	}

	// Init server
	handler := api.Handler(ctx, dep)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	go func() {
		logger.Info("Server listening on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed", "err", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server forced to shutdown", "err", err)
		os.Exit(1)
	}

	logger.Info("Server exited")
}
