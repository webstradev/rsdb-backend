package db

type Tag struct {
	Model
	Tag string `json:"tag" db:"tag"`
}

func (db *Database) GetAllTags() ([]Tag, error) {
	tags := []Tag{}

	err := db.querier.Select(&tags, "SELECT * FROM tags WHERE deleted_at IS NULL")
	if err != nil {
		return nil, err
	}

	return tags, nil
}
