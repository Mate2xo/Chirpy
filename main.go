package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const port = "8080"
	cfg := &apiConfig{}
	mux := initMux(cfg)
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}

func initMux(cfg *apiConfig) *http.ServeMux {
	mux := http.NewServeMux()
	registerRoutes(mux, cfg)
	return mux
}

func registerRoutes(mux *http.ServeMux, cfg *apiConfig) {
	mux.Handle("/app/", cfg.middlewareMetricsInc(fileRoot))
	mux.HandleFunc("GET /api/healthz", healthz)
	mux.HandleFunc("GET /api/metrics", cfg.metrics)
	mux.HandleFunc("POST /api/reset", cfg.reset)
}

var fileRoot = http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
