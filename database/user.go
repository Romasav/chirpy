package database

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

func NewUser(id int, email, password string) (*User, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	newUser := User{
		ID:          id,
		Email:       email,
		Password:    hashedPassword,
		IsChirpyRed: false,
	}
	return &newUser, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("error hashing the password")
	}
	return string(hashedPassword), nil
}

func (user *User) ComparePassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return errors.New("incorrect password")
	}
	return nil
}
