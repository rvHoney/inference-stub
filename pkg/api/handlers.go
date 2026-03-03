// Package api contains HTTP handlers for the Inference-Stub endpoints
package api

import (
	"encoding/json"
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
			initialResp := ChatCompletionStreamResponse{
				ID:                "chatcmpl-mock-stream123",
				Object:            "chat.completion.chunk",
				Created:           time.Now().Unix(),
				Model:             req.Model,
				SystemFingerprint: "fp_mock_stub",
				Choices: []StreamChoice{
					{
						Index: 0,
						Delta: ChatDelta{
							Role: "assistant",
						},
						FinishReason: nil,
					},
				},
			}
			sendStreamEvent(w, flusher, initialResp)

			// Generate N tokens of Lorem Ipsum separated by TPOT
			tokens := []string{"Lorem", " ", "ipsum", ",", " ", "dolor", " ", "sit", " ", "amet", ",", " ", "consectetur", " ", "adipiscing", " ", "elit", "."}

			for _, token := range tokens {
				time.Sleep(cfg.TPOT)

				chunkResp := ChatCompletionStreamResponse{
					ID:                "chatcmpl-mock-stream123",
					Object:            "chat.completion.chunk",
					Created:           time.Now().Unix(),
					Model:             req.Model,
					SystemFingerprint: "fp_mock_stub",
					Choices: []StreamChoice{
						{
							Index: 0,
							Delta: ChatDelta{
								Content: token,
							},
							FinishReason: nil,
						},
					},
				}
				sendStreamEvent(w, flusher, chunkResp)
			}

			// Send final chunk w/ finish_reason
			time.Sleep(cfg.TPOT)
			stopReason := "stop"
			finalResp := ChatCompletionStreamResponse{
				ID:                "chatcmpl-mock-stream123",
				Object:            "chat.completion.chunk",
				Created:           time.Now().Unix(),
				Model:             req.Model,
				SystemFingerprint: "fp_mock_stub",
				Choices: []StreamChoice{
					{
						Index:        0,
						Delta:        ChatDelta{},
						FinishReason: &stopReason,
					},
				},
			}
			sendStreamEvent(w, flusher, finalResp)

			// Send termination signal
			w.Write([]byte("data: [DONE]\n\n"))
			flusher.Flush()
		} else {
			tokenCount := 10
			tpotDelay := time.Duration(tokenCount) * cfg.TPOT
			time.Sleep(tpotDelay)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

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
							Content: "Hello! How can I assist you today?",
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
