package db

import "time"

type Article struct {
	Model
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Link        string    `json:"link" db:"link"`
	Date        time.Time `json:"date" db:"date"`
	Body        string    `json:"body" db:"body"`
}

func (db *Database) CountArticles() (int, error) {
	var count int
	err := db.querier.Get(&count, "SELECT COUNT(*) AS count FROM articles WHERE deleted_at IS NULL")
	return count, err
}
