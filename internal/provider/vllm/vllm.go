package vllm

import (
	"context"
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

}
