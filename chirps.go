package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/Mate2xo/Chirpy/internal/database"
	"github.com/google/uuid"
)

type validationError struct {
	message string
}

func (err validationError) Error() string {
	return err.message
}

type chirpParams struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type chirpResponse struct {
	ID        uuid.UUID `json:"id"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (cfg *apiConfig) postChirp(w http.ResponseWriter, req *http.Request) {
	params := chirpParams{}
	err := json.NewDecoder(req.Body).Decode(&params)
	if err != nil {
		log.Printf("Error: could not decode JSON: %v\n", err)
		respondWithErr(err, http.StatusInternalServerError, w)
	}
	fmt.Printf("Params: %+v\n", params)

	err = validateChirp(params)
	if err != nil {
		respondWithErr(err, http.StatusBadRequest, w)
	}

	cleaned, err := removeBadWords(params)
	if err != nil {
		respondWithErr(err, http.StatusInternalServerError, w)
	}

	createParams := database.CreateChirpParams{Body: cleaned, UserID: params.UserID}
	chirp, err := cfg.dbQueries.CreateChirp(req.Context(), createParams)
	if err != nil {
		respondWithErr(err, http.StatusInternalServerError, w)
	}

	response := chirpResponse{
		ID:        chirp.ID,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
	}
	respondWithJSON(response, http.StatusCreated, w)
}

func validateChirp(params chirpParams) error {
	const maxChirpLength = 140
	if length := len(params.Body); length > maxChirpLength {
		msg := fmt.Sprintf("error: body can be of a maximum length of 140 characters (currently %d)", length)
		return validationError{message: msg}
	}

	return nil
}

func removeBadWords(params chirpParams) (string, error) {
	re, err := regexp.Compile(`(?i)(kerfuffle|sharbert|fornax)[^!]`)
	if err != nil {
		return "", fmt.Errorf("error compiling regex: %w", err)
	}

	cleaned := re.ReplaceAll([]byte(params.Body), []byte("**** "))
	return string(cleaned), nil
}
