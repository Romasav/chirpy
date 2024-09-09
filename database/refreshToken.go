package database

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

type RefreshToken struct {
	UserID    int       `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

func NewRefreshToken(id int) (*RefreshToken, error) {
	token, err := generateRefreshToken()
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(60 * 24 * time.Hour)

	return &RefreshToken{
		UserID:    id,
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

func generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
