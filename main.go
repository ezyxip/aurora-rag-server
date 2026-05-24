package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	configPath := flag.String("config", "./config.json", "path to config file")
	flag.Parse()

	cfg, err := LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	embedder := NewEmbedder(cfg.LlamaURL)
	searcher := NewSearcher(cfg.QdrantURL, cfg.QdrantCollection)
	handler := NewServer(embedder, searcher, cfg.DefaultTopK)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/search", handler.ServeHTTP)

	log.Printf("starting server on %s", cfg.ListenAddr)
	if err := http.ListenAndServe(cfg.ListenAddr, mux); err != nil {
		log.Fatalf("server: %v", err)
	}
}