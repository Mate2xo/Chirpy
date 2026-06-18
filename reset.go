package main

import "net/http"

func (cfg *apiConfig) reset(w http.ResponseWriter, _ *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(200)
	w.Write([]byte("Hits reset to 0"))
}
