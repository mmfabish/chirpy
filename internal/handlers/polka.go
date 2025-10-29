package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type PolkaWebhookPayload struct {
	Event string `json:"event"`
	Data  struct {
		UserID uuid.UUID `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) PolkaWebhookHandler(w http.ResponseWriter, req *http.Request) {
	// decode the payload
	payload := PolkaWebhookPayload{}
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&payload); err != nil {
		log.Fatal("Error occurred while parsing PolkaWebhook payload: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// if event is anything besides user.upgraded, return a 204
	if payload.Event != "user.upgraded" {
		log.Printf("Unsupported event type: %s", payload.Event)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if err := cfg.db.UpgradeUserToChirpyRed(context.Background(), payload.Data.UserID); err != nil {
		log.Printf("User ID not found: %s", payload.Data.UserID)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
