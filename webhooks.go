package main

import (
	"encoding/json"
	"net/http"

	"github.com/Mate2xo/Chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) upgradeUserHook(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		}
	}

	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil || apiKey != cfg.polkaKey {
		respondWithErr(err, http.StatusUnauthorized, w)
		return
	}

	params := parameters{}
	err = json.NewDecoder(req.Body).Decode(&params)
	if err != nil {
		respondWithErr(err, http.StatusBadRequest, w)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = cfg.dbQueries.UpgradeUser(req.Context(), params.Data.UserID)
	if err != nil {
		respondWithErr(err, http.StatusNotFound, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
