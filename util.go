package main

import (
	"encoding/json"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, `{"error": "Failed to encode JSON"}`, http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}

func respondWithError(w http.ResponseWriter, statusCode int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	errorJson := map[string]string{"error": msg}

	jsonData, err := json.Marshal(errorJson)
	if err != nil {
		http.Error(w, `{"error": "Failed to encode JSON"}`, http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}
