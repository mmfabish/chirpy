package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/google/uuid"
	"github.com/mmfabish/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	jwtSecret      string
	subject        uuid.UUID
}

type apiError struct {
	Message string `json:"message"`
}

func NewApiConfig(db *database.Queries, jwtSecret string) apiConfig {
	return apiConfig{
		db:        db,
		jwtSecret: jwtSecret,
	}
}

func RespondWithError(w http.ResponseWriter, statusCode int, message string) {
	respBody := apiError{Message: message}

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

func RespondWithJSON(w http.ResponseWriter, req *http.Request, statusCode int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}
