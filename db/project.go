package db

import (
	"database/sql"
	"log"
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
