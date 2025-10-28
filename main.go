package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mmfabish/chirpy/internal/database"
	"github.com/mmfabish/chirpy/internal/handlers"
)

const filepathRoot = "."

func main() {
	// load environment variables
	godotenv.Load()

	// connect to database
	dbUrl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	cfg := handlers.NewApiConfig(database.New(db))
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

	// user endpoints
	mux.HandleFunc("POST /api/login", cfg.LoginHandler)
	mux.HandleFunc("POST /api/users", cfg.UsersHandler)

	// chirps endpoints
	mux.HandleFunc("GET /api/chirps", cfg.GetChirpsHandler)
	mux.HandleFunc("POST /api/chirps", cfg.CreateChirpHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.GetChirpHandler)

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	server.ListenAndServe()
}
