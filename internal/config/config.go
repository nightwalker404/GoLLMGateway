package config

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type OllamaConfig struct {
	BaseUrl string        `env:"OLLAMA_BASE_URL" envDefault:"http://localhost:11434"`
	Timeout time.Duration `env:"OLLAMA_TIMEOUT" envDefault:"60s"`
}

type VllmConfig struct {
	BaseUrl string        `env:"VLLM_BASE_URL" envDefault:"http://localhost:8000"`
	Timeout time.Duration `env:"VLLM_TIMEOUT" envDefault:"60s"`
}

type config struct {
	ServerAddr      string `env:"SERVER_ADDR" envDefault:"8080"`
	DefaultProvider string `env:"DEFAULT_PROVIDER" envDefault:"ollama"`

	Ollama OllamaConfig
	Vllm   VllmConfig
}

func Load() (config, error) {
	cfg := config{}
	if err := env.FieldParams(&cfg); err != nil {
		return config{}, err
	}
	return cfg, nil
}
