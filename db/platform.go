package db

type Platform struct {
	Model
	Name     string `json:"name" db:"name"`
	Website  string `json:"website" db:"website"`
	Country  string `json:"country" db:"country"`
	Source   string `json:"source" db:"source"`
	Notes    string `json:"notes" db:"notes"`
	Comments string `json:"comments" db:"comments"`
}

func (db *Database) GetAllPlatforms() ([]Platform, error) {
	platforms := []Platform{}

	err := db.querier.Select(&platforms, "SELECT * FROM platforms WHERE deleted_at IS NULL")
	if err != nil {
		return nil, err
	}

	return platforms, nil
}
