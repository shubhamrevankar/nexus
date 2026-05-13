package httpserver

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHealthHandler(t *testing.T) {
	server := httptest.NewServer(NewRouter(slog.New(slog.NewTextHandler(io.Discard, nil)), nil, nil, nil, nil, nil, "", time.Hour))
	defer server.Close()

	response, err := http.Get(server.URL + "/healthz")
	if err != nil {
		t.Fatalf("expected health endpoint to respond: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	var payload healthResponse
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatalf("expected valid json response: %v", err)
	}

	if payload.Service != "api" || payload.Status != "ok" {
		t.Fatalf("unexpected health response: %+v", payload)
	}
}
