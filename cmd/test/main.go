package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/nightwalker404/llm-gateway/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("failed to load config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Ollama models: %#v\n", cfg.Ollama.Models)
	fmt.Printf("vLLM models:   %#v\n", cfg.VLLM.Models)

	baseURL := "http://localhost" + cfg.ServerAddr

	fmt.Println("=== Testing LLM Gateway ===")
	fmt.Printf("Gateway URL: %s\n", baseURL)
	fmt.Printf("Default Provider: %s\n\n", cfg.DefaultProvider)

	// Test List Models
	fmt.Println("→ Testing /v1/models...")
	if err := testListModels(baseURL); err != nil {
		fmt.Printf("  ListModels failed: %v\n", err)
	} else {
		fmt.Println("  ListModels OK")
	}

	// Test Ollama
	fmt.Println("\n→ Testing Ollama...")
	if err := testChat(baseURL, "ollama", "llama3.2:3b"); err != nil {
		fmt.Printf("  Ollama failed: %v\n", err)
	} else {
		fmt.Println("  Ollama OK")
	}

	// Test vLLM
	fmt.Println("\n→ Testing vLLM...")
	if err := testChat(baseURL, "vllm", "Qwen/Qwen2.5-0.5B-Instruct"); err != nil {
		fmt.Printf("  vLLM failed: %v\n", err)
	} else {
		fmt.Println("  vLLM OK")
	}
}

func testChat(baseURL, provider, model string) error {
	body := map[string]any{
		"provider": provider,
		"model":    model,
		"messages": []map[string]string{
			{"role": "user", "content": "Say hello in one short sentence."},
		},
		"temperature": 0.7,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 60 * time.Second}

	resp, err := client.Post(baseURL+"/chat", "application/json", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status %d: %s", resp.StatusCode, string(respBody))
	}

	fmt.Printf("  Response: %s\n", string(respBody))
	return nil
}

func testListModels(baseURL string) error {
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(baseURL + "/v1/models")
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}

	fmt.Printf("  Models: %s\n", string(body))
	return nil
}
