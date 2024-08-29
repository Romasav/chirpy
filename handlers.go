package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Romasav/chirpy/database"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func handlerGetChirp(w http.ResponseWriter, db *database.DB) {
	chirps, err := db.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to load chirps")
		return
	}
	respondWithJSON(w, chirps, http.StatusOK)
}

func handlerGetChirpByID(w http.ResponseWriter, r *http.Request, db *database.DB) {
	chirpIDStr := r.PathValue("chirpID")
	chirpID, err := strconv.Atoi(chirpIDStr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to convert chirpIDStr to int")
		return
	}

	chirp, err := db.GetChirpByID(chirpID)
	if err != nil {
		errorMessage := fmt.Sprintf("The chirp with id = %v was not found", chirpID)
		respondWithError(w, http.StatusNotFound, errorMessage)
		return
	}

	respondWithJSON(w, chirp, http.StatusOK)
}

type chirpRequest struct {
	Body string `json:"body"`
}

func handlerPostChirp(w http.ResponseWriter, r *http.Request, db *database.DB) {
	decoder := json.NewDecoder(r.Body)
	request := chirpRequest{}
	err := decoder.Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	chirp, err := db.CreateChirp(request.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not create chirp")
		return
	}

	respondWithJSON(w, chirp, http.StatusCreated)
}

type userRequest struct {
	Email string `json:"email"`
}

func handlePostUser(w http.ResponseWriter, r *http.Request, db *database.DB) {
	decoder := json.NewDecoder(r.Body)
	request := userRequest{}
	err := decoder.Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	user, err := db.CreateUser(request.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldnt create user")
		return
	}

	respondWithJSON(w, user, http.StatusCreated)
}
