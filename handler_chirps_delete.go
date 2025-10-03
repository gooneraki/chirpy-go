package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gooneraki/chirpy-go/internal/auth"
	"github.com/gooneraki/chirpy-go/internal/database"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userId, err := auth.ValidateJWT(tokenString, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Couldn't validate JWT", err)
		return
	}

	// First, check if the chirp exists and get its details
	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}

	// Check if the user is the author of the chirp
	if chirp.UserID != userId {
		respondWithError(w, http.StatusForbidden, "You can only delete your own chirps", nil)
		return
	}

	// Now delete the chirp
	err = cfg.db.DeleteChirp(r.Context(), database.DeleteChirpParams{
		ID:     chirpID,
		UserID: userId,
	})
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't delete chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
