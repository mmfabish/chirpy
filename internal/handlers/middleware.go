package handlers

import (
	"log"
	"net/http"

	"github.com/mmfabish/chirpy/internal/auth"
)

func (cfg *apiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}

func (cfg *apiConfig) MiddlewareBearerAuth(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		bearerToken, err := auth.GetBearerToken(req.Header)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))

			return
		}

		subject, err := auth.ValidateJWT(bearerToken, cfg.jwtSecret)
		if err != nil {
			log.Printf("Failed login for %s.", subject)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))

			return
		}

		cfg.subject = subject

		log.Printf("Successful login for %s", subject)

		next(w, req)
	})
}
