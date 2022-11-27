package controllers

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/webstradev/rsdb-backend/utils"
)

type Count struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type Counts map[string]int

var tablesToCount = []string{"platforms", "contacts", "articles", "projects"}

func GetCounts(env *utils.Environment) gin.HandlerFunc {
	return func(c *gin.Context) {
		counts := Counts{}

		// Get counts for each table
		for _, table := range tablesToCount {
			// Fetch count from database
			count, err := env.DB.CountRowsForTable(table)
			if err != nil {
				log.Println(err)
				c.AbortWithStatus(500)
				return
			}

			counts[table] = count
		}

		c.JSON(200, counts)
	}
}
