package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/webstradev/rsdb-backend/migrations"
)

type Database struct {
	querier  *sqlx.DB
	migrator *migrations.Sqlx
}

// Setup makes a connection to the db (sql.Open represents not 1 connection but actually represents a thread pool)
func Setup(dsn string, migrator *migrations.Sqlx) (*Database, error) {
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return &Database{querier: db, migrator: migrator}, nil
}

func SetupMockDB(mockDb *sqlx.DB) *Database {
	return &Database{querier: mockDb}
}

func (db *Database) Ping() error {
	err := db.querier.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) Migrate() error {
	err := db.migrator.Migrate(db.querier.DB, "mysql")
	return err
}
