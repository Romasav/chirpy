package main

import (
	"net/http"
)

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	serverMux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(filepathRoot))

	serverMux.Handle("/app/", http.StripPrefix("/app", fileServer))
	serverMux.HandleFunc("/healthz", healthzHandler)

	server := http.Server{
		Addr:    ":" + port,
		Handler: serverMux,
	}
	server.ListenAndServe()
}
