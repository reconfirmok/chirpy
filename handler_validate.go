package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnValues struct {
		CleanBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	parmas := parameters{}
	err := decoder.Decode(&parmas)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(parmas.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return

	}

	cleaned := cleanChirpBody(parmas.Body)
	respondWithJSON(w, http.StatusOK, returnValues{
		CleanBody: cleaned,
	})
}

func cleanChirpBody(body string) string {
	profaneWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	words := strings.Split(body, " ")
	for i, word := range words {
		wordToLower := strings.ToLower(word)
		if _, ok := profaneWords[wordToLower]; ok {
			words[i] = "****"
		}
	}

	cleaned := strings.Join(words, " ")
	return cleaned
}
