package main

import "net/http"

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	dbChrips, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error Retrieving Chrips", err)
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
