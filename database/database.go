package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"sync"
)

type DBStructure struct {
	Chirps        map[int]Chirp        `json:"chirps"`
	Users         map[int]User         `json:"users"`
	RefreshTokens map[int]RefreshToken `json:"refresh_tokens"`
}

func NewDBStructure(chirps map[int]Chirp, users map[int]User, refreshTokens map[int]RefreshToken) (*DBStructure, error) {
	newDBStructure := DBStructure{
		Chirps:        chirps,
		Users:         users,
		RefreshTokens: refreshTokens,
	}
	return &newDBStructure, nil
}

type DB struct {
	path string
	mux  *sync.RWMutex
}

func NewDB(path string) (*DB, error) {
	db := DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	err := db.ensureDB()
	if err != nil {
		return nil, err
	}

	return &db, nil
}

func (db *DB) CreateUser(email, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	newUserId := len(dbStructure.Users) + 1

	newUser, err := NewUser(newUserId, email, password)
	if err != nil {
		return User{}, err
	}

	if _, exists := dbStructure.Users[newUserId]; exists {
		return User{}, errors.New("the user with this email already exists")
	}
	dbStructure.Users[newUserId] = *newUser
	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return *newUser, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	users := dbStructure.Users

	for _, user := range users {
		if user.Email == email {
			return user, nil
		}
	}

	errorMessage := fmt.Sprintf("the user with email = %v dosent exists", email)
	return User{}, errors.New(errorMessage)
}

func (db *DB) UpdateUser(updatedUser User) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	users := dbStructure.Users

	_, exists := users[updatedUser.ID]
	if !exists {
		errorMessage := fmt.Sprintf("the user with id = %v dosent exists", updatedUser.ID)
		return errors.New(errorMessage)
	}

	users[updatedUser.ID] = updatedUser

	err = db.writeDB(dbStructure)
	if err != nil {
		return nil
	}

	return nil
}

func (db *DB) CreateChirp(body string, authorID int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	newID := 0
	for chirpID := range dbStructure.Chirps {
		if chirpID > newID {
			newID = chirpID
		}
	}
	newID++

	chirp, err := NewChirp(body, newID, authorID)
	if err != nil {
		return Chirp{}, err
	}

	dbStructure.Chirps[newID] = *chirp
	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return *chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirpsMap := dbStructure.Chirps
	chirpsSlice := []Chirp{}

	for _, chirp := range chirpsMap {
		chirpsSlice = append(chirpsSlice, chirp)
	}

	sort.Slice(chirpsSlice, func(i, j int) bool {
		return chirpsSlice[i].ID < chirpsSlice[j].ID
	})

	return chirpsSlice, nil
}

func (db *DB) GetChirpByID(chirpID int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirpsMap := dbStructure.Chirps

	chirp, exists := chirpsMap[chirpID]
	if !exists {
		errorMessage := fmt.Sprintf("The chirp with id = %v was not found", chirpID)
		return Chirp{}, errors.New(errorMessage)
	}

	return chirp, nil
}

func (db *DB) DeleteChirpByID(chirpID int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	delete(dbStructure.Chirps, chirpID)

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) CreateRefreshToken(id int) (RefreshToken, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return RefreshToken{}, err
	}

	refreshToken, err := NewRefreshToken(id)
	if err != nil {
		return RefreshToken{}, err
	}

	dbStructure.RefreshTokens[id] = *refreshToken

	err = db.writeDB(dbStructure)
	if err != nil {
		return RefreshToken{}, err
	}

	return *refreshToken, nil
}

func (db *DB) GetRefreshTokenInfo(refreshToken string) (RefreshToken, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return RefreshToken{}, err
	}

	refreshTokens := dbStructure.RefreshTokens

	for _, storedRefToken := range refreshTokens {
		if storedRefToken.Token == refreshToken {
			return storedRefToken, nil
		}
	}
	errorMessage := fmt.Sprintf("refresh token = %v was not found", refreshToken)
	return RefreshToken{}, errors.New(errorMessage)
}

func (db *DB) DeleteRefreshToken(refreshToken string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	refreshTokens := dbStructure.RefreshTokens

	keyToDelete := -1
	for key, storedRefToken := range refreshTokens {
		if storedRefToken.Token == refreshToken {
			keyToDelete = key
		}
	}

	if keyToDelete == -1 {
		errorMessage := fmt.Sprintf("refresh token = %v was not found", refreshToken)
		return errors.New(errorMessage)
	}

	delete(refreshTokens, keyToDelete)

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		initialData, err := NewDBStructure(make(map[int]Chirp), make(map[int]User), make(map[int]RefreshToken))
		if err != nil {
			return err
		}
		return db.writeDB(*initialData)
	}
	return err
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	data, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	return os.WriteFile(db.path, data, 0644)
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	data, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}

	var dbStructure DBStructure
	err = json.Unmarshal(data, &dbStructure)
	if err != nil {
		return DBStructure{}, err
	}

	return dbStructure, nil
}
