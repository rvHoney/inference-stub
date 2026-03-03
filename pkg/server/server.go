// Package server handles TCP communications
package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rvHoney/inference-stub/pkg/api"
)

// Server holds the application state and configuration.
type Server struct {
	httpServer *http.Server
}

// Init initializes a new Server and mounts the multiplexer.
func Init(port int, timeout, ttft, tpot time.Duration) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", api.HealthCheckHandler)

	handlerCfg := api.HandlerConfig{
		TTFT: ttft,
		TPOT: tpot,
	}
	mux.HandleFunc("POST /v1/chat/completions", api.ChatCompletionsHandler(handlerCfg))

	httpSrv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      mux,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	}

	return &Server{
		httpServer: httpSrv,
	}
}

// Start starts listening on the configured HTTP server.
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
