package main

import (
	"net/http"
)

func main() {
	serverMux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("."))

	serverMux.Handle("/", fileServer)

	server := http.Server{
		Addr:    ":8080",
		Handler: serverMux,
	}
	server.ListenAndServe()
}
