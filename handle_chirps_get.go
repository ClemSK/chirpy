package main

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not retrieve chirps")
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:   dbChirp.ID,
			Body: dbChirp.Body,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerChirpsGetById(w http.ResponseWriter, r *http.Request) {
	// dbChirps, err := cfg.DB.GetChirps()
	// if err != nil {
	// 	respondWithError(w, http.StatusInternalServerError, "Could not retrieve chirps")
	// }
	chirpIdStr := chi.URLParam(r, "id")

	chirpId, err := strconv.Atoi(chirpIdStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp id")
		return
	}

	chirp, err := cfg.DB.GetChirp(chirpId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}

	// respondWithJSON(w, http.StatusOK, chirp)
	respondWithJSON(w, http.StatusOK, Chirp{
		ID:   chirp.ID,
		Body: chirp.Body,
	})
}
