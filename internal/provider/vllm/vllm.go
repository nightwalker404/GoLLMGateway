package vllm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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

func New(cfg config.VLLMConfig) *Client {
	return &Client{
		baseURL: cfg.BaseURL,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
		models: cfg.Models,
	}
}

func (c *Client) Name() string { return "vllm" }

func (c *Client) Chat(ctx context.Context, req provider.ChatRequest) (*provider.ChatResponse, error) {
	body := map[string]any{
		"model":    req.Model,
		"messages": req.Messages,
		"stream":   false,
	}

	if req.Temperature > 0 {
		body["temperature"] = req.Temperature
	}
	if req.MaxTokens > 0 {
		body["max_tokens"] = req.MaxTokens
	}

	payload, err := json.Marshal(body)

	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/v1/chat/completions",
		bytes.NewReader(payload),
	)

	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("vllm request faild: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("vllm returned %d: %s", resp.StatusCode, string(b))
	}

	var vllmResp struct {
		ID      string `json:"id"`
		Model   string `json:"model"`
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		} `json:"usage"`
	}

	if len(vllmResp.Choices) == 0 {
		return nil, fmt.Errorf("vllm returned no choices")
	}

	return &provider.ChatResponse{
		ID:      vllmResp.ID,
		Model:   vllmResp.Model,
		Content: vllmResp.Choices[0].Message.Content,
		Usage: provider.Usage{
			PromptTokens:     vllmResp.Usage.PromptTokens,
			CompletionTokens: vllmResp.Usage.CompletionTokens,
			TotalTokens:      vllmResp.Usage.TotalTokens,
		},
	}, nil
}
