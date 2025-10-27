package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mmfabish/chirpy/internal/database"
)

type createChirpRequest struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type createChirpResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

type error struct {
	Message string `json:"error"`
}

func filterMessage(message string) string {
	prohibitedWords := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(message, " ")

	for i := 0; i < len(words); i++ {
		if slices.Contains(prohibitedWords, strings.ToLower(words[i])) {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}

func (cfg *apiConfig) CreateChirpHandler(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	chirpRequest := createChirpRequest{}
	err := decoder.Decode(&chirpRequest)
	if err != nil {
		log.Printf("Error decoding request: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(chirpRequest.Body) > 140 {
		RespondWithError(w, http.StatusBadRequest, "Chirp is too long")
	} else {
		params := database.CreateChirpParams{
			Body:   filterMessage(chirpRequest.Body),
			UserID: chirpRequest.UserID,
		}

		chirp, err := cfg.db.CreateChirp(context.Background(), params)
		if err != nil {
			log.Printf("Error creating user: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := createChirpResponse{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
		RespondWithJSON(w, req, http.StatusCreated, response)
	}
}
