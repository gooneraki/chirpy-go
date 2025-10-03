package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gooneraki/chirpy-go/internal/database"
)

func (cfg *apiConfig) handlerUsersRed(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.db.UpdateUserRed(r.Context(), database.UpdateUserRedParams{
		ID:          params.Data.UserID,
		IsChirpyRed: true,
	})
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't update user red membership", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
