package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {

	chirpsData, err := cfg.db.GetChirps(r.Context())

	if err != nil {
		respondWithError(w, http.StatusServiceUnavailable, "couldn't get all chips", err)
		return
	}
	chirps := make([]Chirp, 0)
	for i := range chirpsData {
		chirps = append(chirps, Chirp{
			ID:        chirpsData[i].ID,
			CreatedAt: chirpsData[i].CreatedAt,
			UpdatedAt: chirpsData[i].UpdatedAt,
			UserID:    chirpsData[i].UserID,
			Body:      chirpsData[i].Body,
		})
	}

	fmt.Printf("chirps: %v\n", chirps)

	respondWithJSON(w, http.StatusOK, chirps)

}
