package db

type User struct {
	Model
	Email    string `json:"email" db:"email"`
	Password string `json:"-" db:"password"`
	Role     string `json:"role" db:"role"`
}

func (db *Database) GetUserWithEmail(email string) (*User, error) {
	user := User{}
	err := db.querier.Get(&user, "SELECT * FROM users WHERE email = ? AND deleted_at IS NULL", email)
	return &user, err
}

func (db *Database) IsUsernameAvailable(email string) (bool, error) {
	var count int64
	err := db.querier.Get(&count, "SELECT COUNT(*) FROM users WHERE email = ? and deleted_at IS NULL", email)
	return count == 0, err
}

func (db *Database) InsertUser(u User) error {
	_, err := db.querier.NamedExec("INSERT INTO users (email, password, role) VALUES (:email, :password, :role)", u)
	return err
}
