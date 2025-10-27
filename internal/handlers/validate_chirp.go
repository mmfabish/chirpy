package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

type parameters struct {
	Body string `json:"body"`
}

type response struct {
	CleanedBody string `json:"cleaned_body"`
}

type error struct {
	Message string `json:"error"`
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

func (cfg *apiConfig) ValidateChirpHandler(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding request: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(params.Body) > 140 {
		RespondWithError(w, http.StatusBadRequest, "Chirp is too long")
	} else {
		payload := response{
			CleanedBody: filterMessage(params.Body),
		}
		RespondWithJSON(w, req, http.StatusOK, payload)
	}
}
