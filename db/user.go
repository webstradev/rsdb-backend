package db

type User struct {
	Model
	Email    string `json:"email" db:"email"`
	Password string `json:"-" db:"password"`
	Role     string `json:"role" db:"role"`
}

func (db *Database) GetUserWithEmail(email string) (*User, error) {
	user := User{}
	err := db.querier.Get(&user, "SELECT * FROM users WHERE email = ?", email)
	return &user, err
}
