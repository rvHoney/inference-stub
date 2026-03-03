// Package api contains HTTP handlers for the Inference-Stub endpoints
package api

import (
	"log/slog"
	"net/http"
	"time"
)

// HandlerConfig holds configuration parameters required by endpoints.
type HandlerConfig struct {
	TTFT time.Duration
	TPOT time.Duration
}

// HealthCheckHandler returns an HTTP 200 to indicate the server is alive.
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// ChatCompletionsHandler handles chat requests with SSE response.
func ChatCompletionsHandler(cfg HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Chat request received", "method", r.Method, "path", r.URL.Path)

		// SSE streaming headers
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// Send the headers over the wire immediately
		w.WriteHeader(http.StatusOK)

		slog.Debug("Connection accepted, holding stream based on mock parameters", "ttft", cfg.TTFT)

		// TODO: Implement Lorem Ipsum token generation looping with TTFT and TPOT delays
	}
}
