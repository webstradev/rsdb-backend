package db

type Category struct {
	Model
	Category string `json:"category" db:"category"`
}

func (db *Database) GetAllCategories() ([]Tag, error) {
	tags := []Tag{}

	err := db.querier.Select(&tags, "SELECT * FROM categories WHERE deleted_at IS NULL")
	if err != nil {
		return nil, err
	}

	return tags, nil
}
