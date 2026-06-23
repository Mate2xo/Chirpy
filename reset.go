package main

import (
	"context"
	"log"
	"net/http"
	"os"
)

func (cfg *apiConfig) reset(w http.ResponseWriter, _ *http.Request) {
	if os.Getenv("PLATFORM") != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	cfg.fileserverHits.Store(0)
	w.WriteHeader(200)
	w.Write([]byte("Hits reset to 0"))
	println("Deleting all users...")
	err := cfg.dbQueries.DeleteUsers(context.Background())
	if err != nil {
		log.Printf("Error deleting users: %v", err)
	}
}
