package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mmfabish/chirpy/internal/auth"
)

type UserLoginParameters struct {
	Password string `json:"password"`
	Email    string `json:"email"`
	Duration int    `json:"expires_in_seconds"`
}

type UserLoginResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

func (cfg *apiConfig) LoginHandler(w http.ResponseWriter, req *http.Request) {
	params := UserLoginParameters{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)

	if err != nil {
		log.Printf("Error decoding JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userEntity, err := cfg.db.GetUserByEmail(context.Background(), params.Email)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "incorrect email or password")
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, userEntity.HashedPassword)
	if err != nil || !match {
		RespondWithError(w, http.StatusUnauthorized, "incorrect email or password")
		return
	}

	var expiresIn time.Duration
	if params.Duration == 0 {
		expiresIn = time.Hour
	} else {
		expiresIn = time.Duration(params.Duration) * time.Second
	}

	token, err := auth.MakeJWT(userEntity.ID, cfg.jwtSecret, expiresIn)
	if err != nil {
		log.Printf("Error generating JWT: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	RespondWithJSON(w, req, http.StatusOK, UserLoginResponse{
		ID:        userEntity.ID,
		CreatedAt: userEntity.CreatedAt,
		UpdatedAt: userEntity.UpdatedAt,
		Email:     userEntity.Email,
		Token:     token,
	})
}
