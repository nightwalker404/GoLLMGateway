package selector

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nightwalker404/llm-gateway/internal/config"
	"github.com/nightwalker404/llm-gateway/internal/provider"
)

type LLMSelector struct {
	cfg       config.SelectorConfig
	providers map[string]provider.Provider
	client    provider.Provider
}

func New(cfg config.SelectorConfig, providers map[string]provider.Provider) (*LLMSelector, error) {
	p, ok := providers[cfg.Provider]
	if !ok {
		return nil, fmt.Errorf("selector provider %q not found", cfg.Provider)
	}

	if !p.IsModelAllowed(cfg.Model) {
		return nil, fmt.Errorf("selector model %q is not allowed in provider %s", cfg.Model, cfg.Provider)
	}

	for name, prov := range providers {
		if len(prov.AllowedModels()) == 0 {
			return nil, fmt.Errorf("provider %q has no models configured: refusing to start", name)
		}
	}

	return &LLMSelector{
		cfg:       cfg,
		providers: providers,
		client:    p,
	}, nil
}

func messagesBuilder(providers map[string]provider.Provider, reqMessages []provider.Message) []provider.Message {
	var options strings.Builder
	for name, p := range providers {
		models := p.AllowedModels()
		if len(models) == 0 {
			continue
		}
		options.WriteString(fmt.Sprintf("- provider: %s → models: %s\n", name, strings.Join(models, ", ")))
	}

	systemContent := fmt.Sprintf(SystemPrompt, options.String())

	messages := make([]provider.Message, 0, len(reqMessages)+1)
	messages = append(messages, provider.Message{
		Role:    "system",
		Content: systemContent,
	})
	messages = append(messages, reqMessages...)

	return messages
}

func cleanContent(content string) string {
	data := strings.TrimSpace(content)
	data = strings.TrimPrefix(data, "```json")
	data = strings.TrimPrefix(data, "```")
	data = strings.TrimSuffix(data, "```")
	data = strings.TrimSpace(data)
	return data
}

func (s *LLMSelector) Select(ctx context.Context, req provider.ChatRequest) (*Selection, error) {
	messages := messagesBuilder(s.providers, req.Messages)

	selectorReq := provider.ChatRequest{
		Model:       s.cfg.Model,
		Messages:    messages,
		Temperature: 0.0,
		MaxTokens:   60,
	}

	selCtx, cancel := context.WithTimeout(ctx, s.cfg.Timeout)
	defer cancel()

	resp, err := s.client.Chat(selCtx, selectorReq)
	if err != nil {
		return s.fallback(req), fmt.Errorf("selector call failed: %w", err)
	}

	content := cleanContent(resp.Content)

	var raw struct {
		Provider string `json:"provider"`
		Model    string `json:"model"`
	}

	if err := json.Unmarshal([]byte(content), &raw); err != nil {
		return s.fallback(req), fmt.Errorf("selector returned invalid JSON: %w (raw: %s)", err, content)
	}

	p, ok := s.providers[raw.Provider]
	if !ok {
		return s.fallback(req), fmt.Errorf("selector chose unknown provider %q", raw.Provider)
	}

	if !p.IsModelAllowed(raw.Model) {
		return s.fallback(req), fmt.Errorf("selector chose disallowed model %q for provider %s", raw.Model, raw.Provider)
	}

	return &Selection{
		Provider: p,
		Model:    raw.Model,
	}, nil
}

func (s *LLMSelector) fallback(req provider.ChatRequest) *Selection {
	if req.Model != "" {
		for _, p := range s.providers {
			if p.IsModelAllowed(req.Model) {
				return &Selection{
					Provider: p,
					Model:    req.Model,
				}
			}
		}
	}

	p, ok := s.providers[s.cfg.FallbackProvider]
	if ok && p.IsModelAllowed(s.cfg.FallbackModel) {
		return &Selection{
			Provider: p,
			Model:    s.cfg.FallbackModel,
		}
	}

	for _, p := range s.providers {
		models := p.AllowedModels()
		if len(models) > 0 {
			return &Selection{
				Provider: p,
				Model:    models[0],
			}
		}
	}

	return nil
}
