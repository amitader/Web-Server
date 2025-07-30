package main

import (
	"net/http"
	"github.com/amitader/web-Server/internal/auth"
)

func (cfg *apiConfig) refresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}
	rToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find token", err)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), rToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get user for refresh token", err)
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
}