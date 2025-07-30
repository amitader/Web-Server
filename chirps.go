package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"errors"
	"github.com/google/uuid"
	"github.com/amitader/web-Server/internal/database"
	"github.com/amitader/web-Server/internal/auth"
	"time"
	"sort"
)
type Chirp struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"update_at"`
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

func (cfg *apiConfig) ChirpsCreation(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	cleaned, err := handlerChirpsValidate(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest,err.Error(), err)
		return
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: cleaned,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError,"couldnt create chirp", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID:chirp.UserID,
	})
}

func (cfg *apiConfig) ChirpsDeletion(w http.ResponseWriter, r *http.Request) {
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find token", err)
		return
	}
	user_id, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}
	stringID := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(stringID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Could not parse chirp ID", err)
		return
	}
	chirp, err := cfg.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound,"Chirp not found", err)
		return
	}
	if user_id != chirp.UserID {
		respondWithError(w, http.StatusForbidden, "Chirp does not belong to user", errors.New("unauthorized deletion attempt"))
		return
	}
	err = cfg.db.DeleteChirp(r.Context(), database.DeleteChirpParams{
		ID: chirpID,
		UserID: user_id,
	})
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Could not delete chirp", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}

func handlerChirpsValidate(body string) (string, error){
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}
	return removeProfane(body), nil
}
func removeProfane (b string) string {
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	words := strings.Fields(b)
	censordWords := []string{}
	for _, word := range words {
		if _, ok := badWords[strings.ToLower(word)]; ok {
			censordWords = append(censordWords, "****")
		} else {
			censordWords = append(censordWords, word)
		}
	}
	return strings.Join(censordWords, " ")

}

func (cfg *apiConfig) getAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError,"couldnt get all chirps", err)
		return
	}
	authorID := uuid.Nil
	authorIDString := r.URL.Query().Get("author_id")
	if authorIDString != "" {
		authorID, err = uuid.Parse(authorIDString)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID", err)
			return
		}
	}
	sortParam := r.URL.Query().Get("sort")
	if sortParam == "" {
		sortParam = "asc"
	}


	resp := []Chirp{}
	for _, chirp := range chirps {
		if authorID != uuid.Nil && chirp.UserID != authorID {
			continue
		}
		val := Chirp{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserID:chirp.UserID,
		}
		resp = append(resp,val)
	}
	sort.Slice(resp, func(i, j int) bool {
		if sortParam == "asc" {
			return resp[i].CreatedAt.Before(resp[j].CreatedAt)
		} 
		return resp[i].CreatedAt.After(resp[j].CreatedAt)
	})

	respondWithJSON(w, http.StatusOK,resp)
}

func (cfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request) {
	stringID := r.PathValue("chirpID")
	parsedID, err := uuid.Parse(stringID)
	if err != nil {
		respondWithError(w, http.StatusNotFound,"couldnt parsed chirp id", err)
		return
	}
	chirp, err := cfg.db.GetChirpByID(r.Context(), parsedID)
	if err != nil {
		respondWithError(w, http.StatusNotFound,"chirp id not found", err)
		return
	}
	respondWithJSON(w, http.StatusOK,Chirp{
		ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserID:chirp.UserID,
	})
}