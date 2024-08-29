package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"unicode/utf8"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

type chirpRequest struct {
	Body string `json:"body"`
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	request := chirpRequest{}
	err := decoder.Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	chirp := request.Body

	if utf8.RuneCountInString(chirp) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanedChirp := cleanChirp(chirp)

	response := map[string]string{"cleaned_body": cleanedChirp}
	respondWithJSON(w, response, http.StatusOK)
}

func cleanChirp(chirp string) string {
	wordsToCensor := [3]string{"kerfuffle", "sharbert", "fornax"}
	chirpWords := strings.Split(chirp, " ")
	for chirpyIndex, chirpWord := range chirpWords {
		for _, wordToCensor := range wordsToCensor {
			if strings.ToLower(chirpWord) == wordToCensor {
				chirpWords[chirpyIndex] = "****"
			}
		}
	}
	return strings.Join(chirpWords, " ")
}
