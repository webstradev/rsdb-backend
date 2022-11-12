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

type Project struct {
	OldId struct {
		ObjectID string `json:"$oid"`
	} `json:"_id"`
	Title       string            `json:"title" db:"title"`
	Link        string            `json:"link" db:"link"`
	Description string            `json:"description" db:"description"`
	Date        string            `json:"date" db:"date"`
	Body        string            `json:"body" db:"body"`
	Platforms   []ObjectReference `json:"platforms"`
}

// Need to migrate dates from 2017 to 09/09/2017 and 2019 to 09/09/2019
func MigrateProjects(env *utils.Environment) {
	// Open our jsonFile
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

	// here
	query := fmt.Sprintf(`
	INSERT INTO projects (title, link, description, date, body, object_id)
	VALUES 
	(?, ?, ?, STR_TO_DATE(?, ?), ?, ?)%s`, strings.Repeat(",(?, ?, ?, STR_TO_DATE(?, ?), ?, ?)", len(projects)-1))

	args := []any{}

	for _, project := range projects {
		args = append(args, project.Title, project.Link, project.Description, project.Date, "%d/%m/%Y", project.Body, project.OldId.ObjectID)
	}

	_, err = env.DB.Querier.Exec(query, args...)
	if err != nil {
		log.Printf("err: %v", err)
	}

	_, err = env.DB.Querier.Exec("UPDATE projects SET date = NULL WHERE date = 0000-00-00")
	if err != nil {
		log.Printf("err: %v", err)
	}
}
