package platforms

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/webstradev/rsdb-backend/utils"
)

func GetContacts(env *utils.Environment) gin.HandlerFunc {
	return func(c *gin.Context) {
		idString := c.Param("platformId")

		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		contacts, err := env.DB.GetContactsForPlatform(id)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, contacts)
	}
}
