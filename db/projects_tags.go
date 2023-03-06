package db

type ProjectTag struct {
	ProjectId int64  `json:"-" db:"project_id"`
	TagId     int64  `json:"id" db:"tag_id"`
	Tag       string `json:"tag" db:"tag"`
}

func (db *Database) GetProjectTags(id int64) ([]ProjectTag, error) {
	tags := []ProjectTag{}

	err := db.querier.Select(&tags, `
	SELECT 
		pt.*,
		t.tag
	FROM 
		tags t 
	LEFT JOIN 
		projects_tags pt ON pt.tag_id = t.id 
	WHERE 
		pt.project_id = ? AND t.deleted_at IS NULL`, id)

	return tags, err
}
