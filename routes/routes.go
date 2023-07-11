package routes

import (
	"fmt"
	"net/http"
	"twitter-clone-api/config/database"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.GET("/", func(c *gin.Context) {
			driver, err := database.ConnectDB()
			if err != nil {
				fmt.Printf("Failed to connect to Neo4j: %v", err)
			} else {
				fmt.Println("Success Connect to DB")
			}
			driver.Close(c)
			c.JSON(http.StatusOK, gin.H{
				"message": "Welcome to Twitter-clone",
			})
		})
	}

	return router
}
