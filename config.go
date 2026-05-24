package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	ListenAddr      string `json:"listen_addr"`
	LlamaURL        string `json:"llama_url"`
	QdrantURL       string `json:"qdrant_url"`
	QdrantCollection string `json:"qdrant_collection"`
	EmbeddingSize   int    `json:"embedding_size"`
	DefaultTopK     int    `json:"default_top_k"`
}

func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config: %w", err)
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}

	if cfg.ListenAddr == "" {
		cfg.ListenAddr = ":8080"
	}
	if cfg.DefaultTopK <= 0 {
		cfg.DefaultTopK = 5
	}
	if cfg.EmbeddingSize <= 0 {
		cfg.EmbeddingSize = 1024
	}

	return &cfg, nil
}