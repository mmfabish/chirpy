package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

const filepathRoot = "."

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, req *http.Request) {
	responseMessage := fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(responseMessage))
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)

	w.WriteHeader(http.StatusOK)
}

func healthCheckHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	mux := http.NewServeMux()
	cfg := apiConfig{}
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))

	mux.Handle("/app/", cfg.middlewareMetricsInc(handler))
	mux.HandleFunc("GET /healthz", healthCheckHandler)
	mux.HandleFunc("GET /metrics", cfg.metricsHandler)
	mux.HandleFunc("POST /reset", cfg.resetHandler)

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	server.ListenAndServe()
}
