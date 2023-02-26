package db

type Platform struct {
	Model
	Name          string `json:"name" db:"name"`
	Website       string `json:"website" db:"website"`
	Country       string `json:"country" db:"country"`
	Source        string `json:"source" db:"source"`
	Notes         string `json:"notes" db:"notes"`
	Comment       string `json:"comment" db:"comment"`
	Categories    string `json:"categories" db:"platform_categories"`
	ContactsCount int    `json:"contactsCount" db:"contacts_count"`
	ArticlesCount int    `json:"articlesCount" db:"articles_count"`
	ProjectsCount int    `json:"projectsCount" db:"projects_count"`
	Privacy       string `json:"privacy" db:"privacy"`
}

func (db *Database) GetPlatforms(page, pageSize int) ([]Platform, error) {
	platforms := []Platform{}

	err := db.querier.Select(&platforms, `
	SELECT 
		p.* , 
		COUNT(DISTINCT c.id) as contacts_count,
		COUNT(DISTINCT pa.article_id) as articles_count,
		COUNT(DISTINCT pp.project_id) as projects_count,
		GROUP_CONCAT(DISTINCT ca.category) AS platform_categories
	FROM 
		platforms p 
	LEFT JOIN 
		contacts c ON c.platform_id = p.id 
	LEFT JOIN 
		platforms_articles pa ON pa.platform_id = p.id
	LEFT JOIN 
		platforms_projects pp ON pp.platform_id = p.id
	LEFT JOIN 
		platforms_categories pc ON pc.platform_id = p.id
	LEFT JOIN
		categories ca ON ca.id = pc.category_id
	WHERE p.deleted_at IS NULL
	GROUP BY p.id
	LIMIT ? OFFSET ?`, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, err
	}

	return platforms, nil
}

func (db *Database) GetPlatform(id int64) (*Platform, error) {
	platform := Platform{}

	err := db.querier.Get(&platform, `
	SELECT 
		p.* , 
		COUNT(DISTINCT c.id) as contacts_count,
		COUNT(DISTINCT pa.article_id) as articles_count,
		COUNT(DISTINCT pp.project_id) as projects_count,
		GROUP_CONCAT(DISTINCT ca.category) AS platform_categories
	FROM 
		platforms p 
	LEFT JOIN 
		contacts c ON c.platform_id = p.id 
	LEFT JOIN 
		platforms_articles pa ON pa.platform_id = p.id
	LEFT JOIN 
		platforms_projects pp ON pp.platform_id = p.id
	LEFT JOIN 
		platforms_categories pc ON pc.platform_id = p.id
	LEFT JOIN
		categories ca ON ca.id = pc.category_id
	WHERE p.deleted_at IS NULL AND p.id = ?
	GROUP BY p.id`, id)
	if err != nil {
		return nil, err
	}

	return &platform, nil
}

func (db *Database) CountPlatforms() (int, error) {
	var count int
	err := db.querier.Get(&count, "SELECT COUNT(*) AS count FROM platforms WHERE deleted_at IS NULL")
	return count, err
}
