package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ClemSK/chirpy/internal/database"
)

func (cfg *apiConfig) handleWebhook(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserId int `json:"user_id"`
		}
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode parameter")
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusOK, struct{}{})
		return
	}

	_, err = cfg.DB.UpgradeChirpyRed(params.Data.UserId)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
			respondWithError(w, http.StatusNotFound, "Could not find user")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Could not update user")
		return
	}
	respondWithJSON(w, http.StatusOK, struct{}{})
}
