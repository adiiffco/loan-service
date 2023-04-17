package health

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Controller(c *gin.Context) {
	log.Println("health check api called")
	c.JSON(http.StatusOK, gin.H{
		"healthCheck": "Successful4",
	})
}
