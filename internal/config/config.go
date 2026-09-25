package config

import (
	"fmt"
	"sync"
	"time"

	"github.com/caarlos0/env/v11"
)

type OllamaConfig struct {
	BaseURL string        `env:"OLLAMA_BASE_URL" envDefault:"http://localhost:11434"`
	Timeout time.Duration `env:"OLLAMA_TIMEOUT" envDefault:"60s"`
}

type VLLMConfig struct {
	BaseURL string        `env:"VLLM_BASE_URL" envDefault:"http://localhost:8000"`
	APIKey  string        `env:"VLLM_API_KEY"`
	Timeout time.Duration `env:"VLLM_TIMEOUT" envDefault:"60s"`
}

type Config struct {
	ServerAddr      string `env:"SERVER_ADDR" envDefault:":8080"`
	DefaultProvider string `env:"DEFAULT_PROVIDER" envDefault:"ollama"`

	Ollama OllamaConfig
	VLLM   VLLMConfig
}

var (
	cfg  Config
	once sync.Once
	err  error
)

func Load() (Config, error) {
	once.Do(func() {
		cfg = Config{}
		if err = env.Parse(&cfg); err != nil {
			return
		}

		if cfg.DefaultProvider != "ollama" && cfg.DefaultProvider != "vllm" {
			err = fmt.Errorf("DEFAULT_PROVIDER must be 'ollama' or 'vllm', got %q", cfg.DefaultProvider)
		}
	})

	return cfg, err
}

func Get() Config {
	return cfg
}
