package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mmfabish/chirpy/internal/handlers"
)

const filepathRoot = "."

func main() {
	// load environment variables
	godotenv.Load()

	cfg, err := handlers.NewApiConfig()
	if err != nil {
		log.Fatal(err)
		return
	}

	mux := http.NewServeMux()

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))

	// admin endpoints
	mux.HandleFunc("GET /admin/metrics", cfg.MetricsHandler)

	// only support reset in dev environment
	if os.Getenv("PLATFORM") == "dev" {
		mux.HandleFunc("POST /admin/reset", cfg.ResetHandler)
	}

	// api endpoints
	mux.Handle("/app/", cfg.MiddlewareMetricsInc(handler))
	mux.HandleFunc("GET /api/healthz", cfg.HealthCheckHandler)

	// auth endpoints
	mux.HandleFunc("POST /api/login", cfg.LoginHandler)
	mux.HandleFunc("POST /api/refresh", cfg.RefreshHandler)
	mux.HandleFunc("POST /api/revoke", cfg.RevokeHandler)

	// user endpoints
	mux.HandleFunc("POST /api/users", cfg.UsersHandler)
	mux.HandleFunc("PUT /api/users", cfg.MiddlewareBearerAuth(cfg.UpdateUsersHandler))

	// chirps endpoints
	mux.HandleFunc("GET /api/chirps", cfg.GetChirpsHandler)
	mux.HandleFunc("POST /api/chirps", cfg.MiddlewareBearerAuth(cfg.CreateChirpHandler))
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.GetChirpHandler)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.MiddlewareBearerAuth(cfg.DeleteChirpHandler))

	// polka endpoints
	mux.HandleFunc("POST /api/polka/webhooks", cfg.MiddlewareApiKeyAuth(cfg.PolkaWebhookHandler))

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	server.ListenAndServe()
}
