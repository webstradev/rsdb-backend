package projects

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/webstradev/rsdb-backend/utils"
)

func GetProjects(env *utils.Environment) gin.HandlerFunc {
	return func(c *gin.Context) {
		page := c.MustGet("page").(int)
		pageSize := c.MustGet("pageSize").(int)

		projects, err := env.DB.GetProjects(page, pageSize)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, projects)
	}
}
