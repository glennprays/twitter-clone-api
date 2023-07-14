package controller

import (
	"log"
	"net/http"
	"time"
	"twitter-clone-api/config/database"
	"twitter-clone-api/middleware"
	"twitter-clone-api/models"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func PostTweet(c *gin.Context) {
	driver, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer driver.Close(c)

	session := driver.NewSession(c, neo4j.SessionConfig{})
	defer session.Close(c)

	var tweet models.Tweet
	if err := c.ShouldBindJSON(&tweet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username, _, _ := middleware.GetUsernameAndRoleFromCookie(c)

	query := `
		MATCH (u:User { username: $username })

		CREATE (t:Tweet {
		content: $content,
		timestamp: datetime()
		})

		CREATE (u)-[:POSTED]->(t)
		return id(t) as nodeId, t.timestamp as timestamp
	`

	result, err := session.Run(c,
		query,
		map[string]any{
			"username": username,
			"content":  tweet.Content,
		})

	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tweet"})
		return
	}

	if result.Next(c) {
		var createdTweet models.Tweet
		nodeID, ok := result.Record().Get("nodeId")
		if ok {
			nodeID, ok := nodeID.(int64)
			if ok {
				createdTweet.ID = &nodeID
			}
		}
		timestamp, ok := result.Record().Get("timestamp")
		if ok {
			timestamp, ok := timestamp.(time.Time)
			if ok {
				createdTweet.Timestamp = &timestamp
			}
		}
		createdTweet.Content = tweet.Content
		c.JSON(http.StatusOK, createdTweet)
	}

}
