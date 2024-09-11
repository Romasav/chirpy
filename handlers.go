package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Romasav/chirpy/database"
	"github.com/golang-jwt/jwt/v5"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func handlerGetChirp(w http.ResponseWriter, r *http.Request, db *database.DB) {
	chirps, err := db.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to load chirps")
		return
	}

	authorIDString := r.URL.Query().Get("author_id")
	if authorIDString != "" {
		authorID, err := strconv.Atoi(authorIDString)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldnt get authorId")
			return
		}
		filteredChirps := []database.Chirp{}
		for _, chirp := range chirps {
			if chirp.AuthorID == authorID {
				filteredChirps = append(filteredChirps, chirp)
			}
		}
		chirps = filteredChirps
	}

	sortMethod := r.URL.Query().Get("sort")
	if sortMethod == "" {
		sortMethod = "asc"
	}
	if sortMethod == "asc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID < chirps[j].ID
		})
	} else if sortMethod == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID > chirps[j].ID
		})
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

func handlerPostChirp(w http.ResponseWriter, r *http.Request, db *database.DB) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Authorization header is required")
		return
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		jwtSecret := os.Getenv("JWT_SECRET")
		return []byte(jwtSecret), nil
	})
	if err != nil || !token.Valid {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
		return
	}

	claims := token.Claims
	userIdString, err := claims.GetSubject()
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token claims")
		return
	}

	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt get userId")
		return
	}

	decoder := json.NewDecoder(r.Body)
	request := struct {
		Body string `json:"body"`
	}{}
	err = decoder.Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	chirp, err := db.CreateChirp(request.Body, userId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not create chirp")
		return
	}

	respondWithJSON(w, chirp, http.StatusCreated)
}

func handlerDeleteChirp(w http.ResponseWriter, r *http.Request, db *database.DB) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Authorization header is required")
		return
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		jwtSecret := os.Getenv("JWT_SECRET")
		return []byte(jwtSecret), nil
	})
	if err != nil || !token.Valid {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
		return
	}

	claims := token.Claims
	userIdString, err := claims.GetSubject()
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token claims")
		return
	}

	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt get userId")
		return
	}

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

	if chirp.AuthorID != userId {
		respondWithError(w, http.StatusForbidden, "you cant delete chirps that were created by someone else")
		return
	}

	err = db.DeleteChirpByID(chirpID)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "error deleting a chirp")
		return
	}

	respondWithJSON(w, struct{}{}, http.StatusNoContent)
}

func handlerPostUser(w http.ResponseWriter, r *http.Request, db *database.DB) {
	decoder := json.NewDecoder(r.Body)
	request := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	err := decoder.Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	user, err := db.CreateUser(request.Email, request.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldnt create user")
		return
	}

	userRespond := struct {
		ID          int    `json:"id"`
		Email       string `json:"email"`
		IsChirpyRed bool   `json:"is_chirpy_red"`
	}{
		ID:          user.ID,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}

	respondWithJSON(w, userRespond, http.StatusCreated)
}

func handlerLoginUser(w http.ResponseWriter, r *http.Request, db *database.DB) {
	decoder := json.NewDecoder(r.Body)
	request := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	err := decoder.Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	user, err := db.GetUserByEmail(request.Email)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "The user was not found")
		return
	}

	err = user.ComparePassword(request.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect Password")
		return
	}

	accessToken, err := generateJWT(user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting a token")
		return
	}

	refreshToken, err := db.CreateRefreshToken(user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating a refresh token")
		return
	}

	userRespond := struct {
		ID           int    `json:"id"`
		Email        string `json:"email"`
		IsChirpyRed  bool   `json:"is_chirpy_red"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}{
		ID:           user.ID,
		Email:        user.Email,
		IsChirpyRed:  user.IsChirpyRed,
		Token:        accessToken,
		RefreshToken: refreshToken.Token,
	}

	respondWithJSON(w, userRespond, http.StatusOK)
}

func generateJWT(userID int) (string, error) {
	expiresInSeconds := 3600
	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		Subject:   fmt.Sprintf("%d", userID),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Duration(expiresInSeconds) * time.Second)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func handlerUpdateUser(w http.ResponseWriter, r *http.Request, db *database.DB) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Authorization header is required")
		return
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		jwtSecret := os.Getenv("JWT_SECRET")
		return []byte(jwtSecret), nil
	})
	if err != nil || !token.Valid {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
		return
	}

	claims := token.Claims
	userIdString, err := claims.GetSubject()
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token claims")
		return
	}

	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt get userId")
		return
	}

	decoder := json.NewDecoder(r.Body)
	request := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	err = decoder.Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	updatedUser, err := database.NewUser(userId, request.Email, request.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt create updated user")
		return
	}

	err = db.UpdateUser(*updatedUser)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	userRespond := struct {
		ID          int    `json:"id"`
		Email       string `json:"email"`
		IsChirpyRed bool   `json:"is_chirpy_red"`
	}{
		ID:          updatedUser.ID,
		Email:       updatedUser.Email,
		IsChirpyRed: updatedUser.IsChirpyRed,
	}

	respondWithJSON(w, userRespond, http.StatusOK)
}

func handlerRefreshToken(w http.ResponseWriter, r *http.Request, db *database.DB) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Authorization header is required")
		return
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	refreshToken, err := db.GetRefreshTokenInfo(tokenStr)
	if err != nil || refreshToken.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "The refresh token is invalid")
		return
	}

	token, err := generateJWT(refreshToken.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting a token")
		return
	}

	respond := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}

	respondWithJSON(w, respond, http.StatusOK)
}

func handlerRevokeToken(w http.ResponseWriter, r *http.Request, db *database.DB) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Authorization header is required")
		return
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	refreshToken, err := db.GetRefreshTokenInfo(tokenStr)
	if err != nil || refreshToken.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "The refresh token is invalid")
		return
	}

	err = db.DeleteRefreshToken(refreshToken.Token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Failed to delete refresh token")
		return
	}

	respondWithJSON(w, struct{}{}, http.StatusNoContent)
}

func handlerWebhooks(w http.ResponseWriter, r *http.Request, db *database.DB) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Authorization header is required")
		return
	}

	apiKey := strings.TrimPrefix(authHeader, "ApiKey ")

	if apiKey != os.Getenv("POLKA_KEY") {
		respondWithError(w, http.StatusUnauthorized, "Incorrect key")
		return
	}

	decoder := json.NewDecoder(r.Body)
	request := struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}{}
	err := decoder.Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if request.Event != "user.upgraded" {
		respondWithJSON(w, struct{}{}, http.StatusNoContent)
		return
	}

	_, err = db.UpgradeToChirpyRed(request.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "the user was not found")
	}

	respondWithJSON(w, struct{}{}, http.StatusNoContent)
}
