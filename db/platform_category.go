package db

type PlatformCategory struct {
	PlatformID int64  `json:"-" db:"platform_id"`
	CategoryID int64  `json:"id" db:"category_id"`
	Category   string `json:"category" db:"category"`
}

func (db *Database) GetPlatformCategories(id int64) ([]PlatformCategory, error) {
	categories := []PlatformCategory{}

	err := db.querier.Select(&categories, `
	SELECT 
		pc.*,
		c.category
	FROM 
		categories c 
	LEFT JOIN 
		platforms_categories pc ON pc.category_id = c.id 
	WHERE 
		pc.platform_id = ? AND c.deleted_at IS NULL`, id)
	if err != nil {
		return nil, err
	}

	return categories, nil
}
