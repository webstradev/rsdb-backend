package db

import (
	"database/sql"
	"log"
	"strings"
)

type Project struct {
	Model
	Title       string            `json:"title" db:"title"`
	Description string            `json:"description" db:"description"`
	Link        string            `json:"link" db:"link"`
	Date        sql.NullTime      `json:"date" db:"date"`
	Body        string            `json:"body" db:"body"`
	Tags        []ProjectTag      `json:"tags" db:"tags"`
	Platforms   []ProjectPlatform `json:"platforms" db:"platforms"`
}

type ProjectWithTagString struct {
	Project
	TagString string `json:"tagString" db:"project_tags"`
}

func (db *Database) CountProjects() (int, error) {
	var count int
	err := db.querier.Get(&count, "SELECT COUNT(*) AS count FROM projects WHERE deleted_at IS NULL")
	return count, err
}

func (p *Project) PopulateTags(db *Database) error {
	tags, err := db.GetProjectTags(p.ID)
	if err != nil {
		return err
	}

	p.Tags = tags

	return nil
}

func (p *Project) PopulatePlatforms(db *Database) error {
	platforms, err := db.GetPlatformsForProject(p.ID)
	if err != nil {
		return err
	}

	p.Platforms = platforms

	return nil
}

func (db *Database) GetProject(id int64) (Project, error) {
	project := Project{}

	err := db.querier.Get(&project, "SELECT p.* FROM projects p WHERE p.id = ? AND p.deleted_at IS NULL", id)
	return project, err
}

func (db *Database) GetProjects(page, pageSize int) ([]ProjectWithTagString, error) {
	projects := []ProjectWithTagString{}

	err := db.querier.Select(&projects, `
	SELECT 
		p.*, 
		COALESCE(GROUP_CONCAT(DISTINCT t.tag), '') AS project_tags
	FROM 
		projects p
	LEFT JOIN 
		projects_tags pt ON pt.project_id = p.id
	LEFT JOIN
		tags t ON t.id = pt.tag_id
	WHERE p.deleted_at IS NULL
	GROUP BY p.id
	LIMIT ? OFFSET ?`, pageSize, (page-1)*pageSize)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return projects, nil
}

func (db *Database) InsertProject(project Project) (int64, error) {
	result, err := db.querier.NamedExec(`
	INSERT INTO projects (title, description, link, date, body)
	VALUES (:title, :description, :link, :date, :body)`, project)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return id, nil
}

func (db *Database) InsertProjectTags(projectId int64, tags []int64) error {
	// If there are no tags, we're done
	if len(tags) == 0 {
		return nil
	}

	// Build a query and args
	query := `INSERT INTO projects_tags (tag_id, project_id) VALUES `
	args := []any{}

	for _, tagId := range tags {
		// Add placeholders for each category
		query += "(?, ?),"

		// Add bind arguments for each category
		args = append(args, tagId, projectId)
	}

	// Remove the last comma
	query = strings.TrimRight(query, ",")

	_, err := db.querier.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) InsertProjectPlatforms(projectId int64, platforms []int64) error {
	// If there are no categories, we're done
	if len(platforms) == 0 {
		return nil
	}

	// Build a query and args
	query := `INSERT INTO platforms_projects (platform_id, project_id) VALUES `
	args := []any{}

	for _, platformId := range platforms {
		// Add placeholders for each category
		query += "(?, ?),"

		// Add bind arguments for each category
		args = append(args, platformId, projectId)
	}

	// Remove the last comma
	query = strings.TrimRight(query, ",")

	_, err := db.querier.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}
