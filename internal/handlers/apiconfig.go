package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/google/uuid"
	"github.com/mmfabish/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	jwtSecret      string
	polkaSecret    string
	subject        uuid.UUID
}

type apiError struct {
	Message string `json:"message"`
}

func NewApiConfig() (*apiConfig, error) {
	// connect to database
	dbUrl := os.Getenv("DB_URL")

	if dbUrl == "" {
		return nil, errors.New("missing environment variable: DB_URL")
	}
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, err
	}

	return &apiConfig{
		db:          database.New(db),
		jwtSecret:   os.Getenv("JWT_SECRET"),
		polkaSecret: os.Getenv("POLKA_KEY"),
	}, nil
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
