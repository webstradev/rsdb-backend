package db

import (
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type Article struct {
	Model
	Title       string            `json:"title" db:"title"`
	Description string            `json:"description" db:"description"`
	Link        string            `json:"link" db:"link"`
	Date        time.Time         `json:"date" db:"date"`
	Body        string            `json:"body" db:"body"`
	Tags        []ArticleTag      `json:"tags"`
	Platforms   []ArticlePlatform `json:"platforms"`
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

func (a *Article) TagIds() []int64 {
	ids := []int64{}

	for _, tag := range a.Tags {
		ids = append(ids, tag.TagId)
	}

	return ids
}

func (a *Article) PlatformIds() []int64 {
	ids := []int64{}

	for _, tag := range a.Platforms {
		ids = append(ids, tag.PlatformId)
	}

	return ids
}

func (a *Article) PopulateTags(db *Database) error {
	tags, err := db.GetArticleTags(a.ID)
	if err != nil {
		return err
	}

	a.Tags = tags

	return nil
}

func (a *Article) PopulatePlatforms(db *Database) error {
	platforms, err := db.GetPlatformsForArticle(a.ID)
	if err != nil {
		return err
	}

	a.Platforms = platforms

	return nil
}

func (db *Database) GetArticle(id int64) (Article, error) {
	article := Article{}

	err := db.querier.Get(&article, "SELECT a.* FROM articles a WHERE a.id = ? AND a.deleted_at IS NULL", id)

	return article, err
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

func (db *Database) InsertArticleTagsTx(tx *sqlx.Tx, articleID int64, tagIDs []int64) error {
	if len(tagIDs) == 0 {
		return nil
	}

	// Build the query and arguments
	query := "INSERT INTO articles_tags (tag_id, article_id) VALUES "
	args := make([]interface{}, 0, len(tagIDs)*2)
	for i, platformID := range tagIDs {
		query += "(?, ?),"
		args = append(args, platformID, articleID)
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

func (db *Database) InsertArticlePlatformsTx(tx *sqlx.Tx, articleID int64, platforms []int64) error {
	if len(platforms) == 0 {
		return nil
	}

	// Build the query and arguments
	query := "INSERT INTO platforms_articles (platform_id, article_id) VALUES "
	args := make([]interface{}, 0, len(platforms)*2)
	for i, platformID := range platforms {
		query += "(?, ?),"
		args = append(args, platformID, articleID)
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

func (db *Database) EditArticle(article Article) error {
	tx, err := db.querier.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(`
    UPDATE articles SET title = :title, description = :description, link = :link, date = :date, body = :body
    WHERE id = :id`, article)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("DELETE FROM articles_tags WHERE article_id = ?", article.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = db.InsertArticleTagsTx(tx, article.ID, article.TagIds())
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("DELETE FROM platforms_articles WHERE article_id = ?", article.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = db.InsertArticlePlatformsTx(tx, article.ID, article.PlatformIds())
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
