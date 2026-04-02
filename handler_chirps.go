package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/brettcross/chirpy/internal/auth"
	"github.com/brettcross/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID 	  uuid.UUID `json:"user_id"`
}

func cleanChirp(msg string, words map[string]struct{}) string {
	profane := "****"
	msgWords := strings.Split(msg, " ")
	for i, msgWord := range msgWords {
		if _, ok := words[strings.ToLower(msgWord)]; ok {
			msgWords[i] = profane
		}
	}
	return strings.Join(msgWords, " ")
}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	profanity := map[string]struct{}{
		"kerfuffle": {},
		"sharbert": {},
		"fornax": {},
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return 
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't get bearer token", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't validate JWT token", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	cleanedChirp := cleanChirp(params.Body, profanity)
	// save to db
	// respond with JSON
	newChirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: cleanedChirp,
		UserID: userID,
	})
	if err != nil {
		// TODO: replace with respondWithError
		log.Printf("Error creating chirp in db: %s", err)
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID: newChirp.ID,
		CreatedAt: newChirp.CreatedAt,
		UpdatedAt: newChirp.UpdatedAt,
		Body: newChirp.Body,
		UserID: newChirp.UserID,
	})
}

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		// TODO: replace with respondWithError
		log.Printf("Error retrieving chirps from db: %s", err)
	}
	
	responseChirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		responseChirps = append(responseChirps, Chirp{
			ID: dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body: dbChirp.Body,
			UserID: dbChirp.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK, responseChirps)
}

func (cfg *apiConfig) handlerChirpsRetrieveByID(w http.ResponseWriter, r *http.Request) {
	chirpIDStr := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}
	dbChirp, err := cfg.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID: dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body: dbChirp.Body,
		UserID: dbChirp.UserID,
	})

}

func (cfg *apiConfig) handlerChirpsDeleteByID(w http.ResponseWriter, r *http.Request) {
	// authentication 
	authToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't get auth token", err)
		return
	}

	// validation
	userID, err := auth.ValidateJWT(authToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token", err)
		return
	}
	
	// get chirp ID from path
	chirpIDStr := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	// get chirp from db
	dbChirp, err := cfg.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	// verify chirp belongs to requesting user
	if dbChirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "", err)
		return
	}

	err = cfg.db.DeleteChirpByID(r.Context(), dbChirp.ID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "coudn't find chirp", err)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}