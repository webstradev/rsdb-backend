package models

type Contact struct {
	Model
	Name    string `json:"name" db:"name"`
	Title   string `json:"title" db:"title"`
	Email   string `json:"email" db:"email"`
	Phone   string `json:"phone" db:"phone"`
	Phone2  string `json:"phone2" db:"phone2"`
	Address string `json:"address" db:"address"`
	Notes   string `json:"notes" db:"notes"`
	Source  string `json:"source" db:"source"`
	Privacy string `json:"privacy" db:"privacy"`
}
