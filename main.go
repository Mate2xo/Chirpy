package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}

func (cfg *apiConfig) metrics(w http.ResponseWriter, _ *http.Request) {
	fmt.Printf("hits are %v\n", cfg.fileserverHits.Load())
	_, err := fmt.Fprintf(w, "Hits: %v", cfg.fileserverHits.Load())
	if err != nil {
		log.Fatal(err)
	}
}

func (cfg *apiConfig) reset(_ http.ResponseWriter, _ *http.Request) {
	cfg.fileserverHits.Store(0)
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
	mux.HandleFunc("/healthz", healthz)
	mux.HandleFunc("/metrics", cfg.metrics)
	mux.HandleFunc("/reset", cfg.reset)
}

var fileRoot = http.StripPrefix("/app/", http.FileServer(http.Dir(".")))

func healthz(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(http.StatusText(200)))
}
