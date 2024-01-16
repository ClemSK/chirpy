package database

import "errors"

type Chirp struct {
	AuthorID int    `json:"author_id"`
	Body     string `json:"body"`
	ID       int    `json:"id"`
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string, authorID int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	// Initialize Chirps map if it's nil
	if dbStructure.Chirps == nil {
		dbStructure.Chirps = make(map[int]Chirp)
	}

	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		AuthorID: authorID,
		Body:     body,
		ID:       id,
	}
	dbStructure.Chirps[id] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}
	return chirp, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}
	return chirps, nil
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	chirp, found := dbStructure.Chirps[id]
	if !found {
		return Chirp{}, errors.New("chirp not found")
	}
	return chirp, nil
}
