package httpserver

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

type healthResponse struct {
	Service string `json:"service"`
	Status  string `json:"status"`
	Time    string `json:"time"`
}

func NewRouter(logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthHandler)

	return requestLogger(logger, mux)
}

func healthHandler(response http.ResponseWriter, request *http.Request) {
	writeJSON(response, http.StatusOK, healthResponse{
		Service: "api",
		Status:  "ok",
		Time:    time.Now().UTC().Format(time.RFC3339),
	})
}

func requestLogger(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		started := time.Now()
		next.ServeHTTP(response, request)
		logger.Info(
			"http request handled",
			slog.String("method", request.Method),
			slog.String("path", request.URL.Path),
			slog.Duration("duration", time.Since(started)),
		)
	})
}

func writeJSON(response http.ResponseWriter, status int, body any) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status)

	if err := json.NewEncoder(response).Encode(body); err != nil {
		http.Error(response, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

