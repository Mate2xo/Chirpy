package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/Mate2xo/Chirpy/internal/auth"
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

type chirpsResponse struct {
	elements []chirpResponse
}

func (cfg *apiConfig) indexChirps(w http.ResponseWriter, req *http.Request) {
	chirps, err := cfg.dbQueries.AllChirps(req.Context())
	if err != nil {
		respondWithErr(err, http.StatusInternalServerError, w)
	}

	collection := chirpsResponse{}
	for _, chirp := range chirps {
		collection.elements = append(
			collection.elements, chirpResponse{
				ID:        chirp.ID,
				Body:      chirp.Body,
				UserID:    chirp.UserID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
			},
		)
	}

	respondWithJSON(collection.elements, http.StatusOK, w)
}

func (cfg *apiConfig) showChirp(w http.ResponseWriter, req *http.Request) {
	id, err := uuid.Parse(req.PathValue("id"))
	if err != nil {
		respondWithErr(err, http.StatusInternalServerError, w)
		return
	}

	chirp, err := cfg.dbQueries.GetChirp(req.Context(), id)
	if err != nil {
		respondWithErr(err, http.StatusNotFound, w)
		return
	}

	member := chirpResponse{
		ID:        chirp.ID,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
	}

	respondWithJSON(member, http.StatusOK, w)
}

func (cfg *apiConfig) createChirp(w http.ResponseWriter, req *http.Request) {
	params := chirpParams{}
	err := json.NewDecoder(req.Body).Decode(&params)
	if err != nil {
		log.Printf("Error: could not decode JSON: %v\n", err)
		respondWithErr(err, http.StatusInternalServerError, w)
	}

	jwt, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithErr(err, http.StatusUnauthorized, w)
	}
	userID, err := auth.ValidateJWT(jwt, cfg.jwtSecret)
	if err != nil {
		respondWithErr(err, http.StatusUnauthorized, w)
		return
	}

	err = validateChirp(params)
	if err != nil {
		respondWithErr(err, http.StatusBadRequest, w)
	}

	cleaned, err := removeBadWords(params)
	if err != nil {
		respondWithErr(err, http.StatusInternalServerError, w)
	}

	createParams := database.CreateChirpParams{Body: cleaned, UserID: userID}
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

func (cfg *apiConfig) deleteChirp(w http.ResponseWriter, req *http.Request) {
	id, err := uuid.Parse(req.PathValue("id"))
	if err != nil {
		respondWithErr(err, http.StatusBadRequest, w)
		return
	}

	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithErr(err, http.StatusUnauthorized, w)
		return
	}
	currentUserID, err := auth.ValidateJWT(bearerToken, cfg.jwtSecret)
	if err != nil {
		respondWithErr(err, http.StatusUnauthorized, w)
		return
	}

	chirp, err := cfg.dbQueries.GetChirp(req.Context(), id)
	if err != nil {
		respondWithErr(err, http.StatusNotFound, w)
		return
	}
	if currentUserID != chirp.UserID {
		respondWithErr(err, http.StatusForbidden, w)
		return
	}

	err = cfg.dbQueries.DeleteChrip(req.Context(), id)
	if err != nil {
		respondWithErr(err, http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
