package db

import (
	"database/sql"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
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

func (p *Project) TagIds() []int64 {
	ids := []int64{}

	for _, tag := range p.Tags {
		ids = append(ids, tag.TagId)
	}

	return ids
}

func (p *Project) PlatformIds() []int64 {
	ids := []int64{}

	for _, tag := range p.Platforms {
		ids = append(ids, tag.PlatformId)
	}

	return ids
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

func (db *Database) InsertProjectTagsTx(tx *sqlx.Tx, projectID int64, tagIDs []int64) error {
	if len(tagIDs) == 0 {
		return nil
	}

	// Build the query and arguments
	query := "INSERT INTO projects_tags (tag_id, project_id) VALUES "
	args := make([]interface{}, 0, len(tagIDs)*2)
	for i, platformID := range tagIDs {
		query += "(?, ?),"
		args = append(args, platformID, projectID)
		if i == len(tagIDs)-1 {
			query = query[:len(query)-1]
		}
	}

	// Execute the query
	_, err := tx.Exec(query, args...)
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

func (db *Database) InsertProjectPlatformsTx(tx *sqlx.Tx, projectID int64, platforms []int64) error {
	if len(platforms) == 0 {
		return nil
	}

	// Build the query and arguments
	query := "INSERT INTO platforms_projects (platform_id, project_id) VALUES "
	args := make([]interface{}, 0, len(platforms)*2)
	for i, platformID := range platforms {
		query += "(?, ?),"
		args = append(args, platformID, projectID)
		if i == len(platforms)-1 {
			query = query[:len(query)-1]
		}
	}

	// Execute the query
	_, err := tx.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) EditProject(project Project) error {
	tx, err := db.querier.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(`
    UPDATE projects SET title = :title, description = :description, link = :link, date = :date, body = :body
    WHERE id = :id`, project)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("DELETE FROM projects_tags WHERE project_id = ?", project.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = db.InsertProjectTagsTx(tx, project.ID, project.TagIds())
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("DELETE FROM platforms_projects WHERE project_id = ?", project.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = db.InsertProjectPlatformsTx(tx, project.ID, project.PlatformIds())
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) DeleteProject(id int64) error {
	_, err := db.querier.Exec("UPDATE projects SET deleted_at = CURRENT_TIMESTAMP() WHERE id = ?", id)
	return err
}
