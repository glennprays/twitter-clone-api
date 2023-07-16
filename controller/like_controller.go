package controller

import (
	"log"
	"net/http"
	"strconv"
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
	tweetIDInt, err := strconv.ParseInt(tweetID, 10, 64)
	if err != nil {
		log.Fatalln(err)
	}

	query := `
		MATCH (u:User { username: $username }), (t:Tweet)
		WHERE id(t) = $tweetID
		CREATE (u)-[:LIKES { timestamp: datetime() }]->(t)
	`

	result, err := session.Run(c, query, map[string]interface{}{
		"username": username,
		"tweetID":  tweetIDInt,
	})
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like"})
		return
	}

	// Check if the like was successfully created
	if result.Err() == nil {
		log.Println(result.Err())
		c.JSON(http.StatusOK, gin.H{"message": "Tweet liked"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like tweet"})
	}
}

func UnLikeTweet(c *gin.Context) {

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
	tweetIDInt, err := strconv.ParseInt(tweetID, 10, 64)
	if err != nil {
		log.Fatalln(err)
	}

	query := `
		MATCH (u:User { username: $username })-[l:LIKES]->(t:Tweet)
		WHERE id(t) = $tweetID
		DELETE l
	`

	result, err := session.Run(c, query, map[string]interface{}{
		"username": username,
		"tweetID":  tweetIDInt,
	})
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Unlike"})
		return
	}

	// Check if the like was successfully created
	if result.Err() == nil {
		log.Println(result.Err())
		c.JSON(http.StatusOK, gin.H{"message": "Tweet unliked"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unlike tweet"})
	}
}

func GetTweetLikes(c *gin.Context) {
	driver, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer driver.Close(c)

	session := driver.NewSession(c, neo4j.SessionConfig{})
	defer session.Close(c)

	tweetID := c.Param("tweetID")
	tweetIDInt, err := strconv.ParseInt(tweetID, 10, 64)
	if err != nil {
		log.Fatalln(err)
	}

	query := `
		MATCH (u:User)-[:LIKES]->(t:Tweet)
		WHERE id(t) = $tweetID
		RETURN u.username AS username
	`

	result, err := session.Run(c, query, map[string]interface{}{
		"tweetID": tweetIDInt,
	})
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get data"})
		return
	}

	var likes []string
	for result.Next(c) {
		record := result.Record()
		username, ok := record.Get("username")
		if ok {
			likes = append(likes, username.(string))
		}
	}

	c.JSON(http.StatusOK, gin.H{"likes": likes})
}
