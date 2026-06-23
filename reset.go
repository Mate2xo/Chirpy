package main

import (
	"context"
	"log"
	"net/http"
)

func (cfg *apiConfig) reset(w http.ResponseWriter, _ *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	cfg.fileserverHits.Store(0)
	w.WriteHeader(200)
	w.Write([]byte("Hits reset to 0"))
	println("Deleting all users...")
	err := cfg.dbQueries.Reset(context.Background())
	if err != nil {
		log.Printf("Error resetting DB: %v", err)
	}
}
