package config

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type OllamaConfig struct {
	BaseURL string        `env:"OLLAMA_BASE_URL" envDefault:"http://localhost:11434"`
	Timeout time.Duration `env:"OLLAMA_TIMEOUT" envDefault:"60s"`
	Models  []string      `env:"OLLAMA_MODELS" envSeparator:","`
}

type VLLMConfig struct {
	BaseURL string        `env:"VLLM_BASE_URL" envDefault:"http://localhost:8000"`
	APIKey  string        `env:"VLLM_API_KEY"`
	Timeout time.Duration `env:"VLLM_TIMEOUT" envDefault:"60s"`
	Models  []string      `env:"VLLM_MODELS" envSeparator:","`
}

type SelectorConfig struct {
	Provider         string        `env:"SELECTOR_PROVIDER" envDefault:"vllm"`
	Model            string        `env:"SELECTOR_MODEL" envDefault:"Qwen/Qwen2.5-0.5B-Instruct"`
	Timeout          time.Duration `env:"SELECTOR_TIMEOUT" envDefault:"5s"`
	FallbackProvider string        `env:"SELECTOR_FALLBACK_PROVIDER" envDefault:"ollama"`
	FallbackModel    string        `env:"SELECTOR_FALLBACK_MODEL" envDefault:"llama3.2:3b"`
}

type Config struct {
	ServerAddr      string `env:"SERVER_ADDR" envDefault:":8080"`
	DefaultProvider string `env:"DEFAULT_PROVIDER" envDefault:"ollama"`

	Ollama   OllamaConfig
	VLLM     VLLMConfig
	Selector SelectorConfig
}

var (
	cfg  Config
	once sync.Once
	err  error
)

func cleanModels(models []string) []string {
	result := make([]string, 0, len(models))
	for _, m := range models {
		m = strings.TrimSpace(m)
		if m != "" {
			result = append(result, m)
		}
	}
	return result
}

func Load() (Config, error) {
	once.Do(func() {
		_ = godotenv.Load()

		cfg = Config{}
		if err = env.Parse(&cfg); err != nil {
			return
		}

		cfg.Ollama.Models = cleanModels(cfg.Ollama.Models)
		cfg.VLLM.Models = cleanModels(cfg.VLLM.Models)
		cfg.Selector.Provider = strings.ToLower(cfg.Selector.Provider)
		cfg.Selector.FallbackProvider = strings.ToLower(cfg.Selector.FallbackProvider)

		if cfg.DefaultProvider != "ollama" && cfg.DefaultProvider != "vllm" {
			err = fmt.Errorf("DEFAULT_PROVIDER must be 'ollama' or 'vllm', got %q", cfg.DefaultProvider)
		}
	})

	return cfg, err
}

func Get() Config {
	return cfg
}
