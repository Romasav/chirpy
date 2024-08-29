package main

import (
	"fmt"
	"net/http"
	"sync"
)

type apiConfig struct {
	fileserverHits int
	mutex          sync.Mutex
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.mutex.Lock()
		cfg.fileserverHits++
		cfg.mutex.Unlock()

		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.mutex.Lock()
	defer cfg.mutex.Unlock()
	cfg.fileserverHits = 0
}

func (cfg *apiConfig) handlerAdminMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := fmt.Sprintf(`
		<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
		</html>`, cfg.fileserverHits)

	w.Write([]byte(html))
}
