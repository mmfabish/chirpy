package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mmfabish/chirpy/internal/database"
)

type CreateChirpParams struct {
	Body string `json:"body"`
}

type ChirpDTO struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
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
			UserID: cfg.subject,
		}

		log.Printf("Creating chirp for user %s", cfg.subject)
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
	var chirps []database.Chirp
	var err error

	authorId := req.URL.Query().Get("author_id")
	if authorId != "" {
		log.Printf("Getting chirps for author %s", authorId)

		userID, err := uuid.Parse(authorId)
		if err != nil {
			log.Printf("Error parsing user ID from author_id parameter: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		chirps, err = cfg.db.GetChirpsByAuthorID(context.Background(), userID)
		if err != nil {
			log.Printf("Error retrieving chirps for author %s: %s", authorId, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		chirps, err = cfg.db.GetAllChirps(context.Background())
		if err != nil {
			log.Printf("Error retrieving chirps: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	data := make([]ChirpDTO, len(chirps))
	for i, chirp := range chirps {
		data[i] = mapChirpEntityToDTO(&chirp)
	}

	sortDirection := strings.ToLower(req.URL.Query().Get("sort"))
	switch sortDirection {
	case "asc":
		sort.Slice(data, func(i, j int) bool { return data[i].CreatedAt.Before(data[j].CreatedAt) })

	case "desc":
		sort.Slice(data, func(i, j int) bool { return data[j].CreatedAt.Before(data[i].CreatedAt) })
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

func (cfg *apiConfig) DeleteChirpHandler(w http.ResponseWriter, req *http.Request) {
	chirpID, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		log.Printf("Error converting string to UUID: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	chirp, err := cfg.db.GetChirpByID(context.Background(), chirpID)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Chirp does not exist.")
	}

	if chirp.UserID != cfg.subject {
		log.Printf("Unauthorized removal attempt of chrip %s by user %s", chirpID, cfg.subject)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if err := cfg.db.DeleteChirp(context.Background(), chirpID); err != nil {
		log.Printf("Error deleting chirp: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
