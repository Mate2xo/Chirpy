package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func validateChirp(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	type parameters struct {
		Body string `json:"body"`
	}
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithErr(err, http.StatusInternalServerError, w)
		return
	}

	if length := len(params.Body); length > 140 {
		err = fmt.Errorf("error: body can be of a maximum length of 140 characters (currently %d)", length)
		respondWithErr(err, http.StatusBadRequest, w)
		return
	}

	respondWithJSON(validated{Valid: true}, http.StatusOK, w)
}

type validated struct {
	Valid bool `json:"valid"`
}
