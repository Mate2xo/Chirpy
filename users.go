package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Mate2xo/Chirpy/internal/auth"
	"github.com/Mate2xo/Chirpy/internal/database"
	"github.com/google/uuid"
)

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}
type userParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) postUser(w http.ResponseWriter, req *http.Request) {
	params := userParams{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error: could not decode JSON: %v\n", err)
		respondWithErr(err, http.StatusBadRequest, w)
		return
	}

	password, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithErr(err, http.StatusInternalServerError, w)
		return
	}

	user, err := cfg.dbQueries.CreateUser(context.Background(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: password,
	})
	if err != nil {
		log.Printf("Error: could not create user: %v\n", err)
		respondWithErr(err, http.StatusInternalServerError, w)
		return
	}

	userResponse := UserResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		Email:     user.Email,
	}

	fmt.Printf("Created user with email %s\n", user.Email)
	respondWithJSON(userResponse, http.StatusCreated, w)
}

func (cfg *apiConfig) loginUser(w http.ResponseWriter, req *http.Request) {
	params := userParams{}
	err := json.NewDecoder(req.Body).Decode(&params)
	if err != nil {
		respondWithErr(err, http.StatusInternalServerError, w)
		return
	}

	user, err := cfg.dbQueries.GetUserByEmail(req.Context(), params.Email)
	if err != nil {
		respondWithErr(err, http.StatusNotFound, w)
		return
	}

	if ok, err := auth.CheckPasswordHash(params.Password, user.HashedPassword); !ok || err != nil {
		err = errors.New("incorrect email or password")
		respondWithErr(err, http.StatusUnauthorized, w)
		return
	}

	payload := UserResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	respondWithJSON(payload, http.StatusOK, w)
}
