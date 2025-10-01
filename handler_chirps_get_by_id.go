package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsRetrieveById(w http.ResponseWriter, r *http.Request) {

	chirpIDstring := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDstring)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	dbChirps, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't retrieve chirp with id", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        dbChirps.ID,
		CreatedAt: dbChirps.CreatedAt,
		UpdatedAt: dbChirps.UpdatedAt,
		UserID:    dbChirps.UserID,
		Body:      dbChirps.Body,
	})
}
