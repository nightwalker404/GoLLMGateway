# GoLLMGateway

Simple and maintainable **Go** gateway for local LLM backends (**Ollama** + **vLLM**).

It provides a unified HTTP API so your applications only need to talk to one endpoint.  
When no provider is specified, a small model acts as a **selector** and decides which provider + model should handle the request.

---

## Features

- Unified chat API for Ollama and vLLM
- **Intelligent provider/model selection** using a tiny LLM
- Explicit provider selection still supported (`"provider": "ollama"` or `"vllm"`)
- Model allow-listing per provider (strict – empty list is not allowed)
- Clean configuration via environment variables / `.env`
- Graceful shutdown
- Health check + model listing endpoints
- Docker & Docker Compose support for the backends

---

## Quick Start

### 1. Clone & configure

```bash
git clone https://github.com/nightwalker404/GoLLMGateway.git
cd GoLLMGateway

cp .env.example .env
# edit .env if needed

---
```

---
                        ┌─────────────────────┐
                        │   Client / App      │
                        └─────────┬───────────┘
                                  │
                                  ▼
                        ┌─────────────────────┐
                        │   GoLLMGateway      │
                        │  (HTTP :8080)       │
                        └─────────┬───────────┘
                                  │
                    ┌─────────────┴─────────────┐
                    │                           │
          Explicit provider?              No provider
                    │                           │
                    ▼                           ▼
          Use given provider          Tiny Model Selector
          + given model               (decides provider + model)
                    │                           │
                    └─────────────┬─────────────┘
                                  │
                    ┌─────────────┴─────────────┐
                    │                           │
                    ▼                           ▼
            ┌──────────────┐           ┌──────────────┐
            │    Ollama    │           │     vLLM     │
            └──────────────┘           └──────────────┘