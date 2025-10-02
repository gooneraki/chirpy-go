package main

import (
	"net/http"
	"time"

	"github.com/gooneraki/chirpy-go/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't get bearer token", err)
		return
	}

	refreshTokenDb, err := cfg.db.GetRefreshTokenByToken(r.Context(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't find refresh token in db", err)
		return
	}

	if refreshTokenDb.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "refresh token has expired", nil)
		return
	}

	if refreshTokenDb.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "refresh token has been revoked", nil)
		return
	}

	// Create a new JWT access token
	accessToken, err := auth.MakeJWT(refreshTokenDb.UserID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create JWT", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})

}
