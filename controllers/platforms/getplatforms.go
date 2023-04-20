package platforms

import (
	"log"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/webstradev/rsdb-backend/utils"
)

func GetPlatforms(env *utils.Environment) gin.HandlerFunc {
	return func(c *gin.Context) {
		page := c.MustGet("page").(int)
		pageSize := c.MustGet("pageSize").(int)

		platforms, err := env.DB.GetPlatforms(page, pageSize)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		count, err := env.DB.CountPlatforms()
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{"platforms": platforms, "total": count, "numPages": math.Ceil(float64(count) / float64(pageSize))})
	}
}
