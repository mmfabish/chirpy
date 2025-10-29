package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/mmfabish/chirpy/internal/auth"
)

type RefreshResponse struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) RefreshHandler(w http.ResponseWriter, req *http.Request) {
	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	refreshToken, err := cfg.db.GetRefreshTokenByToken(context.Background(), bearerToken)
	if err != nil || time.Now().After(refreshToken.ExpiresAt) || refreshToken.RevokedAt.Valid {
		RespondWithError(w, http.StatusUnauthorized, "")
		return
	}

	accessToken, err := auth.MakeJWT(refreshToken.UserID, cfg.jwtSecret, time.Hour)
	if err != nil {
		log.Printf("Error generating JWT: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	RespondWithJSON(w, req, http.StatusOK, RefreshResponse{
		Token: accessToken,
	})
}
