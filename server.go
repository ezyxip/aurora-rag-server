package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type searchRequest struct {
	Query string `json:"query"`
	TopK  int    `json:"top_k"`
}

type searchResponse struct {
	Results []searchDocument `json:"results"`
}

type Server struct {
	embedder *Embedder
	searcher *Searcher
	topK     int
}

func NewServer(embedder *Embedder, searcher *Searcher, topK int) *Server {
	return &Server{
		embedder: embedder,
		searcher: searcher,
		topK:     topK,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req searchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Query == "" {
		http.Error(w, "query is required", http.StatusBadRequest)
		return
	}

	topK := s.topK
	if req.TopK > 0 {
		topK = req.TopK
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	vec, err := s.embedder.Embed(ctx, req.Query)
	if err != nil {
		log.Printf("embedding error: %v", err)
		http.Error(w, "embedding failed", http.StatusInternalServerError)
		return
	}

	docs, err := s.searcher.Search(ctx, vec, topK)
	if err != nil {
		log.Printf("search error: %v", err)
		http.Error(w, "search failed", http.StatusInternalServerError)
		return
	}

	resp := searchResponse{Results: docs}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}