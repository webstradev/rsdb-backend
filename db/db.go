package db

import (
	"github.com/jmoiron/sqlx"
)

// Setup makes a connection to the db (sql.Open represents not 1 connection but actually represents a thread pool)
func Setup(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Ping the DB first
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
