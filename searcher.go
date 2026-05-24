package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type qdrantSearchRequest struct {
	Vector    []float32 `json:"vector"`
	Limit     int       `json:"limit"`
	WithPayload bool    `json:"with_payload"`
}

type qdrantSearchResponse struct {
	Result []qdrantScoredPoint `json:"result"`
}

type qdrantScoredPoint struct {
	ID      any                    `json:"id"`
	Score   float64                `json:"score"`
	Version int                    `json:"version"`
	Payload map[string]any         `json:"payload,omitempty"`
}

type searchDocument struct {
	Index   int     `json:"index"`
	Title   string  `json:"title"`
	Content string  `json:"content"`
	Source  string  `json:"source"`
	Score   float64 `json:"score"`
}

type Searcher struct {
	baseURL    string
	collection string
	client     *http.Client
}

func NewSearcher(baseURL, collection string) *Searcher {
	return &Searcher{
		baseURL:    baseURL,
		collection: collection,
		client:     &http.Client{},
	}
}

func (s *Searcher) Search(ctx context.Context, vector []float32, limit int) ([]searchDocument, error) {
	reqBody := qdrantSearchRequest{
		Vector:      vector,
		Limit:       limit,
		WithPayload: true,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal qdrant search: %w", err)
	}

	url := fmt.Sprintf("%s/collections/%s/points/search", s.baseURL, s.collection)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create qdrant request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("qdrant search: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("qdrant returned %d: %s", resp.StatusCode, string(raw))
	}

	var qr qdrantSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&qr); err != nil {
		return nil, fmt.Errorf("decode qdrant response: %w", err)
	}

	docs := make([]searchDocument, 0, len(qr.Result))
	for i, pt := range qr.Result {
		d := searchDocument{
			Index: i,
			Score: pt.Score,
		}
		if pt.Payload != nil {
			if v, ok := pt.Payload["title"]; ok {
				d.Title, _ = v.(string)
			}
			if v, ok := pt.Payload["content"]; ok {
				d.Content, _ = v.(string)
			}
			if v, ok := pt.Payload["source"]; ok {
				d.Source, _ = v.(string)
			}
		}
		docs = append(docs, d)
	}

	return docs, nil
}