package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/reconfirmok/chirpy/internal/database"
)

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	authorID, err := authorIDFromRequest(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid author ID", err)
		return
	}

	var dbChrips []database.Chirp

	if authorID != uuid.Nil {
		dbChrips, err = cfg.db.GetChirpsByAuthor(r.Context(), authorID)
	} else {
		dbChrips, err = cfg.db.GetChirps(r.Context())
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error retrieving chrips", err)
		return
	}

	chirpsList := make([]Chirp, len(dbChrips))
	for i, chrip := range dbChrips {
		chirpsList[i] = Chirp{
			ID:        chrip.ID,
			CreatedAt: chrip.CreatedAt,
			UpdatedAt: chrip.UpdatedAt,
			Body:      chrip.Body,
			UserID:    chrip.UserID,
		}
	}

	respondWithJSON(w, http.StatusOK, chirpsList)
}

func authorIDFromRequest(r *http.Request) (uuid.UUID, error) {
	authorIDStr := r.URL.Query().Get("author_id")
	if authorIDStr == "" {
		return uuid.Nil, nil
	}

	authorId, err := uuid.Parse(authorIDStr)
	if err != nil {
		return uuid.Nil, err
	}
	return authorId, nil
}
