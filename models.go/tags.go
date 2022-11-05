package models

type Tags struct {
	Model
	Tag string `json:"tag" db:"tag"`
}
