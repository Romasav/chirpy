package main

import (
	"net/http"
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
	serverMux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	server := http.Server{
		Addr:    ":" + port,
		Handler: serverMux,
	}
	server.ListenAndServe()
}
