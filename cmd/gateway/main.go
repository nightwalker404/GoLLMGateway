package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/nightwalker404/llm-gateway/internal/config"
	"github.com/nightwalker404/llm-gateway/internal/gateway"
	"github.com/nightwalker404/llm-gateway/internal/httpapi"
	"github.com/nightwalker404/llm-gateway/internal/provider"
	"github.com/nightwalker404/llm-gateway/internal/provider/ollama"
	"github.com/nightwalker404/llm-gateway/internal/provider/vllm"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	providers := map[string]provider.Provider{
		"ollama": ollama.New(cfg.Ollama),
		"vllm":   vllm.New(cfg.VLLM),
	}

	gw := gateway.New(&cfg, providers)

	server := httpapi.NewServer(&cfg, gw)

	go func() {
		log.Printf("LLM Gateway started on %s | default provider: %s", cfg.ServerAddr, cfg.DefaultProvider)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	if err := server.Shutdown(); err != nil {
		log.Printf("shutdown error: %v", err)
	}
	log.Println("server stopped")
}
