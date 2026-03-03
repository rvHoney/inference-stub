/*
Main entry point for the inference-stub server.
It handles CLI flags and starts the TCP server.
*/
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rvHoney/inference-stub/internal/config"
	"github.com/rvHoney/inference-stub/internal/logger"
	"github.com/rvHoney/inference-stub/pkg/server"
)

func main() {
	cfg, err := config.Parse(os.Args[1:])
	if err != nil {
		os.Exit(1)
	}

	logger.Init(cfg)

	slog.Debug("inference-stub initialized", "config", cfg)

	srv := server.Init(cfg.Port, cfg.Timeout, cfg.TTFT, cfg.TPOT)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("Starting server", "port", cfg.Port)
		if err := srv.Start(); err != nil {
			slog.Error("Server failed", "error", err)
		}
	}()

	// Interruption handling
	<-ctx.Done()
	slog.Info("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("Graceful shutdown failed", "error", err)
	} else {
		slog.Info("Server exited gracefully")
	}
}
