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

type UserLoginParameters struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserLoginResponse struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
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

	// create access token
	accessToken, err := auth.MakeJWT(userEntity.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		log.Printf("Error generating JWT: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// create refresh token
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("Error generating refresh token: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cfg.db.CreateRefreshToken(context.Background(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    userEntity.ID,
		ExpiresAt: time.Now().Add(time.Hour * 1440),
	})

	RespondWithJSON(w, req, http.StatusOK, UserLoginResponse{
		ID:           userEntity.ID,
		CreatedAt:    userEntity.CreatedAt,
		UpdatedAt:    userEntity.UpdatedAt,
		Email:        userEntity.Email,
		IsChirpyRed:  userEntity.IsChirpyRed.Bool,
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
}
