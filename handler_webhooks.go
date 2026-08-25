package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/reconfirmok/chirpy/internal/auth"
)

const (
	userEvent = "user.upgraded"
)

func (cfg *apiConfig) handlerWebhooks(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error fetching api key from header", err)
		return
	}
	if apiKey != cfg.polka {
		respondWithError(w, http.StatusUnauthorized, "API key is invalid", err)
		return
	}

	params := parameters{}
	decode := json.NewDecoder(r.Body)
	err = decode.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding params", err)
		return
	}

	if params.Event != userEvent {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.db.UpgradeToChirpyRed(r.Context(), params.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Couldn't find user", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "error updating user to red", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
