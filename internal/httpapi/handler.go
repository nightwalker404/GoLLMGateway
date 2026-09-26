package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/nightwalker404/llm-gateway/internal/gateway"
	"github.com/nightwalker404/llm-gateway/internal/provider"
)

type Handler struct {
	gw *gateway.Gateway
}

type chatRequestBody struct {
	Provider    string             `json:"provider,omitempty"`
	Model       string             `json:"model"`
	Messages    []provider.Message `json:"messages"`
	Temperature float64            `json:"temperature,omitempty"`
	MaxTokens   int                `json:"max_tokens,omitempty"`
}

func NewHandler(gw *gateway.Gateway) *Handler {
	return &Handler{gw: gw}
}

func (h *Handler) Chat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body chatRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json: "+err.Error(), http.StatusBadRequest)
		return
	}

	if body.Model == "" || len(body.Messages) == 0 {
		http.Error(w, "model and messages are required", http.StatusBadRequest)
		return
	}

	req := provider.ChatRequest{
		Model:       body.Model,
		Messages:    body.Messages,
		Temperature: body.Temperature,
		MaxTokens:   body.MaxTokens,
	}

	resp, err := h.gw.Chat(r.Context(), body.Provider, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
