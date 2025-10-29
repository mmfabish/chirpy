package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mmfabish/chirpy/internal/auth"
	"github.com/mmfabish/chirpy/internal/database"
)

type CreateUserParams struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UpdateUserParams struct {
	CreateUserParams
}

type UserDTO struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

func mapUserEntityToDTO(userEntity *database.User) UserDTO {
	return UserDTO{
		ID:          userEntity.ID,
		CreatedAt:   userEntity.CreatedAt,
		UpdatedAt:   userEntity.UpdatedAt,
		Email:       userEntity.Email,
		IsChirpyRed: userEntity.IsChirpyRed.Bool,
	}
}

func (cfg *apiConfig) UsersHandler(w http.ResponseWriter, req *http.Request) {
	// decode the request
	params := CreateUserParams{}

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&params); err != nil {
		log.Printf("Error decoding JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// create the user
	hashed_password, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dbCreateUserParams := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashed_password,
	}

	userEntity, err := cfg.db.CreateUser(context.Background(), dbCreateUserParams)
	if err != nil {
		log.Printf("Error creating user: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userDTO := mapUserEntityToDTO(&userEntity)

	RespondWithJSON(w, req, http.StatusCreated, userDTO)
}

func (cfg *apiConfig) UpdateUsersHandler(w http.ResponseWriter, req *http.Request) {
	// decode the request
	params := UpdateUserParams{}

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&params); err != nil {
		log.Printf("Error decoding JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	hashed_password, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := cfg.db.UpdateUser(context.Background(), database.UpdateUserParams{
		Email:          params.Email,
		HashedPassword: hashed_password,
		ID:             cfg.subject,
	}); err != nil {
		log.Printf("Error hashing password: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userEntity, err := cfg.db.GetUserByEmail(req.Context(), params.Email)
	if err != nil {
		log.Printf("Error retrieving User entity: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	RespondWithJSON(w, req, http.StatusOK, mapUserEntityToDTO(&userEntity))
}
