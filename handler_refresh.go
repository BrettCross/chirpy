package main

import (
	"net/http"
	"time"

	"github.com/brettcross/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	headerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't get token from header", err)
		return
	}

	dbToken, err := cfg.db.GetRefreshToken(r.Context(), headerToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "token doesn't exist", err)
		return
	}

	if dbToken.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "token expired", err)
		return
	}

	if dbToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "token revoked", err)
		return
	}

	jwtToken, err := auth.MakeJWT(dbToken.UserID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't make JWT token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: jwtToken,
	})

}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	headerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't get token from header", err)
		return
	}
	err = cfg.db.RevokeRefreshToken(r.Context(), headerToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't revoke token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}