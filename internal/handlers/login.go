package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/mmfabish/chirpy/internal/auth"
)

func (cfg *apiConfig) LoginHandler(w http.ResponseWriter, req *http.Request) {
	params, err := DecodeCreateUserParams(req.Body)
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

	userDTO := mapUserEntityToDTO(&userEntity)
	RespondWithJSON(w, req, http.StatusOK, userDTO)
}
