package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/reconfirmok/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAT time.Time `json:"created_at"`
	UpdatedAT time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChrips(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	parmas := parameters{}
	err := decoder.Decode(&parmas)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}

	validChrip, err := validateChirp(parmas.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	chrip, err := cfg.db.CreateChrip(r.Context(), database.CreateChripParams{
		Body:   validChrip,
		UserID: parmas.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating chrip", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chrip.ID,
		CreatedAT: chrip.CreatedAt,
		UpdatedAT: chrip.UpdatedAt,
		Body:      chrip.Body,
		UserID:    chrip.UserID,
	})

}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chrips is too long")
	}

	cleaned := cleanChirpBody(body)
	return cleaned, nil

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
