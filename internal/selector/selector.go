package selector

import (
	"context"

	"github.com/nightwalker404/llm-gateway/internal/provider"
)

type Selection struct {
	Provider provider.Provider `json:"provider"`
	Model    string            `json:"model"`
}

type Selector interface {
	Select(ctx context.Context, req provider.ChatRequest) (*Selection, error)
}
