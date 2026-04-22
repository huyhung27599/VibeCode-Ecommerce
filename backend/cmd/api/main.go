package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vibecode/ecommerce/backend/internal/config"
	"github.com/vibecode/ecommerce/backend/internal/server"
	"github.com/vibecode/ecommerce/backend/pkg/database"
	"github.com/vibecode/ecommerce/backend/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	log := logger.New(cfg.App.Env)
	slog.SetDefault(log)

	db, err := database.NewPostgres(cfg.Database, cfg.App.Env)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	if cfg.App.AutoMigrate {
		if err := database.AutoMigrate(db); err != nil {
			log.Error("auto-migration failed", "error", err)
			os.Exit(1)
		}
		log.Info("database migrations applied")
	}

	rdb, err := database.NewRedis(cfg.Redis)
	if err != nil {
		log.Warn("failed to connect to redis", "error", err)
	}

	app := server.New(cfg, log, db, rdb)

	srv := &http.Server{
		Addr:              ":" + cfg.App.Port,
		Handler:           app.Router(),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	go func() {
		log.Info("server starting", "port", cfg.App.Port, "env", cfg.App.Env)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("forced shutdown", "error", err)
	}

	if sqlDB, err := db.DB(); err == nil {
		_ = sqlDB.Close()
	}
	if rdb != nil {
		_ = rdb.Close()
	}

	log.Info("server stopped")
}
