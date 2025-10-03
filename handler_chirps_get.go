package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gooneraki/chirpy-go/internal/database"
)

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	authorId := r.URL.Query().Get("author_id")

	var dbChirps []database.Chirp
	var err error
	if authorId == "" {
		dbChirps, err = cfg.db.GetChirps(r.Context())
	} else {
		authorUUID, parseErr := uuid.Parse(authorId)
		if parseErr != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID", parseErr)
			return
		}
		dbChirps, err = cfg.db.GetChirpsByAuthor(r.Context(), authorUUID)
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
		return
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			UserID:    dbChirp.UserID,
			Body:      dbChirp.Body,
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}
