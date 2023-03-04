package db

import (
	"log"
	"time"
)

type Article struct {
	Model
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Link        string    `json:"link" db:"link"`
	Date        time.Time `json:"date" db:"date"`
	Body        string    `json:"body" db:"body"`
	Tags        []Tag     `json:"tags" db:"tags"`
}

type ArticleWithTagString struct {
	Article
	TagString string `json:"tagString" db:"article_tags"`
}

func (db *Database) CountArticles() (int, error) {
	var count int
	err := db.querier.Get(&count, "SELECT COUNT(*) AS count FROM articles WHERE deleted_at IS NULL")
	return count, err
}

func (db *Database) GetArticles(page, pageSize int) ([]ArticleWithTagString, error) {
	articles := []ArticleWithTagString{}

	err := db.querier.Select(&articles, `
	SELECT 
		a.*, 
		COALESCE(GROUP_CONCAT(DISTINCT t.tag), '') AS article_tags
	FROM 
		articles a 
	LEFT JOIN 
		articles_tags at ON at.article_id = a.id
	LEFT JOIN
		tags t ON t.id = at.tag_id
	WHERE a.deleted_at IS NULL
	GROUP BY a.id
	LIMIT ? OFFSET ?`, pageSize, (page-1)*pageSize)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return articles, nil
}
