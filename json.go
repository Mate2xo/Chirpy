package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}

func respondWithErr(err error, code int, w http.ResponseWriter) {
	payload := errorResponse{Error: fmt.Sprintf("%v", err)}
	respondWithJSON(payload, code, w)
}

func respondWithJSON(payload any, code int, w http.ResponseWriter) {
	w.Header().Set("Content-Tyte", "application/json")

	body, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error marshalling JSON: %s", err)
		return
	}

	w.WriteHeader(code)
	w.Write([]byte(body))
}
