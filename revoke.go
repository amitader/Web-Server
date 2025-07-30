package main

import (
	"net/http"
	"github.com/amitader/web-Server/internal/auth"
)

func (cfg *apiConfig) revoke(w http.ResponseWriter, r *http.Request) {
	rToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't find token", err)
		return
	}
	_, err = cfg.db.RevokeToken(r.Context(), rToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke session", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}