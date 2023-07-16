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

func FollowUser(c *gin.Context) {

	driver, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer driver.Close(c)

	session := driver.NewSession(c, neo4j.SessionConfig{})
	defer session.Close(c)

	username, _, _ := middleware.GetUsernameAndRoleFromCookie(c)
	userID := c.Param("userID")
	userIDInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		log.Fatalln(err)
	}
	query := `
		MATCH (u:User { username: $username }), (t:User)
		WHERE id(t) = $userID
		CREATE (u)-[:FOLLOWS]->(t)
	`

	result, err := session.Run(c, query, map[string]interface{}{
		"username": username,
		"userID":   userIDInt,
	})
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to follow"})
		return
	}

	// Check if the like was successfully created
	if result.Err() == nil {
		log.Println(result.Err())
		c.JSON(http.StatusOK, gin.H{"message": "User Followed"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to follow"})
	}
}

func UnFollowUser(c *gin.Context) {

	driver, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer driver.Close(c)

	session := driver.NewSession(c, neo4j.SessionConfig{})
	defer session.Close(c)

	username, _, _ := middleware.GetUsernameAndRoleFromCookie(c)
	userID := c.Param("userID")
	userIDInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		log.Fatalln(err)
	}
	query := `
		MATCH (u:User { username: $username })-[f:FOLLOWS]->(t:User)
		WHERE id(t) = $userID
		DELETE f
	`

	result, err := session.Run(c, query, map[string]interface{}{
		"username": username,
		"userID":   userIDInt,
	})
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to follow"})
		return
	}

	// Check if the like was successfully created
	if result.Err() == nil {
		log.Println(result.Err())
		c.JSON(http.StatusOK, gin.H{"message": "User Unfollowed"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to follow"})
	}
}

func GetFollowers(c *gin.Context) {
	driver, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer driver.Close(c)

	session := driver.NewSession(c, neo4j.SessionConfig{})
	defer session.Close(c)

	userID := c.Param("userID")
	userIDInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		log.Fatalln(err)
	}

	query := `
		MATCH (u:User)-[:FOLLOWS]->(t:User)
		WHERE id(t) = $userID
		RETURN u.username AS username
	`

	result, err := session.Run(c, query, map[string]interface{}{
		"userID": userIDInt,
	})
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get data"})
		return
	}

	var usernames []string
	for result.Next(c) {
		record := result.Record()
		username, ok := record.Get("username")
		if ok {
			usernames = append(usernames, username.(string))
		}
	}

	c.JSON(http.StatusOK, gin.H{"Followers": usernames})
}

func GetFollowing(c *gin.Context) {
	driver, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer driver.Close(c)

	session := driver.NewSession(c, neo4j.SessionConfig{})
	defer session.Close(c)

	userID := c.Param("userID")
	userIDInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		log.Fatalln(err)
	}

	query := `
		MATCH (u:User)-[:FOLLOWS]->(t:User)
		WHERE id(u) = $userID
		RETURN t.username AS username
	`

	result, err := session.Run(c, query, map[string]interface{}{
		"userID": userIDInt,
	})
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get data"})
		return
	}

	var usernames []string
	for result.Next(c) {
		record := result.Record()
		username, ok := record.Get("username")
		if ok {
			usernames = append(usernames, username.(string))
		}
	}

	c.JSON(http.StatusOK, gin.H{"Following": usernames})
}
