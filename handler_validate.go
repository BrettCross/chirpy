package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

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

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	profanity := map[string]struct{}{
		"kerfuffle": {},
		"sharbert": {},
		"fornax": {},
	}

	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
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
	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: cleanedChirp,
	})
}