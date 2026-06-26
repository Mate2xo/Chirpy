package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Mate2xo/Chirpy/internal/auth"
	"github.com/Mate2xo/Chirpy/internal/database"
)

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

	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret)
	if err != nil {
		respondWithErr(err, http.StatusInternalServerError, w)
		return
	}
	refreshToken, err := cfg.dbQueries.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token:     auth.MakeRefrehToken(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Local().Add(60 * time.Hour * 24),
	})
	if err != nil {
		respondWithErr(err, http.StatusInternalServerError, w)
		return
	}

	payload := UserResponse{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
	}
	respondWithJSON(payload, http.StatusOK, w)
}

func (cfg *apiConfig) refreshUser(w http.ResponseWriter, req *http.Request) {
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithErr(err, http.StatusBadRequest, w)
		return
	}

	user, err := cfg.dbQueries.GetUserByRefreshToken(req.Context(), refreshToken)
	if err != nil {
		respondWithErr(err, http.StatusUnauthorized, w)
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret)
	if err != nil {
		respondWithErr(err, http.StatusInternalServerError, w)
		return
	}

	payload := UserResponse{AccessToken: accessToken}
	respondWithJSON(payload, http.StatusOK, w)
}

func (cfg *apiConfig) revokeUser(w http.ResponseWriter, req *http.Request) {
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithErr(err, http.StatusBadRequest, w)
		return
	}

	err = cfg.dbQueries.RevokeRefreshToken(req.Context(), refreshToken)
	if err != nil {
		respondWithErr(err, http.StatusInternalServerError, w)
		return
	}
	respondWithJSON(nil, http.StatusNoContent, w)
}
