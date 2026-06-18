package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

type validated struct {
	CleanedBody string `json:"cleaned_body"`
}
type parameters struct {
	Body string `json:"body"`
}

func validateChirp(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
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

	payload, err := removeBadWords(params)
	if err != nil {
		respondWithErr(err, http.StatusInternalServerError, w)
	}

	respondWithJSON(validated{CleanedBody: payload}, http.StatusOK, w)
}

func removeBadWords(params parameters) (string, error) {
	re, err := regexp.Compile(`(?i)(kerfuffle|sharbert|fornax)[^!]`)
	if err != nil {
		return "", fmt.Errorf("error compiling regex: %w", err)
	}

	cleaned := re.ReplaceAll([]byte(params.Body), []byte("**** "))
	return string(cleaned), nil
}
