package main

import (
	"net/http"

	"github.com/gooneraki/chirpy-go/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't get bearer token", err)
		return
	}

	_, err = cfg.db.RevokeRefreshToken(r.Context(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't revoke token", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, struct{}{})

}
