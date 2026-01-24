package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type jsonResponse struct {
	Message string `json:"message"`
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", logging(http.HandlerFunc(rootHandler)))
	mux.Handle("/healthz", logging(http.HandlerFunc(healthHandler)))

	port := getenv("PORT", "8080")
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("starting server on port %s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, jsonResponse{Message: "hello from sample-go-api"})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, jsonResponse{Message: "ok"})
}

func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
