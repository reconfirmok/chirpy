package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	_, err := w.Write([]byte("Hits reset to 0"))
	if err != nil {
		log.Fatalf("Error reseting hits: %v\n", err)
	}
}
