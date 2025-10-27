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

type CreateChirpParams struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type ChirpDTO struct {
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

func mapChirpEntityToDTO(chirpEntity *database.Chirp) ChirpDTO {
	return ChirpDTO{
		ID:        chirpEntity.ID,
		CreatedAt: chirpEntity.CreatedAt,
		UpdatedAt: chirpEntity.UpdatedAt,
		Body:      chirpEntity.Body,
		UserID:    chirpEntity.UserID,
	}
}

func (cfg *apiConfig) CreateChirpHandler(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	params := CreateChirpParams{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding request: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(params.Body) > 140 {
		RespondWithError(w, http.StatusBadRequest, "Chirp is too long")
	} else {
		params := database.CreateChirpParams{
			Body:   filterMessage(params.Body),
			UserID: params.UserID,
		}

		chirp, err := cfg.db.CreateChirp(context.Background(), params)
		if err != nil {
			log.Printf("Error creating chirp: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := mapChirpEntityToDTO(&chirp)

		RespondWithJSON(w, req, http.StatusCreated, response)
	}
}

func (cfg *apiConfig) GetChirpsHandler(w http.ResponseWriter, req *http.Request) {
	chirps, err := cfg.db.GetAllChirps(context.Background())
	if err != nil {
		log.Printf("Error retrieving chirps: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := make([]ChirpDTO, len(chirps))
	for i, chirp := range chirps {
		data[i] = mapChirpEntityToDTO(&chirp)
	}

	RespondWithJSON(w, req, http.StatusOK, data)
}

func (cfg *apiConfig) GetChirpHandler(w http.ResponseWriter, req *http.Request) {
	chirpID, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		log.Printf("Error converting string to UUID: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	chirp, err := cfg.db.GetChirpByID(context.Background(), chirpID)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Not found")
		return
	}

	RespondWithJSON(w, req, http.StatusOK, mapChirpEntityToDTO(&chirp))
}
