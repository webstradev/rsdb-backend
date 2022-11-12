package temp

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/webstradev/rsdb-backend/utils"
)

func MigrateProjectsTags(env *utils.Environment) {
	// Open our articles jsonFile
	jsonFile, err := os.Open("./temp/projects.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Fatalf("failed open file: %v", err)
	}

	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, _ := io.ReadAll(jsonFile)

	projects := []Project{}

	json.Unmarshal(byteValue, &projects)

	tags := []string{}

	for _, project := range projects {
		tags = append(tags, project.Tags...)
	}

	args := []any{}
	for _, tag := range tags {
		args = append(args, strings.TrimSpace(tag))
	}

	// here
	query := fmt.Sprintf(`
	INSERT IGNORE INTO tags (tag)
	VALUES 
	(?)%s`, strings.Repeat(",(?)", len(tags)-1))

	_, err = env.DB.Querier.Exec(query, args...)
	if err != nil {
		log.Printf("err: %v", err)
	}

	dbCategories := []struct {
		ID  int64  `db:"id"`
		Tag string `db:"tag"`
	}{}
	categoryToId := map[string]int64{}

	err = env.DB.Querier.Select(&dbCategories, "SELECT id, tag FROM tags")
	if err != nil {
		log.Printf("could not fetch categories: %v", err)
	}

	for _, dbCategory := range dbCategories {
		categoryToId[dbCategory.Tag] = dbCategory.ID
	}

	for _, project := range projects {
		// skip this article if there are no tags
		if len(project.Tags) == 0 {
			continue
		}

		projectId := 0
		err = env.DB.Querier.Get(&projectId, "SELECT id FROM articles WHERE object_id = ?", project.OldId.ObjectID)
		if err != nil {
			log.Printf("1:%v", err)
		}

		// go to next article if we cant find the article in the database
		if projectId == 0 {
			continue
		}

		query := fmt.Sprintf(`
		INSERT INTO projects_tags (project_id, tag_id)
		VALUES 
		(?, ?)%s`, strings.Repeat(",(?,?)", len(tags)-1))

		args := []any{}
		for _, tag := range project.Tags {
			args = append(args, projectId, categoryToId[tag])
		}

		_, err = env.DB.Querier.Exec(query, args...)
		if err != nil {
			log.Printf("err: %v", err)
		}
		log.Printf("migrated %d tags for project %d", len(tags), projectId)
	}
}
