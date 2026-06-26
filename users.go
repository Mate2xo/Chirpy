package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Mate2xo/Chirpy/internal/auth"
	"github.com/Mate2xo/Chirpy/internal/database"
)

type UserResponse struct {
	database.User
	AccessToken  string `json:"token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
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
		User: database.User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			Email:     user.Email,
		},
	}

	fmt.Printf("Created user with email %s\n", user.Email)
	respondWithJSON(userResponse, http.StatusCreated, w)
}
