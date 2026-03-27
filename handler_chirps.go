package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

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

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
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
		UserID: params.UserID,
	})
	if err != nil {
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