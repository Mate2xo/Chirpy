package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Mate2xo/Chirpy/internal/database"
	"github.com/google/uuid"
)

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) postUser(w http.ResponseWriter, req *http.Request) {
	params := database.CreateUserParams{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error: could not decode JSON: %v\n", err)
		respondWithErr(err, http.StatusBadRequest, w)
		return
	}
	params.ID = uuid.New()
	fmt.Printf("-- Params: %+v\n", params)

	user, err := cfg.dbQueries.CreateUser(context.Background(), params)
	if err != nil {
		log.Printf("Error: could not create user: %v\n", err)
		respondWithErr(err, http.StatusInternalServerError, w)
		return
	}

	userResponse := UserResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	fmt.Printf("Created user with email %s\n", user.Email)
	respondWithJSON(userResponse, http.StatusCreated, w)
}
