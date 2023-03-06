package db

type ProjectPlatform struct {
	ProjectId    int64  `json:"-" db:"project_id"`
	PlatformId   int64  `json:"id" db:"platform_id"`
	PlatformName string `json:"platform" db:"platform_name"`
}

func (db *Database) GetPlatformsForProject(id int64) ([]ProjectPlatform, error) {
	platforms := []ProjectPlatform{}

	err := db.querier.Select(&platforms, `
	SELECT 
		pp.*,
		p.name as platform_name
	FROM 
		platforms p 
	LEFT JOIN 
		platforms_projects pp ON pp.platform_id = p.id 
	WHERE 
		pp.project_id = ? AND p.deleted_at IS NULL`, id)
	if err != nil {
		return nil, err
	}

	return platforms, nil
}
