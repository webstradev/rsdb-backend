package db

import (
	"log"
	"strings"
)

type Platform struct {
	Model
	Name          string             `json:"name" db:"name"`
	Website       string             `json:"website" db:"website"`
	Country       string             `json:"country" db:"country"`
	Source        string             `json:"source" db:"source"`
	Notes         string             `json:"notes" db:"notes"`
	Comment       string             `json:"comment" db:"comment"`
	Categories    []PlatformCategory `json:"categories"`
	ContactsCount int                `json:"contactsCount" db:"contacts_count"`
	ArticlesCount int                `json:"articlesCount" db:"articles_count"`
	ProjectsCount int                `json:"projectsCount" db:"projects_count"`
	Privacy       string             `json:"privacy" db:"privacy"`
}

type PlatformWithCategoryString struct {
	Platform
	CategoryString string `json:"categoryString" db:"platform_categories"`
}

func (p *Platform) PopulateCategories(db *Database) error {
	categories, err := db.GetPlatformCategories(p.ID)
	if err != nil {
		return err
	}

	p.Categories = categories

	return nil
}

func (db *Database) GetPlatforms(page, pageSize int) ([]PlatformWithCategoryString, error) {
	platforms := []PlatformWithCategoryString{}

	err := db.querier.Select(&platforms, `
	SELECT 
		p.* , 
		COUNT(DISTINCT c.id) as contacts_count,
		COUNT(DISTINCT pa.article_id) as articles_count,
		COUNT(DISTINCT pp.project_id) as projects_count,
		COALESCE(GROUP_CONCAT(DISTINCT ca.category), '') AS platform_categories
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
		log.Println(err)
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
		COUNT(DISTINCT pp.project_id) as projects_count
	FROM 
		platforms p 
	LEFT JOIN 
		contacts c ON c.platform_id = p.id 
	LEFT JOIN 
		platforms_articles pa ON pa.platform_id = p.id
	LEFT JOIN  
		platforms_projects pp ON pp.platform_id = p.id
	WHERE p.deleted_at IS NULL AND p.id = ? GROUP BY p.id`, id)
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

func (db *Database) CreatePlatform(name, website, country, source, notes, comment, privacy string) (int64, error) {
	// Create the platform
	result, err := db.querier.Exec(`INSERT INTO platforms (name, website, country, source, notes, comment, privacy) VALUES (?, ?, ?, ?, ?, ?, ?)`, name, website, country, source, notes, comment, privacy)
	if err != nil {
		return -1, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (db *Database) EditPlatform(name, website, country, source, notes, comment, privacy string, id int64) error {
	// Update basic platform information
	_, err := db.querier.Exec(`UPDATE platforms SET name = ?, website = ?, country = ?, source = ?, notes = ?, comment = ?, privacy = ? WHERE id = ?`, name, website, country, source, notes, comment, privacy, id)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) InsertPlatformCategories(platformId int64, categories []int64) error {
	// If there are no categories, we're done
	if len(categories) == 0 {
		return nil
	}

	// Build a query and args
	query := `INSERT INTO platforms_categories (platform_id, category_id) VALUES `
	args := []any{}

	for _, categoryId := range categories {
		// Add placeholders for each category
		query += "(?, ?),"

		// Add bind arguments for each category
		args = append(args, platformId, categoryId)
	}

	// Remove the last comma
	query = strings.TrimRight(query, ",")

	_, err := db.querier.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdatePlatformCategories(platformId int64, categories []int64) error {
	// Delete all existing categories for this platform
	_, err := db.querier.Exec("DELETE FROM platforms_categories WHERE platform_id = ?", platformId)
	if err != nil {
		return err
	}

	// Insert the new categories
	err = db.InsertPlatformCategories(platformId, categories)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) DeletePlatform(platformId int64) error {
	_, err := db.querier.Exec("UPDATE platforms SET deleted_at = CURRENT_TIMESTAMP() WHERE id = ?", platformId)
	if err != nil {
		return err
	}

	return nil
}
