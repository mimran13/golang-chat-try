package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"go-chat/internal/config"
	"go-chat/internal/shared"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	cfg , err := config.LoadConfig()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	setLogger(cfg.App.Env)
	db := setupDB(*cfg)
	defer db.Close()

	c := NewContainer(db)

	var shuttingDown atomic.Bool

	server := &http.Server{
		Addr:              ":" + cfg.App.Port,
		Handler:           setupRoutes(c, &shuttingDown),
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	<-ctx.Done()
	slog.Info("shutdown started")
	shuttingDown.Store(true)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	slog.Info("stopping http server")
	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
		server.Close()
	}

	slog.Info("shutdown complete")
}

func setLogger(appEnv string) {
	level := slog.LevelInfo

	if appEnv == shared.EnvDevelopment {
		level = slog.LevelDebug
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{                                                                                                                                                
      Level: level,                       
  	}))

 	 slog.SetDefault(logger)
}

func setupDB(cfg config.Config) *sql.DB {
	runMigrations(cfg.Database)
	db := connectDB(cfg.Database)
	return db;
}