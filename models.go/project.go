package models

import "time"

type Project struct {
	Model
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Link        string    `json:"link" db:"link"`
	Date        time.Time `json:"date" db:"date"`
	Body        string    `json:"body" db:"body"`
}
