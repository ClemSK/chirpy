package database

import "time"

type Revocation struct {
	Token     string `json:token`
	RevokedAt string `json:revoked_at`
}

func (db *DB) RevokedToken(token string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}
	revocation := Revocation{
		Token:     token,
		RevokedAt: time.Now().UTC(),
	}
	dbStructure.Revocation[token] = revocation
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}
	return nil
}
