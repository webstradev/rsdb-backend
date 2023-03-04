package db

import (
	"log"
	"strings"
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

func (db *Database) InsertArticle(article Article) (int64, error) {
	result, err := db.querier.NamedExec(`
	INSERT INTO articles (title, description, link, date, body)
	VALUES (:title, :description, :link, :date, :body)`, article)
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

func (db *Database) InsertArticleTags(articleId int64, tags []int64) error {
	// If there are no categories, we're done
	if len(tags) == 0 {
		return nil
	}

	// Build a query and args
	query := `INSERT INTO articles_tags (tag_id, article_id) VALUES `
	args := []any{}

	for _, tagId := range tags {
		// Add placeholders for each category
		query += "(?, ?),"

		// Add bind arguments for each category
		args = append(args, tagId, articleId)
	}

	// Remove the last comma
	query = strings.TrimRight(query, ",")

	_, err := db.querier.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) InsertArticlePlatforms(articleId int64, platforms []int64) error {
	// If there are no categories, we're done
	if len(platforms) == 0 {
		return nil
	}

	// Build a query and args
	query := `INSERT INTO platforms_articles (platform_id, article_id) VALUES `
	args := []any{}

	for _, platformId := range platforms {
		// Add placeholders for each category
		query += "(?, ?),"

		// Add bind arguments for each category
		args = append(args, platformId, articleId)
	}

	// Remove the last comma
	query = strings.TrimRight(query, ",")

	_, err := db.querier.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}
