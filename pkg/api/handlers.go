// Package api contains HTTP handlers for the Inference-Stub endpoints
package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/robin-vidal/inference-stub/pkg/lorem"
)

// HandlerConfig holds configuration parameters required by endpoints.
type HandlerConfig struct {
	TTFT      time.Duration
	TPOT      time.Duration
	Generator *lorem.Generator
}

// HealthCheckHandler returns an HTTP 200 to indicate the server is alive.
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// ChatCompletionsHandler handles chat requests with SSE response.
func ChatCompletionsHandler(cfg HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Chat request received", "method", r.Method, "path", r.URL.Path)

		var req ChatCompletionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			slog.Error("Failed to decode request body", "error", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		slog.Debug("Parsed ChatCompletionRequest", "model", req.Model, "stream", req.Stream, "messages_count", len(req.Messages))

		// Simulate TTFT
		time.Sleep(cfg.TTFT)

		if req.Stream {
			handleStreamResponse(w, req, cfg)
		} else {
			handleNonStreamResponse(w, req, cfg)
		}
	}
}

func handleStreamResponse(w http.ResponseWriter, req ChatCompletionRequest, cfg HandlerConfig) {
	// SSE streaming headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)

	slog.Debug("Connection accepted, holding stream based on mock parameters", "tpot", cfg.TPOT)

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Send initial empty chunk with assistant role
	sendStreamEvent(w, flusher, createStreamChunk(req.Model, ChatDelta{Role: "assistant"}, nil))

	// Generate N tokens of Lorem Ipsum separated by TPOT
	tokens := cfg.Generator.GenerateTokens()

	for _, token := range tokens {
		time.Sleep(cfg.TPOT)
		sendStreamEvent(w, flusher, createStreamChunk(req.Model, ChatDelta{Content: token}, nil))
	}

	// Send final chunk w/ finish_reason
	time.Sleep(cfg.TPOT)
	stopReason := "stop"
	sendStreamEvent(w, flusher, createStreamChunk(req.Model, ChatDelta{}, &stopReason))

	// Send termination signal
	w.Write([]byte("data: [DONE]\n\n"))
	flusher.Flush()
}

// createStreamChunk build the ChatCompletionStreamResponse boilerplate
func createStreamChunk(model string, delta ChatDelta, finishReason *string) ChatCompletionStreamResponse {
	return ChatCompletionStreamResponse{
		ID:                "chatcmpl-mock-stream123",
		Object:            "chat.completion.chunk",
		Created:           time.Now().Unix(),
		Model:             model,
		SystemFingerprint: "fp_mock_stub",
		Choices: []StreamChoice{
			{
				Index:        0,
				Delta:        delta,
				FinishReason: finishReason,
			},
		},
	}
}

func handleNonStreamResponse(w http.ResponseWriter, req ChatCompletionRequest, cfg HandlerConfig) {
	tokens := cfg.Generator.GenerateTokens()
	tokenCount := len(tokens)
	tpotDelay := time.Duration(tokenCount) * cfg.TPOT
	time.Sleep(tpotDelay)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	content := strings.Join(tokens, "")

	resp := ChatCompletionResponse{
		ID:                "chatcmpl-mock123",
		Object:            "chat.completion",
		Created:           time.Now().Unix(),
		Model:             req.Model,
		SystemFingerprint: "fp_mock_stub",
		Choices: []Choice{
			{
				Index: 0,
				Message: ChatMessage{
					Role:    "assistant",
					Content: content,
				},
				LogProbs:     nil,
				FinishReason: "stop",
			},
		},
		Usage: Usage{
			PromptTokens:     19,
			CompletionTokens: tokenCount,
			TotalTokens:      19 + tokenCount,
		},
		ServiceTier: "default",
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("Failed to encode response", "error", err)
	}
}

// sendStreamEvent convert a response chunk into JSON and writes it as an SSE data event.
func sendStreamEvent(w http.ResponseWriter, flusher http.Flusher, resp ChatCompletionStreamResponse) {
	data, err := json.Marshal(resp)
	if err != nil {
		slog.Error("Failed to marshal stream event", "error", err)
		return
	}

	// SSE format "data: {json}\n\n"
	w.Write([]byte("data: "))
	w.Write(data)
	w.Write([]byte("\n\n"))
	flusher.Flush()
}
