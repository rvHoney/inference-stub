/*
Main entry point for the inference-stub server.
It handles CLI flags and starts the TCP server.
*/
package main

import (
	"log/slog"
	"os"

	"github.com/rvHoney/inference-stub/internal/config"
)

func main() {
	cfg, err := config.Parse(os.Args[1:])
	if err != nil {
		slog.Error("failed to parse config", "error", err)
		os.Exit(1)
	}

	slog.Info("inference-stub initialized", "config", cfg)
}
