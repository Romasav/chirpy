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
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

func NewDBStructure(chirps map[int]Chirp, users map[int]User) (*DBStructure, error) {
	newDBStructure := DBStructure{
		Chirps: chirps,
		Users:  users,
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

func (db *DB) CreateUser(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	newUserId := len(dbStructure.Users) + 1

	newUser, err := NewUser(newUserId, email)
	if err != nil {
		return User{}, err
	}
	dbStructure.Users[newUserId] = *newUser
	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return *newUser, nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	newId := len(dbStructure.Chirps) + 1

	chirp, err := NewChirp(body, newId)
	if err != nil {
		return Chirp{}, err
	}

	dbStructure.Chirps[newId] = *chirp
	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, nil
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

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		initialData, err := NewDBStructure(make(map[int]Chirp), make(map[int]User))
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
