package gateway

import (
	"context"
	"fmt"

	"github.com/nightwalker404/llm-gateway/internal/config"
	"github.com/nightwalker404/llm-gateway/internal/provider"
	"github.com/nightwalker404/llm-gateway/internal/selector"
)

type Gateway struct {
	cfg       *config.Config
	providers map[string]provider.Provider
	selector  selector.Selector
}

func New(cfg *config.Config, providers map[string]provider.Provider, sel selector.Selector) *Gateway {
	return &Gateway{
		cfg:       cfg,
		providers: providers,
		selector:  sel,
	}
}

func (g *Gateway) selectProvider(name string, ctx context.Context, req provider.ChatRequest) (*selector.Selection, error) {
	if name != "" {
		p, ok := g.providers[name]
		if !ok {
			return nil, fmt.Errorf("unknown provider: %s", name)
		}
		return &selector.Selection{
			Provider: p,
			Model:    req.Model,
		}, nil
	}

	sel, err := g.selector.Select(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to select provider: %w", err)
	}
	return sel, nil
}

func (g *Gateway) ListModels() map[string][]string {
	result := make(map[string][]string)
	for name, p := range g.providers {
		result[name] = p.AllowedModels()
	}
	return result
}

func (g *Gateway) Chat(ctx context.Context, providerName string, req provider.ChatRequest) (*provider.ChatResponse, error) {
	sel, err := g.selectProvider(providerName, ctx, req)
	if err != nil {
		return nil, err
	}

	req.Model = sel.Model

	if !sel.Provider.IsModelAllowed(req.Model) {
		return nil, fmt.Errorf("model %q is not allowed for provider %s", req.Model, sel.Provider.Name())
	}

	return sel.Provider.Chat(ctx, req)

}
