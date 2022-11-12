package db

type Platform struct {
	Model
	Name          string `json:"name" db:"name"`
	Website       string `json:"website" db:"website"`
	Country       string `json:"country" db:"country"`
	Source        string `json:"source" db:"source"`
	Notes         string `json:"notes" db:"notes"`
	Comments      string `json:"comments" db:"comments"`
	Categories    string `json:"categories" db:"platform_categories"`
	ContactsCount int    `json:"contactsCount" db:"contacts_count"`
	ArticlesCount int    `json:"articlesCount" db:"articles_count"`
	ProjectsCount int    `json:"projectsCount" db:"projects_count"`
}

func (db *Database) GetAllPlatforms() ([]Platform, error) {
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
	GROUP BY p.id`)
	if err != nil {
		return nil, err
	}

	return platforms, nil
}
