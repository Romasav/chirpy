package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Romasav/chirpy/database"
)

func main() {
	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if *dbg {
		fmt.Println("Debug mode enabled: Deleting database...")
		err := os.Remove("database.json")
		if err != nil && !os.IsNotExist(err) {
			fmt.Printf("Failed to delete database: %v\n", err)
		}

	}

	const filepathRoot = "."
	const port = "8080"
	apiConfig := apiConfig{}
	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal("Could not create new database")
	}

	serverMux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(filepathRoot))

	serverMux.Handle("/app/", apiConfig.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))
	serverMux.HandleFunc("GET /api/healthz", handlerReadiness)
	serverMux.HandleFunc("/api/reset", apiConfig.handlerReset)
	serverMux.HandleFunc("GET /admin/metrics", apiConfig.handlerAdminMetrics)
	serverMux.HandleFunc("POST /api/chirps", func(w http.ResponseWriter, r *http.Request) { handlerPostChirp(w, r, db) })
	serverMux.HandleFunc("GET /api/chirps", func(w http.ResponseWriter, r *http.Request) { handlerGetChirp(w, db) })
	serverMux.HandleFunc("GET /api/chirps/{chirpID}", func(w http.ResponseWriter, r *http.Request) { handlerGetChirpByID(w, r, db) })
	serverMux.HandleFunc("POST /api/users", func(w http.ResponseWriter, r *http.Request) { handlePostUser(w, r, db) })

	server := http.Server{
		Addr:    ":" + port,
		Handler: serverMux,
	}
	server.ListenAndServe()
}
