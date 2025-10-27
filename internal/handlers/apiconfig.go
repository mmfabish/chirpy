package handlers

import (
	"sync/atomic"

	"github.com/mmfabish/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
}

func NewApiConfig(db *database.Queries) apiConfig {
	return apiConfig{
		db: db,
	}
}
