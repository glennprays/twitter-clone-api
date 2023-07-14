package routes

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// User represents a user node in the Neo4j database.
type User struct {
	Name     string    `json:"name"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Created  time.Time `json:"created"`
}

func GetPeopleHandler(driver neo4j.Driver) gin.HandlerFunc {
	return func(c *gin.Context) {
		people, err := getPeople(driver)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		c.JSON(http.StatusOK, people)
	}
}

// getPeople retrieves all people from the Neo4j database.
func getPeople(driver neo4j.Driver) ([]User, error) {
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	result, err := session.Run("MATCH (p:User) RETURN p.name AS name, p.email AS email, p.username AS username, p.created AS created, p.password AS password", nil)
	if err != nil {
		return nil, err
	}

	var people []User

	for result.Next() {
		record := result.Record()

		nameVal, nameExists := record.Get("name")
		emailVal, emailExists := record.Get("email")
		usernameVal, usernameExists := record.Get("username")
		createdVal, createdExists := record.Get("created")
		passwordVal, passwordExists := record.Get("password")

		if !nameExists || !emailExists || !usernameExists || !createdExists || !passwordExists {
			continue
		}

		user := User{
			Name:     nameVal.(string),
			Email:    emailVal.(string),
			Username: usernameVal.(string),
			Created:  createdVal.(time.Time),
			Password: passwordVal.(string),
		}
		people = append(people, user)
	}

	return people, nil
}

func CreateUserHandler(driver neo4j.Driver) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user User

		err := c.ShouldBindJSON(&user)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
			return
		}

		err = createUser(driver, user)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
	}
}

// // createUser creates a new user node in the Neo4j database.
func createUser(driver neo4j.Driver, user User) error {
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err := session.Run("CREATE (u:User {name: $name, username: $username, email: $email, password: $password, created: datetime()})", map[string]interface{}{
		"name":     user.Name,
		"username": user.Username,
		"email":    user.Email,
		"password": user.Password,
	})
	if err != nil {
		return err
	}

	return nil
}
