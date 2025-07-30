package main

import (
	"net/http"
)
func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Forbidden\n"))
	}
	cfg.fileserverHits.Store(0)
	cfg.db.DeleteAllUsers(r.Context())
	w.WriteHeader(http.StatusOK)
    w.Write([]byte("Counter reset\n"))
}