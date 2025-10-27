package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/mmfabish/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
}

func NewApiConfig(db *database.Queries) apiConfig {
	return apiConfig{
		db: db,
	}
}

func RespondWithError(w http.ResponseWriter, statusCode int, message string) {
	respBody := error{Message: message}

	data, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}

func RespondWithJSON(w http.ResponseWriter, req *http.Request, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
