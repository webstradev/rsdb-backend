package db

const (
	REGISTRATION_TYPE = "REGISTRATION"
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
