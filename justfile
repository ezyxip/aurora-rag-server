default:
  aseek-server --help

# Build the server binary
build:
  go build -o aseek-server .

# Run Qdrant via docker-compose
up-db:
  docker compose up -d qdrant

# Stop all Docker services
down-db:
  docker compose down

embed_model := "nomic-ai/nomic-embed-text-v1.5-GGUF"
reranker_model := "BAAI/bge-reranker-v2-m3-GGUF"
llama_port_embed := "8081"
llama_port_rerank := "8082"
llama_n_gpu_layers := "99"

# Start llama-server for embeddings
up-embed:
  llama-server \
    --host 0.0.0.0 \
    --port {{llama_port_embed}} \
    --model {{embed_model}} \
    --gpu-layers {{llama_n_gpu_layers}} \
    --embedding \
    --pooling mean

# Start llama-server for reranking
up-rerank:
  llama-server \
    --host 0.0.0.0 \
    --port {{llama_port_rerank}} \
    --model {{reranker_model}} \
    --gpu-layers {{llama_n_gpu_layers}} \
    --reranking

# Start all backend services
up-all: up-db up-embed up-rerank

# Run the server (assumes Qdrant and llama-server are already running)
run:
  ./aseek-server

# Build + run
start: build run

# Quick sanity check: embed -> search
test-query query="test":
  curl -s -X POST http://localhost:8080/api/search \
    -H 'Content-Type: application/json' \
    -d '{"query":"{{query}}","top_k":3}' | jq .

# Show all commands
list:
  @just --list
