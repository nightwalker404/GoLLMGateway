package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/nightwalker404/llm-gateway/internal/config"
	"github.com/nightwalker404/llm-gateway/internal/provider"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	models     []string
}

func (c *Client) AllowedModels() []string {
	return c.models
}

func (c *Client) IsModelAllowed(model string) bool {
	if len(c.models) == 0 {
		return true
	}
	for _, m := range c.models {
		if m == model {
			return true
		}
	}
	return false
}

func New(cfg config.OllamaConfig) *Client {
	return &Client{
		baseURL: cfg.BaseURL,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
		models: cfg.Models,
	}
}

func (c *Client) Name() string { return "ollama" }

func (c *Client) Chat(ctx context.Context, req provider.ChatRequest) (*provider.ChatResponse, error) {
	body := map[string]any{
		"model":    req.Model,
		"messages": req.Messages,
		"stream":   false,
	}

	if req.Temperature > 0 {
		body["options"] = map[string]any{
			"temperature": req.Temperature,
		}
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/api/chat",
		bytes.NewReader(payload),
	)

	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)

	if err != nil {
		return nil, fmt.Errorf("ollama request failed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama returned %d: %s", resp.StatusCode, string(b))
	}

	var ollamaResp struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		Model string `json:"model"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, fmt.Errorf("decode ollama response: %w", err)
	}

	return &provider.ChatResponse{
		ID:      fmt.Sprintf("ollama-%d", time.Now().UnixNano()),
		Model:   ollamaResp.Model,
		Content: ollamaResp.Message.Content,
	}, nil

}
