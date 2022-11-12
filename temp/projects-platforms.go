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

func MigrateProjectsToPlatforms(env *utils.Environment) {
	// Open our jsonFile
	jsonFile, err := os.Open("./temp/projects.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Fatalf("failed open file: %v", err)
	}

	// read our opened jsonFile as a byte array.
	byteValue, _ := io.ReadAll(jsonFile)

	projects := []Project{}

	json.Unmarshal(byteValue, &projects)

	for _, project := range projects {
		// go to next platform if there are no Platforms
		if len(project.Platforms) == 0 {
			continue
		}

		projectId := 0
		err = env.DB.Querier.Get(&projectId, "SELECT id FROM projects WHERE object_id = ?", project.OldId.ObjectID)
		if err != nil {
			log.Printf("1:%v", err)
		}

		// go to next article if we cant find the article in the database
		if projectId == 0 {
			continue
		}

		platformObjectIds := []any{}
		for _, platform := range project.Platforms {
			platformObjectIds = append(platformObjectIds, platform.ObjectID)
		}

		platformIds := []any{}

		err = env.DB.Querier.Select(&platformIds, fmt.Sprintf(`SELECT id FROM platforms WHERE object_id IN (?%s)`, strings.Repeat(",?", len(project.Platforms)-1)), platformObjectIds...)
		if err != nil {
			log.Printf("2:%v", err)
		}

		if len(platformIds) == 0 {
			continue
		}

		args := []any{}
		for _, plt := range platformIds {
			args = append(args, projectId, plt)
		}

		query := fmt.Sprintf(`
		INSERT INTO 
			platforms_projects (project_id, platform_id) 
		VALUES
			(?,?)%s`, strings.Repeat(",(?,?)", len(platformIds)-1))

		_, err = env.DB.Querier.Exec(query, args...)
		if err != nil {
			log.Printf("Error executing query: %s", query)
			log.Printf("3: %v", err)
		}

		log.Printf("migrated %d platforms for project %d", len(platformIds), projectId)
	}

}
