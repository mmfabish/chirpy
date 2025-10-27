package handlers

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) MetricsHandler(w http.ResponseWriter, req *http.Request) {
	responseTemplate := `
		<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
		</html>
		`
	responseMessage := fmt.Sprintf(responseTemplate, cfg.fileserverHits.Load())

	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(responseMessage))
}
