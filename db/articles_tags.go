package db

type ArticleTag struct {
	ArticleId int64  `json:"-" db:"article_id"`
	TagId     int64  `json:"id" db:"tag_id"`
	Tag       string `json:"tag" db:"tag"`
}

func (db *Database) GetArticleTags(id int64) ([]ArticleTag, error) {
	tags := []ArticleTag{}

	err := db.querier.Select(&tags, `
	SELECT 
		at.*,
		t.tag
	FROM 
		tags t 
	LEFT JOIN 
		articles_tags at ON at.tag_id = t.id 
	WHERE 
		at.article_id = ? AND t.deleted_at IS NULL`, id)
	if err != nil {
		return nil, err
	}

	return tags, nil
}
