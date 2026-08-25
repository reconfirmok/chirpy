package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/reconfirmok/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	chirpIDStr := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error fetching JWT from header", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwt)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error validating JWT", err)
		return
	}

	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp is not found", err)
		return
	}

	if dbChirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Current user cant delete this chirp", err)
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
