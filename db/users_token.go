package db

import (
	"time"
)

const (
	REGISTRATION_TYPE         = "REGISTRATION"
	REGISTRATION_TOKEN_EXPIRY = 24 * time.Hour
	PASSWORD_RESET_TYPE       = "PASSWORD_RESET"
)

type UsersToken struct {
	ModelLite
	HashedToken string `json:"-" db:"hashed_token"`
	Type        string `json:"-" db:"type"`
	CreatedBy   int64  `json:"-" db:"created_by"`
	UserId      int64  `json:"-" db:"user_id"`
	Used        int64  `json:"-" db:"used"`
}

func (db *Database) InsertRegistrationToken(hashedToken string, createdBy int64) error {
	_, err := db.querier.Exec("INSERT INTO users_tokens (hashed_token, type, created_by) VALUES (?, ?, ?)", hashedToken, REGISTRATION_TYPE, createdBy)
	return err
}

func (db *Database) InsertPasswordResetToken(hashedToken string, userId, createdBy int64) error {
	_, err := db.querier.Exec("INSERT INTO users_tokens (hashed_token, type, user_id, created_by) VALUES (?, ?, ?, ?)", hashedToken, PASSWORD_RESET_TYPE, userId, createdBy)
	return err
}

func (db *Database) ValidateRegistrationToken(hashedToken string) (bool, error) {
	return db.validateToken(hashedToken, REGISTRATION_TYPE)
}

func (db *Database) validateToken(hashedToken, tokenType string) (bool, error) {
	var count int64
	// Make sure there is an unused token with the given hash and type that is not expired
	err := db.querier.Get(&count, `
	SELECT 
		COUNT(*) 
	FROM 
		users_tokens 
	WHERE 
		hashed_token = ? AND 
		type = ? AND 
		used = 0 AND 
		created_at > ?
	`, hashedToken, tokenType, time.Now().Add(-REGISTRATION_TOKEN_EXPIRY))
	return count > 0, err
}

func (db *Database) ConsumeToken(hashedToken string) error {
	_, err := db.querier.Exec("UPDATE users_tokens SET used = 1 WHERE hashed_token = ?", hashedToken)
	return err
}
