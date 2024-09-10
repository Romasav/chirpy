package database

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

type Chirp struct {
	ID       int    `json:"id"`
	Body     string `json:"body"`
	AuthorID int    `json:"author_id"`
}

func NewChirp(body string, id int, authorID int) (*Chirp, error) {
	validatedBody, err := validateChirp(body)
	if err != nil {
		return nil, err
	}

	newChirp := Chirp{
		ID:       id,
		Body:     validatedBody,
		AuthorID: authorID,
	}

	return &newChirp, nil
}

func validateChirp(chirp string) (string, error) {
	chirpLength := utf8.RuneCountInString(chirp)
	if chirpLength > 140 {
		errorMessage := fmt.Sprintf("The chirp(len = %v) exeeds the rune limit of 140", chirpLength)
		return "", errors.New(errorMessage)
	}

	cleanedChirp := cleanChirp(chirp)

	return cleanedChirp, nil
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
