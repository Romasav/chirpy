package main

import (
	"fmt"
	"net/http"
	"sync"
)

func main() {
	const filepathRoot = "."
	const port = "8080"
	apiConfig := apiConfig{}

	serverMux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(filepathRoot))

	serverMux.Handle("/app/", apiConfig.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))
	serverMux.HandleFunc("GET /api/healthz", handlerReadiness)
	serverMux.HandleFunc("/api/reset", apiConfig.handlerReset)
	serverMux.HandleFunc("GET /admin/metrics", apiConfig.handlerAdminMetrics)

	server := http.Server{
		Addr:    ":" + port,
		Handler: serverMux,
	}
	server.ListenAndServe()
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

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

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hits: %d\n", cfg.fileserverHits)
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
