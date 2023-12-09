package database

import "errors"

type User struct {
	ID             int    `json:"id"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
}

var ErrAlreadyExists = errors.New("already exists")

func (db *DB) CreateUser(email, hashedPassword string) (User, error) {
	if _, err := db.GetUserByEmail(email); !errors.Is(err, ErrNotExist) {
		return User{}, ErrAlreadyExists
	}
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	// Initialize Chirps map if it's nil
	if dbStructure.Users == nil {
		dbStructure.Users = make(map[int]User)
	}

	id := len(dbStructure.Users) + 1
	user := User{
		ID:             id,
		Email:          email,
		HashedPassword: hashedPassword,
	}
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (db *DB) GetUsers() ([]User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	users := make([]User, 0, len(dbStructure.Users))
	for _, user := range dbStructure.Users {
		users = append(users, user)
	}
	return users, nil
}

func (db *DB) GetUser(id int) (User, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	user, found := dbStructure.Users[id]
	if !found {
		return User{}, errors.New("chirp not found")
	}
	return user, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dBStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dBStructure.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return User{}, ErrNotExist
}