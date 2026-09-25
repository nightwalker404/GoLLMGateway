package main

import (
	"log"

	"github.com/nightwalker404/llm-gateway/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("config: %v", err)
	}
}
