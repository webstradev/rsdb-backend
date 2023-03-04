package db

type ArticlePlatform struct {
	ArticleId    int64  `json:"-" db:"article_id"`
	PlatformId   int64  `json:"id" db:"platform_id"`
	PlatformName string `json:"platform" db:"platform_name"`
}

func (db *Database) GetPlatformsForArticle(id int64) ([]ArticlePlatform, error) {
	platforms := []ArticlePlatform{}

	err := db.querier.Select(&platforms, `
	SELECT 
		pa.*,
		p.name as platform_name
	FROM 
		platforms p 
	LEFT JOIN 
		platforms_articles pa ON pa.platform_id = p.id 
	WHERE 
		pa.article_id = ? AND p.deleted_at IS NULL`, id)
	if err != nil {
		return nil, err
	}

	return platforms, nil
}
