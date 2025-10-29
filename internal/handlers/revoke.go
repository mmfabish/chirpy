package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/mmfabish/chirpy/internal/auth"
)

func (cfg *apiConfig) RevokeHandler(w http.ResponseWriter, req *http.Request) {
	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if err := cfg.db.RevokeToken(context.Background(), bearerToken); err == nil {
		w.WriteHeader(http.StatusNoContent)
	} else {
		log.Printf("Error revoking token: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
