package httpapi

import (
	"context"
	"net/http"
	"time"

	"github.com/nightwalker404/llm-gateway/internal/config"
	"github.com/nightwalker404/llm-gateway/internal/gateway"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) ListenAndServe() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}

func NewServer(cfg *config.Config, gw *gateway.Gateway) *Server {
	mux := http.NewServeMux()
	h := NewHandler(gw)

	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/v1/chat/completions", h.Chat)
	mux.HandleFunc("/chat", h.Chat)

	return &Server{
		httpServer: &http.Server{
			Addr:         cfg.ServerAddr,
			Handler:      mux,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 120 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}
