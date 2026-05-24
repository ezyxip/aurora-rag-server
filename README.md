# aseek-server

Go-прослойка для векторного поиска документов. Принимает текстовый запрос, получает эмбеддинг через llama-server (`/v1/embeddings`), ищет ближайшие вектора в Qdrant, возвращает результаты.

## Архитектура

```
POST /api/search  →  aseek-server  →  llama-server (embeddings)
                                    →  Qdrant (vector search)
                         (reranking — отдельный сервис на llama-server)
```

## Конфигурация

`config.json` (по умолч. `./config.json`, переопределяется флагом `-config`):

```json
{
  "listen_addr": ":8080",
  "llama_url": "http://localhost:8081",
  "qdrant_url": "http://localhost:6333",
  "qdrant_collection": "documents",
  "embedding_size": 1024,
  "default_top_k": 5
}
```

## Быстрый старт

```bash
# 1. Запустить Qdrant
just up-db

# 2. Запустить llama-server для эмбеддингов
just up-embed

# 3. Собрать и запустить сервер
just start

# 4. Проверить
just test-query "что такое семафор"
```

Все сервисы разом:

```bash
just up-all     # Qdrant + 2× llama-server (embeddings + reranking)
just start      # build + run сервера
```

## API

`POST /api/search`

```json
{"query":"<текст>","top_k":5}
```

Ответ:

```json
{
  "results": [
    {"title":"...", "content":"...", "source":"...", "score":0.95}
  ]
}
```

## Justfile

| Команда | Описание |
|---|---|
| `just build` | Собрать бинарник |
| `just up-db` | Запустить Qdrant |
| `just up-embed` | Запустить llama-server (embeddings) |
| `just up-rerank` | Запустить llama-server (reranking) |
| `just up-all` | Запустить всё |
| `just down-db` | Остановить Docker-сервисы |
| `just run` | Запустить сервер |
| `just start` | build + run |
| `just test-query "..."` | curl-проверка |

## Зависимости

- Go 1.21+
- [just](https://just.systems/) (make-альтернатива)
- Docker (для Qdrant)
- [llama.cpp](https://github.com/ggml-org/llama.cpp) (бинарник `llama-server` в `$PATH`)