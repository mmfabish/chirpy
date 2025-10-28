package handlers

import (
	"context"
	"encoding/json"
	"io"
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

type UserDTO struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func mapUserEntityToDTO(userEntity *database.User) UserDTO {
	return UserDTO{
		ID:        userEntity.ID,
		CreatedAt: userEntity.CreatedAt,
		UpdatedAt: userEntity.UpdatedAt,
		Email:     userEntity.Email,
	}
}

func DecodeCreateUserParams(requestBody io.ReadCloser) (*CreateUserParams, error) {
	params := CreateUserParams{}

	decoder := json.NewDecoder(requestBody)
	if err := decoder.Decode(&params); err != nil {
		return nil, err
	}

	return &params, nil
}

func (cfg *apiConfig) UsersHandler(w http.ResponseWriter, req *http.Request) {
	// decode the request
	params, err := DecodeCreateUserParams(req.Body)
	if err != nil {
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
