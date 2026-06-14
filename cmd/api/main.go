package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"go-chat/internal/config"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	setLogger(cfg.App.Env)
	db, err := setupDB(*cfg)
	if err != nil {
		slog.Error("db setup failed", "error", err)
		os.Exit(1)
	}
	defer func() { _ = db.Close() }()

	c := NewContainer(db, cfg)

	var shuttingDown atomic.Bool

	server := &http.Server{
		Addr:              ":" + cfg.App.Port,
		Handler:           setupRoutes(c, &shuttingDown, cfg.CORS.AllowedOrigins),
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	slog.Info("Server starting on port..." + cfg.App.Port)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("http server", "error", err)
			stop()
		}
	}()

	<-ctx.Done()
	slog.Info("shutdown started")
	shuttingDown.Store(true)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	slog.Info("stopping http server")
	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
		_ = server.Close()
	}

	slog.Info("shutdown complete")
}

func setLogger(appEnv string) {
	level := slog.LevelInfo

	if appEnv == config.EnvDevelopment {
		level = slog.LevelDebug
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))

	slog.SetDefault(logger)
}

func setupDB(cfg config.Config) (*sql.DB, error) {
	db, err := connectDB(cfg.Database)
	if err != nil {
		return nil, err
	}
	return db, nil
}
