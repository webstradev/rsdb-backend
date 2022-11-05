package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Database struct {
	querier *sqlx.DB
}

// Setup makes a connection to the db (sql.Open represents not 1 connection but actually represents a thread pool)
func Setup(dsn string) (*Database, error) {
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return &Database{querier: db}, nil
}

func (db *Database) Ping() error {
	err := db.querier.Ping()
	if err != nil {
		return err
	}
	return nil
}
