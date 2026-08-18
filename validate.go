package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func handlerValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnValues struct {
		Valid bool `json:"valid"`
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

	respondWithJSON(w, http.StatusOK, returnValues{
		Valid: true,
	})
}

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}

	if code > 499 {
		log.Printf("Responding with 5XX error: %s\n", msg)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}
