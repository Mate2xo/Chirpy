package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Mate2xo/Chirpy/internal/database"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
}

func main() {
	const port = "8080"
	cfg := &apiConfig{}
	mux := initMux(cfg)
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	db := initDB(cfg)
	defer db.Close()

	log.Printf("Serving on port %s", port)
	log.Fatal(server.ListenAndServe())
}

func initMux(cfg *apiConfig) *http.ServeMux {
	mux := http.NewServeMux()
	registerRoutes(mux, cfg)
	return mux
}

func initDB(cfg *apiConfig) *sql.DB {
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error: could not connect to DB at %s", dbURL)
	}
	dbQueries := database.New(db)
	cfg.dbQueries = dbQueries

	return db
}

func registerRoutes(mux *http.ServeMux, cfg *apiConfig) {
	mux.Handle("/app/", cfg.middlewareMetricsInc(fileRoot))

	mux.HandleFunc("GET /admin/metrics", cfg.metrics)
	mux.HandleFunc("POST /admin/reset", cfg.reset)

	mux.HandleFunc("GET /api/healthz", healthz)
	mux.HandleFunc("POST /api/validate_chirp", validateChirp)
}

var fileRoot = http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
