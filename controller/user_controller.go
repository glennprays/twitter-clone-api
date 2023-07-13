package controller

import (
	"log"
	"net/http"
	"twitter-clone-api/config/database"
	"twitter-clone-api/middleware"
	"twitter-clone-api/models"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func LoginUser(c *gin.Context) {
	driver, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer driver.Close(c)

	session := driver.NewSession(c, neo4j.SessionConfig{})
	defer session.Close(c)

	var login models.LoginRequest
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := session.Run(c,
		"MATCH (u:User {username: $username, password: $password}) RETURN count(u) = 1",
		map[string]any{
			"username": login.Username,
			"password": login.Password,
		})
	if err != nil {
		log.Fatal(err)
		return
	}

	if result.Next(c) {
		count := result.Record().Values[0].(bool)
		log.Println(count)
		if count {
			middleware.CreateToken(c, login.Username, "user", 3600)

			c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
			return
		}
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
}
