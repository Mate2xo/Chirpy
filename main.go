package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	mux := initMux()
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}

func initMux() *http.ServeMux {
	mux := http.NewServeMux()
	registerRoutes(mux)
	return mux
}

func registerRoutes(mux *http.ServeMux) {
	mux.Handle("/{$}", root())
}

func root() http.Handler {
	return http.FileServer(http.Dir("."))
}
