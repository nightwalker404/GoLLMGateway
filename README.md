# GoLLMGateway

Simple, maintainable **Go** gateway for local LLM backends (**Ollama** and **vLLM**).

It provides a unified HTTP API so your applications only need to talk to one endpoint, regardless of which backend is serving the request.

---

## Features

- Unified chat API for Ollama + vLLM
- Provider selection per request (`"provider": "ollama"` or `"vllm"`)
- Model allow-listing per provider
- Clean configuration via environment variables / `.env`
- Graceful shutdown
- Docker & Docker Compose support
- Health check + model listing endpoints

---

## Quick Start

### 1. Clone & configure

```bash
git clone https://github.com/nightwalker404/GoLLMGateway.git
cd GoLLMGateway

cp .env.example .env
# edit .env if needed