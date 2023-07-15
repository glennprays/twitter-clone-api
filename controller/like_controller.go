package controller

import (
	"log"
	"net/http"
	"twitter-clone-api/config/database"
	"twitter-clone-api/middleware"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func LikeTweet(c *gin.Context) {

	driver, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer driver.Close(c)

	session := driver.NewSession(c, neo4j.SessionConfig{})
	defer session.Close(c)

	username, _, _ := middleware.GetUsernameAndRoleFromCookie(c)
	tweetID := c.Param("tweetID")

	query := `
		MATCH (u:User { username: $username })
		MATCH (t:Tweet)
		WHERE id(t) = $tweetID
		CREATE (u)-[:LIKES { timestamp: datetime() }]->(t)
		return id(t)
	`

	result, err := session.Run(c, query, map[string]interface{}{
		"username": username,
		"tweetID":  tweetID,
	})
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like"})
		return
	}

	// Check if the like was successfully created
	if result.Err() != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Tweet liked"})
	} else {
		log.Println(result.Err())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like tweet"})
	}
}
