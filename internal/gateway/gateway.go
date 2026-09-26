package gateway

import (
	"context"
	"fmt"

	"github.com/nightwalker404/llm-gateway/internal/config"
	"github.com/nightwalker404/llm-gateway/internal/provider"
)

type Gateway struct {
	cfg       *config.Config
	providers map[string]provider.Provider
}

func New(cfg *config.Config, providers map[string]provider.Provider) *Gateway {
	return &Gateway{
		cfg:       cfg,
		providers: providers,
	}
}

func (g *Gateway) getProvider(name string) (provider.Provider, error) {
	if name == "" {
		name = g.cfg.DefaultProvider
	}

	p, ok := g.providers[name]
	if !ok {
		return nil, fmt.Errorf("unknown provider: %s", name)
	}
	return p, nil
}

func (g *Gateway) ListModels() map[string][]string {
	result := make(map[string][]string)
	for name, p := range g.providers {
		result[name] = p.AllowedModels()
	}
	return result
}

func (g *Gateway) Chat(ctx context.Context, providerName string, req provider.ChatRequest) (*provider.ChatResponse, error) {
	p, err := g.getProvider(providerName)
	if err != nil {
		return nil, err
	}

	if !p.IsModelAllowed(req.Model) {
		return nil, fmt.Errorf("model %q is not allowed for provider %s", req.Model, p.Name())
	}

	return p.Chat(ctx, req)

}
