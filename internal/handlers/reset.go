package handlers

import (
	"context"
	"net/http"
)

func (cfg *apiConfig) ResetHandler(w http.ResponseWriter, req *http.Request) {
	// reset the hit counter
	cfg.fileserverHits.Store(0)

	// delete all users in database
	cfg.db.DeleteAllUsers(context.Background())

	w.WriteHeader(http.StatusOK)
}
